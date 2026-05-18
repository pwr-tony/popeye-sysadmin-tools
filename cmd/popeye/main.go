package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwr-tony/popeye/internal/cmdb"
	"github.com/pwr-tony/popeye/internal/commands"
	"github.com/pwr-tony/popeye/internal/docs"
	"github.com/pwr-tony/popeye/internal/tui"
)

func main() {
	baseDir := getBaseDir()

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

func getBaseDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}

	exeDir := filepath.Dir(exe)

	if _, err := os.Stat(filepath.Join(exeDir, "configs")); err == nil {
		return exeDir
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}

	return cwd
}
