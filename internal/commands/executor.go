package commands

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"text/template"
)

type ExecutionResult struct {
	Output   string
	ExitCode int
	Error    error
}

type Executor struct {
	OutputChan chan string
}

func NewExecutor() *Executor {
	return &Executor{
		OutputChan: make(chan string, 100),
	}
}

func (e *Executor) BuildCommand(cmd *Command, params map[string]string) (string, error) {
	for _, p := range cmd.Params {
		if _, exists := params[p.Name]; !exists {
			if p.Default != "" {
				params[p.Name] = p.Default
			} else if p.Required {
				return "", fmt.Errorf("missing required parameter: %s", p.Name)
			}
		}
	}

	tmpl, err := template.New("cmd").Parse(cmd.Command)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, params)
	if err != nil {
		return "", err
	}

	cmdStr := buf.String()
	if cmd.Sudo {
		cmdStr = "sudo " + cmdStr
	}

	return cmdStr, nil
}

func (e *Executor) Execute(cmdStr string) ExecutionResult {
	result := ExecutionResult{}

	cmd := exec.Command("bash", "-c", cmdStr)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		result.Error = err
		return result
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		result.Error = err
		return result
	}

	if err := cmd.Start(); err != nil {
		result.Error = err
		return result
	}

	var output strings.Builder

	go e.streamOutput(stdout, &output)
	go e.streamOutput(stderr, &output)

	err = cmd.Wait()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.ExitCode()
		} else {
			result.Error = err
		}
	}

	result.Output = output.String()
	return result
}

func (e *Executor) streamOutput(r io.Reader, output *strings.Builder) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		output.WriteString(line + "\n")
		select {
		case e.OutputChan <- line:
		default:
		}
	}
}

func (e *Executor) Close() {
	close(e.OutputChan)
}
