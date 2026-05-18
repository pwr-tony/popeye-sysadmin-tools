package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwr-tony/popeye/internal/cmdb"
	"github.com/pwr-tony/popeye/internal/commands"
	"github.com/pwr-tony/popeye/internal/docs"
	"github.com/pwr-tony/popeye/internal/tui"
	"gopkg.in/yaml.v3"
)

var baseDir string

func main() {
	baseDir = getBaseDir()

	if len(os.Args) < 2 {
		runTUI()
		return
	}

	switch os.Args[1] {
	case "list":
		cmdList()
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: popeye add <category|command>")
			os.Exit(1)
		}
		switch os.Args[2] {
		case "category":
			cmdAddCategory()
		case "command":
			cmdAddCommand()
		default:
			fmt.Printf("Unknown: %s\n", os.Args[2])
			os.Exit(1)
		}
	case "help", "-h", "--help":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println(`Popeye - Sysadmin Tools Manager

Usage:
  popeye              Launch TUI
  popeye list         List all commands
  popeye add category Add new category
  popeye add command  Add new command
  popeye help         Show this help`)
}

func runTUI() {
	cmdStore := commands.NewCommandStore()
	if err := cmdStore.LoadFromDir(filepath.Join(baseDir, "configs", "commands")); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not load commands: %v\n", err)
	}

	cmdbStore := cmdb.NewStore(filepath.Join(baseDir, "data", "cmdb.json"))
	if err := cmdbStore.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not load CMDB: %v\n", err)
	}

	docStore := docs.NewDocStore()
	if err := docStore.LoadFromDir(filepath.Join(baseDir, "docs", "procedures")); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not load docs: %v\n", err)
	}

	model := tui.NewModel(cmdStore, cmdbStore, docStore)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func cmdList() {
	cmdStore := commands.NewCommandStore()
	if err := cmdStore.LoadFromDir(filepath.Join(baseDir, "configs", "commands")); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	for _, cat := range cmdStore.Categories {
		fmt.Printf("\n%s %s\n", cat.Icon, strings.ToUpper(cat.Category))
		fmt.Println(strings.Repeat("-", 40))
		for _, cmd := range cat.Commands {
			fmt.Printf("  %-25s %s\n", cmd.Name, cmd.Description)
		}
	}
	fmt.Println()
}

func cmdAddCategory() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Category name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Icon (emoji): ")
	icon, _ := reader.ReadString('\n')
	icon = strings.TrimSpace(icon)

	if name == "" {
		fmt.Println("Category name is required")
		os.Exit(1)
	}
	if icon == "" {
		icon = "📁"
	}

	cat := commands.CommandCategory{
		Category: name,
		Icon:     icon,
		Commands: []commands.Command{},
	}

	data, err := yaml.Marshal(cat)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	filename := filepath.Join(baseDir, "configs", "commands", name+".yaml")
	if err := os.WriteFile(filename, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created category: %s %s\n", icon, name)
	fmt.Printf("File: %s\n", filename)
}

func cmdAddCommand() {
	cmdStore := commands.NewCommandStore()
	if err := cmdStore.LoadFromDir(filepath.Join(baseDir, "configs", "commands")); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(cmdStore.Categories) == 0 {
		fmt.Println("No categories found. Create one first: popeye add category")
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Available categories:")
	for i, cat := range cmdStore.Categories {
		fmt.Printf("  [%d] %s %s\n", i+1, cat.Icon, cat.Category)
	}
	fmt.Print("\nSelect category (number): ")
	var catIdx int
	fmt.Scanln(&catIdx)
	catIdx--

	if catIdx < 0 || catIdx >= len(cmdStore.Categories) {
		fmt.Println("Invalid selection")
		os.Exit(1)
	}

	selectedCat := cmdStore.Categories[catIdx]

	fmt.Print("Command name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Description: ")
	desc, _ := reader.ReadString('\n')
	desc = strings.TrimSpace(desc)

	fmt.Print("Command (shell): ")
	cmdStr, _ := reader.ReadString('\n')
	cmdStr = strings.TrimSpace(cmdStr)

	fmt.Print("Requires sudo? (y/N): ")
	sudoStr, _ := reader.ReadString('\n')
	sudo := strings.ToLower(strings.TrimSpace(sudoStr)) == "y"

	newCmd := commands.Command{
		Name:        name,
		Description: desc,
		Command:     cmdStr,
		Sudo:        sudo,
		Params:      []commands.Param{},
	}

	fmt.Print("Add parameters? (y/N): ")
	addParams, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(addParams)) == "y" {
		for {
			fmt.Print("\nParameter name (empty to finish): ")
			pName, _ := reader.ReadString('\n')
			pName = strings.TrimSpace(pName)
			if pName == "" {
				break
			}

			fmt.Print("Required? (y/N): ")
			reqStr, _ := reader.ReadString('\n')
			required := strings.ToLower(strings.TrimSpace(reqStr)) == "y"

			fmt.Print("Default value (empty for none): ")
			defVal, _ := reader.ReadString('\n')
			defVal = strings.TrimSpace(defVal)

			fmt.Print("Prompt text: ")
			prompt, _ := reader.ReadString('\n')
			prompt = strings.TrimSpace(prompt)

			newCmd.Params = append(newCmd.Params, commands.Param{
				Name:     pName,
				Type:     "string",
				Required: required,
				Default:  defVal,
				Prompt:   prompt,
			})
		}
	}

	selectedCat.Commands = append(selectedCat.Commands, newCmd)

	data, err := yaml.Marshal(selectedCat)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	filename := filepath.Join(baseDir, "configs", "commands", selectedCat.Category+".yaml")
	if err := os.WriteFile(filename, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nAdded command: %s/%s\n", selectedCat.Category, name)
	fmt.Printf("File updated: %s\n", filename)
}

func getBaseDir() string {
	// Try executable directory
	exe, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exe)
		if _, err := os.Stat(filepath.Join(exeDir, "configs")); err == nil {
			return exeDir
		}
		// Try parent of executable (for bin/)
		parentDir := filepath.Dir(exeDir)
		if _, err := os.Stat(filepath.Join(parentDir, "configs")); err == nil {
			return parentDir
		}
	}

	// Try current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	if _, err := os.Stat(filepath.Join(cwd, "configs")); err == nil {
		return cwd
	}

	// Try parent of cwd
	parentCwd := filepath.Dir(cwd)
	if _, err := os.Stat(filepath.Join(parentCwd, "configs")); err == nil {
		return parentCwd
	}

	return cwd
}
