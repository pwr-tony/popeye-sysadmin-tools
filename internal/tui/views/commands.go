package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwr-tony/popeye/internal/commands"
	"github.com/pwr-tony/popeye/internal/ui"
)

type commandItem struct {
	Category string
	Icon     string
	Command  commands.Command
}

func (i commandItem) Title() string {
	return fmt.Sprintf("%s %s/%s", i.Icon, i.Category, i.Command.Name)
}

func (i commandItem) Description() string {
	return i.Command.Description
}

func (i commandItem) FilterValue() string {
	return i.Category + "/" + i.Command.Name
}

type CommandsView struct {
	list     list.Model
	store    *commands.CommandStore
	width    int
	height   int
	selected *commandItem
}

func NewCommandsView(store *commands.CommandStore) *CommandsView {
	items := make([]list.Item, 0)

	for _, cat := range store.Categories {
		for _, cmd := range cat.Commands {
			items = append(items, commandItem{
				Category: cat.Category,
				Icon:     cat.Icon,
				Command:  cmd,
			})
		}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = ui.SelectedStyle
	delegate.Styles.SelectedDesc = ui.MutedStyle

	l := list.New(items, delegate, 0, 0)
	l.Title = "Commands"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = ui.TitleStyle

	return &CommandsView{
		list:  l,
		store: store,
	}
}

func (v *CommandsView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.list.SetSize(width-4, height-4)
}

func (v *CommandsView) Update(msg tea.Msg) (*CommandsView, tea.Cmd) {
	var cmd tea.Cmd
	v.list, cmd = v.list.Update(msg)
	return v, cmd
}

func (v *CommandsView) View() string {
	return v.list.View()
}

func (v *CommandsView) SelectedCommand() *commandItem {
	if item, ok := v.list.SelectedItem().(commandItem); ok {
		return &item
	}
	return nil
}

func (v *CommandsView) GetCommandDetails(item *commandItem) string {
	var b strings.Builder

	b.WriteString(ui.TitleStyle.Render(fmt.Sprintf("%s %s", item.Icon, item.Command.Name)))
	b.WriteString("\n\n")
	b.WriteString(ui.SubtitleStyle.Render("Category: "))
	b.WriteString(item.Category)
	b.WriteString("\n\n")
	b.WriteString(ui.SubtitleStyle.Render("Description:"))
	b.WriteString("\n")
	b.WriteString(item.Command.Description)
	b.WriteString("\n\n")

	if len(item.Command.Params) > 0 {
		b.WriteString(ui.SubtitleStyle.Render("Parameters:"))
		b.WriteString("\n")
		for _, p := range item.Command.Params {
			req := ""
			if p.Required {
				req = " (required)"
			}
			def := ""
			if p.Default != "" {
				def = fmt.Sprintf(" [default: %s]", p.Default)
			}
			b.WriteString(fmt.Sprintf("  • %s: %s%s%s\n", p.Name, p.Type, req, def))
			if p.Prompt != "" {
				b.WriteString(fmt.Sprintf("    %s\n", ui.MutedStyle.Render(p.Prompt)))
			}
		}
		b.WriteString("\n")
	}

	b.WriteString(ui.SubtitleStyle.Render("Command:"))
	b.WriteString("\n")
	b.WriteString(ui.MutedStyle.Render(item.Command.Command))

	if item.Command.Sudo {
		b.WriteString("\n\n")
		b.WriteString(ui.ErrorStyle.Render("⚠ Requires sudo"))
	}

	return b.String()
}
