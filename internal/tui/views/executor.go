package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwr-tony/popeye/internal/commands"
	"github.com/pwr-tony/popeye/internal/ui"
)

type ExecutorState int

const (
	StateParams ExecutorState = iota
	StateRunning
	StateDone
)

type OutputMsg string
type DoneMsg struct {
	ExitCode int
	Error    error
}

type ExecutorView struct {
	command      *commands.Command
	category     string
	icon         string
	executor     *commands.Executor
	params       map[string]string
	paramInputs  []textinput.Model
	currentParam int
	output       textarea.Model
	state        ExecutorState
	width        int
	height       int
	exitCode     int
	execError    error
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

	ta := textarea.New()
	ta.SetWidth(80)
	ta.SetHeight(20)
	ta.ShowLineNumbers = false

	state := StateParams
	if len(cmd.Params) == 0 {
		state = StateRunning
	}

	return &ExecutorView{
		command:      cmd,
		category:     category,
		icon:         icon,
		executor:     commands.NewExecutor(),
		params:       params,
		paramInputs:  inputs,
		currentParam: 0,
		output:       ta,
		state:        state,
	}
}

func (v *ExecutorView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.output.SetWidth(width - 10)
	v.output.SetHeight(height - 15)
}

func (v *ExecutorView) Update(msg tea.Msg) (*ExecutorView, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if v.state == StateParams {
			switch msg.String() {
			case "tab", "down", "enter":
				if v.currentParam < len(v.paramInputs)-1 {
					v.paramInputs[v.currentParam].Blur()
					v.currentParam++
					v.paramInputs[v.currentParam].Focus()
				} else if msg.String() == "enter" {
					return v, v.executeCommand()
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

	case OutputMsg:
		current := v.output.Value()
		v.output.SetValue(current + string(msg) + "\n")
		return v, v.waitForOutput()

	case DoneMsg:
		v.state = StateDone
		v.exitCode = msg.ExitCode
		v.execError = msg.Error
		return v, nil
	}

	return v, nil
}

func (v *ExecutorView) executeCommand() tea.Cmd {
	for i, p := range v.command.Params {
		v.params[p.Name] = v.paramInputs[i].Value()
	}

	v.state = StateRunning

	return func() tea.Msg {
		cmdStr, err := v.executor.BuildCommand(v.command, v.params)
		if err != nil {
			return DoneMsg{Error: err}
		}

		result := v.executor.Execute(cmdStr)
		return DoneMsg{
			ExitCode: result.ExitCode,
			Error:    result.Error,
		}
	}
}

func (v *ExecutorView) waitForOutput() tea.Cmd {
	return func() tea.Msg {
		line, ok := <-v.executor.OutputChan
		if !ok {
			return nil
		}
		return OutputMsg(line)
	}
}

func (v *ExecutorView) View() string {
	var b strings.Builder

	b.WriteString(ui.TitleStyle.Render(fmt.Sprintf("%s %s/%s", v.icon, v.category, v.command.Name)))
	b.WriteString("\n\n")

	switch v.state {
	case StateParams:
		b.WriteString(ui.SubtitleStyle.Render("Enter parameters:"))
		b.WriteString("\n\n")

		for i, p := range v.command.Params {
			label := p.Name
			if p.Required {
				label += " *"
			}
			b.WriteString(fmt.Sprintf("%s: ", label))
			b.WriteString(v.paramInputs[i].View())
			b.WriteString("\n")
		}

		b.WriteString("\n")
		b.WriteString(ui.HelpStyle.Render("Tab/Enter: next • Shift+Tab: prev • Enter on last: execute"))

	case StateRunning:
		b.WriteString(ui.SubtitleStyle.Render("Running..."))
		b.WriteString("\n\n")
		b.WriteString(ui.OutputStyle.Render(v.output.Value()))

	case StateDone:
		if v.execError != nil {
			b.WriteString(ui.ErrorStyle.Render(fmt.Sprintf("Error: %v", v.execError)))
		} else if v.exitCode != 0 {
			b.WriteString(ui.ErrorStyle.Render(fmt.Sprintf("Exit code: %d", v.exitCode)))
		} else {
			b.WriteString(ui.SuccessStyle.Render("Command completed successfully"))
		}
		b.WriteString("\n\n")
		b.WriteString(ui.OutputStyle.Render(v.output.Value()))
		b.WriteString("\n")
		b.WriteString(ui.HelpStyle.Render("Press Esc to go back"))
	}

	return b.String()
}

func (v *ExecutorView) IsReady() bool {
	return v.state == StateRunning && len(v.command.Params) == 0
}

func (v *ExecutorView) StartExecution() tea.Cmd {
	return v.executeCommand()
}
