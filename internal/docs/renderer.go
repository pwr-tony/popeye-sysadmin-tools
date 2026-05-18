package docs

import (
	"github.com/charmbracelet/glamour"
)

type Renderer struct {
	renderer *glamour.TermRenderer
}

func NewRenderer() (*Renderer, error) {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		return nil, err
	}

	return &Renderer{renderer: r}, nil
}

func (r *Renderer) Render(content string) (string, error) {
	return r.renderer.Render(content)
}
