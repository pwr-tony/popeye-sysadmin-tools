package views

import (
	"fmt"
	"strings"

	"github.com/pwr-tony/popeye/internal/ui"
)

type HomeView struct {
	width  int
	height int
}

func NewHomeView() *HomeView {
	return &HomeView{}
}

func (v *HomeView) SetSize(width, height int) {
	v.width = width
	v.height = height
}

func (v *HomeView) View() string {
	var b strings.Builder

	logo := `
    ____
   / __ \____  ____  ___  __  _____
  / /_/ / __ \/ __ \/ _ \/ / / / _ \
 / ____/ /_/ / /_/ /  __/ /_/ /  __/
/_/    \____/ .___/\___/\__, /\___/
           /_/         /____/
`

	b.WriteString(ui.TitleStyle.Render(logo))
	b.WriteString("\n\n")
	b.WriteString(ui.SubtitleStyle.Render("Sysadmin Tools Manager"))
	b.WriteString("\n\n")

	menu := []string{
		"[1] Commands  - Execute sysadmin commands",
		"[2] CMDB      - Manage infrastructure inventory",
		"[3] Docs      - View procedures and documentation",
	}

	for _, item := range menu {
		b.WriteString(ui.NormalStyle.Render(item))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(ui.HelpStyle.Render("Press number to navigate • ? for help • q to quit"))

	stats := fmt.Sprintf("\n\nv0.1.0 • %d commands available", 0)
	b.WriteString(ui.MutedStyle.Render(stats))

	return ui.BoxStyle.Render(b.String())
}
