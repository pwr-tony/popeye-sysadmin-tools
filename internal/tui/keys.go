package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Quit       key.Binding
	Help       key.Binding
	Commands   key.Binding
	CMDB       key.Binding
	Docs       key.Binding
	Enter      key.Binding
	Back       key.Binding
	Search     key.Binding
	Tab        key.Binding
	Up         key.Binding
	Down       key.Binding
	Left       key.Binding
	Right      key.Binding
	PageUp     key.Binding
	PageDown   key.Binding
}

var Keys = KeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Commands: key.NewBinding(
		key.WithKeys("1"),
		key.WithHelp("1", "commands"),
	),
	CMDB: key.NewBinding(
		key.WithKeys("2"),
		key.WithHelp("2", "cmdb"),
	),
	Docs: key.NewBinding(
		key.WithKeys("3"),
		key.WithHelp("3", "docs"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "search"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch panel"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "right"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("pgup", "page up"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("pgdown", "page down"),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Help, k.Commands, k.CMDB, k.Docs}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Enter, k.Back, k.Search, k.Tab},
		{k.Commands, k.CMDB, k.Docs},
		{k.Quit, k.Help},
	}
}
