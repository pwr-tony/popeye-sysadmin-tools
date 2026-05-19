package main

import (
	"context"
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed configs/commands/*.yaml
var configsFS embed.FS

type App struct {
	ctx        context.Context
	categories []CommandCategory
	userDir    string
	docsDir    string
}

type CommandCategory struct {
	Category string    `yaml:"category" json:"category"`
	Icon     string    `yaml:"icon" json:"icon"`
	Commands []Command `yaml:"commands" json:"commands"`
	IsUser   bool      `yaml:"-" json:"isUser"`
}

type Command struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
	Command     string `yaml:"command" json:"command"`
	Sudo        bool   `yaml:"sudo" json:"sudo"`
}

type CommandInfo struct {
	Category    string `json:"category"`
	Icon        string `json:"icon"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Command     string `json:"command"`
	Sudo        bool   `json:"sudo"`
	IsUser      bool   `json:"isUser"`
}

func NewApp() *App {
	home, _ := os.UserHomeDir()
	userDir := filepath.Join(home, ".config", "popeye", "commands")
	docsDir := filepath.Join(home, "cmdb")
	return &App{userDir: userDir, docsDir: docsDir}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.initializeUserCommands()
	a.loadUserCommands()
}

func (a *App) initializeUserCommands() {
	if err := os.MkdirAll(a.userDir, 0755); err != nil {
		return
	}

	entries, err := fs.ReadDir(configsFS, "configs/commands")
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		destPath := filepath.Join(a.userDir, entry.Name())
		if _, err := os.Stat(destPath); err == nil {
			continue // already exists, skip
		}

		data, err := configsFS.ReadFile("configs/commands/" + entry.Name())
		if err != nil {
			continue
		}

		os.WriteFile(destPath, data, 0644)
	}
}

func (a *App) loadUserCommands() {
	if _, err := os.Stat(a.userDir); os.IsNotExist(err) {
		return
	}

	entries, err := os.ReadDir(a.userDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		data, err := os.ReadFile(filepath.Join(a.userDir, entry.Name()))
		if err != nil {
			continue
		}

		var cat CommandCategory
		if err := yaml.Unmarshal(data, &cat); err != nil {
			continue
		}
		cat.IsUser = true
		a.categories = append(a.categories, cat)
	}
}

func (a *App) GetCommands() []CommandInfo {
	var result []CommandInfo

	for _, cat := range a.categories {
		for _, cmd := range cat.Commands {
			cmdStr := strings.TrimSpace(cmd.Command)
			if cmd.Sudo {
				cmdStr = "sudo " + cmdStr
			}

			result = append(result, CommandInfo{
				Category:    cat.Category,
				Icon:        cat.Icon,
				Name:        cmd.Name,
				Description: cmd.Description,
				Command:     cmdStr,
				Sudo:        cmd.Sudo,
				IsUser:      cat.IsUser,
			})
		}
	}

	return result
}

func (a *App) GetCategories() []map[string]interface{} {
	var result []map[string]interface{}
	seen := make(map[string]bool)

	for _, cat := range a.categories {
		if seen[cat.Category] {
			continue
		}
		seen[cat.Category] = true
		result = append(result, map[string]interface{}{
			"name":   cat.Category,
			"icon":   cat.Icon,
			"isUser": cat.IsUser,
		})
	}
	return result
}

func (a *App) AddCategory(name, icon string) error {
	if err := os.MkdirAll(a.userDir, 0755); err != nil {
		return err
	}

	cat := CommandCategory{
		Category: name,
		Icon:     icon,
		Commands: []Command{},
		IsUser:   true,
	}

	data, err := yaml.Marshal(cat)
	if err != nil {
		return err
	}

	filename := filepath.Join(a.userDir, name+".yaml")
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return err
	}

	a.categories = append(a.categories, cat)
	return nil
}

func (a *App) AddCommand(category, name, description, command string, sudo bool) error {
	if err := os.MkdirAll(a.userDir, 0755); err != nil {
		return err
	}

	var targetCat *CommandCategory
	var isNewCat bool

	for i := range a.categories {
		if a.categories[i].Category == category && a.categories[i].IsUser {
			targetCat = &a.categories[i]
			break
		}
	}

	if targetCat == nil {
		newCat := CommandCategory{
			Category: category,
			Icon:     "📁",
			Commands: []Command{},
			IsUser:   true,
		}
		a.categories = append(a.categories, newCat)
		targetCat = &a.categories[len(a.categories)-1]
		isNewCat = true
	}

	newCmd := Command{
		Name:        name,
		Description: description,
		Command:     command,
		Sudo:        sudo,
	}

	targetCat.Commands = append(targetCat.Commands, newCmd)

	data, err := yaml.Marshal(targetCat)
	if err != nil {
		return err
	}

	filename := filepath.Join(a.userDir, targetCat.Category+".yaml")
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return err
	}

	if !isNewCat {
		for i := range a.categories {
			if a.categories[i].Category == category && a.categories[i].IsUser {
				a.categories[i] = *targetCat
				break
			}
		}
	}

	return nil
}

func (a *App) UpdateCommand(category, originalName, name, description, command string, sudo bool) error {
	for i := range a.categories {
		if a.categories[i].Category != category || !a.categories[i].IsUser {
			continue
		}

		for j := range a.categories[i].Commands {
			if a.categories[i].Commands[j].Name == originalName {
				a.categories[i].Commands[j] = Command{
					Name:        name,
					Description: description,
					Command:     command,
					Sudo:        sudo,
				}

				data, err := yaml.Marshal(a.categories[i])
				if err != nil {
					return err
				}

				filename := filepath.Join(a.userDir, category+".yaml")
				return os.WriteFile(filename, data, 0644)
			}
		}
	}
	return nil
}

func (a *App) DeleteCommand(category, name string) error {
	for i := range a.categories {
		if a.categories[i].Category != category || !a.categories[i].IsUser {
			continue
		}

		for j := range a.categories[i].Commands {
			if a.categories[i].Commands[j].Name == name {
				a.categories[i].Commands = append(
					a.categories[i].Commands[:j],
					a.categories[i].Commands[j+1:]...,
				)

				data, err := yaml.Marshal(a.categories[i])
				if err != nil {
					return err
				}

				filename := filepath.Join(a.userDir, category+".yaml")
				return os.WriteFile(filename, data, 0644)
			}
		}
	}
	return nil
}

func (a *App) DeleteCategory(category string) error {
	for i := range a.categories {
		if a.categories[i].Category == category && a.categories[i].IsUser {
			a.categories = append(a.categories[:i], a.categories[i+1:]...)
			filename := filepath.Join(a.userDir, category+".yaml")
			return os.Remove(filename)
		}
	}
	return nil
}

// Documentation methods

type DocFile struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Category string `json:"category"`
}

type DocCategory struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func (a *App) GetDocCategories() []DocCategory {
	var categories []DocCategory

	entries, err := os.ReadDir(a.docsDir)
	if err != nil {
		return categories
	}

	// Count root files
	rootCount := 0
	for _, entry := range entries {
		if !entry.IsDir() && (strings.HasSuffix(entry.Name(), ".txt") || strings.HasSuffix(entry.Name(), ".md")) {
			rootCount++
		}
	}
	if rootCount > 0 {
		categories = append(categories, DocCategory{Name: "root", Count: rootCount})
	}

	// Count files in subdirectories
	for _, entry := range entries {
		if entry.IsDir() {
			subPath := filepath.Join(a.docsDir, entry.Name())
			subEntries, err := os.ReadDir(subPath)
			if err != nil {
				continue
			}
			count := 0
			for _, subEntry := range subEntries {
				if !subEntry.IsDir() && (strings.HasSuffix(subEntry.Name(), ".txt") || strings.HasSuffix(subEntry.Name(), ".md")) {
					count++
				}
			}
			if count > 0 {
				categories = append(categories, DocCategory{Name: entry.Name(), Count: count})
			}
		}
	}

	return categories
}

func (a *App) GetDocs(category string) []DocFile {
	var docs []DocFile

	var targetDir string
	if category == "root" || category == "" {
		targetDir = a.docsDir
	} else {
		targetDir = filepath.Join(a.docsDir, category)
	}

	entries, err := os.ReadDir(targetDir)
	if err != nil {
		return docs
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".txt") || strings.HasSuffix(entry.Name(), ".md") {
			docs = append(docs, DocFile{
				Name:     entry.Name(),
				Path:     filepath.Join(targetDir, entry.Name()),
				Category: category,
			})
		}
	}

	return docs
}

func (a *App) GetDocContent(path string) (string, error) {
	// Security: ensure path is within docsDir
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(absPath, a.docsDir) {
		return "", os.ErrPermission
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (a *App) SaveDocContent(path, content string) error {
	// Security: ensure path is within docsDir
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(absPath, a.docsDir) {
		return os.ErrPermission
	}

	return os.WriteFile(absPath, []byte(content), 0644)
}

func (a *App) CreateDoc(category, name string) error {
	var targetDir string
	if category == "root" || category == "" {
		targetDir = a.docsDir
	} else {
		targetDir = filepath.Join(a.docsDir, category)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return err
		}
	}

	if !strings.HasSuffix(name, ".txt") && !strings.HasSuffix(name, ".md") {
		name = name + ".txt"
	}

	filePath := filepath.Join(targetDir, name)
	return os.WriteFile(filePath, []byte(""), 0644)
}

func (a *App) DeleteDoc(path string) error {
	// Security: ensure path is within docsDir
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(absPath, a.docsDir) {
		return os.ErrPermission
	}

	return os.Remove(absPath)
}
