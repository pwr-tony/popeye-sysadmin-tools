package commands

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type CommandStore struct {
	Categories []CommandCategory
}

func NewCommandStore() *CommandStore {
	return &CommandStore{
		Categories: make([]CommandCategory, 0),
	}
}

func (s *CommandStore) LoadFromDir(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return err
	}

	for _, file := range files {
		cat, err := s.loadFile(file)
		if err != nil {
			return err
		}
		s.Categories = append(s.Categories, cat)
	}

	return nil
}

func (s *CommandStore) loadFile(path string) (CommandCategory, error) {
	var cat CommandCategory

	data, err := os.ReadFile(path)
	if err != nil {
		return cat, err
	}

	err = yaml.Unmarshal(data, &cat)
	return cat, err
}

func (s *CommandStore) GetCommand(category, name string) *Command {
	for _, cat := range s.Categories {
		if cat.Category == category {
			for i := range cat.Commands {
				if cat.Commands[i].Name == name {
					return &cat.Commands[i]
				}
			}
		}
	}
	return nil
}

func (s *CommandStore) AllCommands() []struct {
	Category string
	Icon     string
	Command  Command
} {
	var result []struct {
		Category string
		Icon     string
		Command  Command
	}

	for _, cat := range s.Categories {
		for _, cmd := range cat.Commands {
			result = append(result, struct {
				Category string
				Icon     string
				Command  Command
			}{
				Category: cat.Category,
				Icon:     cat.Icon,
				Command:  cmd,
			})
		}
	}

	return result
}
