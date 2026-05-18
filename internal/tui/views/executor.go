package views

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwr-tony/popeye/internal/commands"
	"github.com/pwr-tony/popeye/internal/ui"
)

type ExecutorState int

const (
	StateParams ExecutorState = iota
	StateShow
)

type ExecutorView struct {
	command      *commands.Command
	category     string
	icon         string
	params       map[string]string
	paramInputs  []textinput.Model
	currentParam int
	state        ExecutorState
	width        int
	height       int
	finalCmd     string
	copied       bool
}

func NewExecutorView(cmd *commands.Command, category, icon string) *ExecutorView {
	params := make(map[string]string)
	inputs := make([]textinput.Model, len(cmd.Params))

	for i, p := range cmd.Params {
		ti := textinput.New()
		ti.Placeholder = p.Prompt
		if p.Prompt == "" {
			ti.Placeholder = p.Name
		}
		ti.CharLimit = 256
		if p.Default != "" {
			ti.SetValue(p.Default)
		}
		if i == 0 {
			ti.Focus()
		}
		inputs[i] = ti
	}

	state := StateParams
	if len(cmd.Params) == 0 {
		state = StateShow
	}

	return &ExecutorView{
		command:      cmd,
		category:     category,
		icon:         icon,
		params:       params,
		paramInputs:  inputs,
		currentParam: 0,
		state:        state,
	}
}

func (v *ExecutorView) SetSize(width, height int) {
	v.width = width
	v.height = height
}

func (v *ExecutorView) Update(msg tea.Msg) (*ExecutorView, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if v.state == StateParams {
			switch msg.String() {
			case "tab", "down":
				if v.currentParam < len(v.paramInputs)-1 {
					v.paramInputs[v.currentParam].Blur()
					v.currentParam++
					v.paramInputs[v.currentParam].Focus()
				}
				return v, nil
			case "enter":
				if v.currentParam < len(v.paramInputs)-1 {
					v.paramInputs[v.currentParam].Blur()
					v.currentParam++
					v.paramInputs[v.currentParam].Focus()
				} else {
					v.buildCommand()
					v.state = StateShow
				}
				return v, nil
			case "shift+tab", "up":
				if v.currentParam > 0 {
					v.paramInputs[v.currentParam].Blur()
					v.currentParam--
					v.paramInputs[v.currentParam].Focus()
				}
				return v, nil
			}

			var cmd tea.Cmd
			v.paramInputs[v.currentParam], cmd = v.paramInputs[v.currentParam].Update(msg)
			return v, cmd
		}

		if v.state == StateShow {
			switch msg.String() {
			case "c":
				if err := clipboard.WriteAll(v.finalCmd); err == nil {
					v.copied = true
				}
				return v, nil
			}
		}
	}

	return v, nil
}

func (v *ExecutorView) buildCommand() {
	for i, p := range v.command.Params {
		val := v.paramInputs[i].Value()
		if val == "" && p.Default != "" {
			val = p.Default
		}
		v.params[p.Name] = val
	}

	cmdStr := v.command.Command
	for name, value := range v.params {
		cmdStr = strings.ReplaceAll(cmdStr, "{{."+name+"}}", value)
	}

	cmdStr = strings.TrimSpace(cmdStr)

	if v.command.Sudo {
		cmdStr = "sudo " + cmdStr
	}

	v.finalCmd = cmdStr
}

func (v *ExecutorView) View() string {
	var b strings.Builder

	b.WriteString(ui.TitleStyle.Render(fmt.Sprintf("%s %s/%s", v.icon, v.category, v.command.Name)))
	b.WriteString("\n\n")
	b.WriteString(ui.MutedStyle.Render(v.command.Description))
	b.WriteString("\n\n")

	switch v.state {
	case StateParams:
		b.WriteString(ui.SubtitleStyle.Render("Parameters:"))
		b.WriteString("\n\n")

		for i, p := range v.command.Params {
			label := p.Name
			if p.Required {
				label += " *"
			}
			if p.Default != "" {
				label += fmt.Sprintf(" [%s]", p.Default)
			}
			b.WriteString(fmt.Sprintf("  %s: ", label))
			b.WriteString(v.paramInputs[i].View())
			b.WriteString("\n")
		}

		b.WriteString("\n")
		b.WriteString(ui.HelpStyle.Render("Tab/Enter: next field • Enter on last: show command • Esc: back"))

	case StateShow:
		if v.finalCmd == "" {
			v.buildCommand()
		}

		b.WriteString(ui.SubtitleStyle.Render("Command:"))
		b.WriteString("\n\n")
		b.WriteString(ui.OutputStyle.Render(v.finalCmd))
		b.WriteString("\n\n")

		if v.command.Sudo {
			b.WriteString(ui.ErrorStyle.Render("⚠ Requires sudo"))
			b.WriteString("\n\n")
		}

		if v.copied {
			b.WriteString(ui.SuccessStyle.Render("✓ Copied to clipboard!"))
		} else {
			b.WriteString(ui.HelpStyle.Render("c: copy to clipboard • Esc: back"))
		}
	}

	return b.String()
}

func (v *ExecutorView) IsReady() bool {
	return v.state == StateShow && len(v.command.Params) == 0
}

func (v *ExecutorView) StartExecution() tea.Cmd {
	return nil
}
