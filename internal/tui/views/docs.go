package views

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwr-tony/popeye/internal/docs"
	"github.com/pwr-tony/popeye/internal/ui"
)

type docItem struct {
	doc docs.Document
}

func (i docItem) Title() string {
	return i.doc.Name
}

func (i docItem) Description() string {
	lines := strings.Split(i.doc.Content, "\n")
	if len(lines) > 0 {
		first := strings.TrimPrefix(lines[0], "# ")
		if len(first) > 50 {
			return first[:50] + "..."
		}
		return first
	}
	return ""
}

func (i docItem) FilterValue() string {
	return i.doc.Name
}

type DocsView struct {
	list       list.Model
	viewport   viewport.Model
	store      *docs.DocStore
	renderer   *docs.Renderer
	width      int
	height     int
	viewing    bool
	currentDoc *docs.Document
}

func NewDocsView(store *docs.DocStore, renderer *docs.Renderer) *DocsView {
	items := make([]list.Item, 0)

	for _, doc := range store.GetDocs() {
		items = append(items, docItem{doc: doc})
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = ui.SelectedStyle
	delegate.Styles.SelectedDesc = ui.MutedStyle

	l := list.New(items, delegate, 0, 0)
	l.Title = "Documentation"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = ui.TitleStyle

	vp := viewport.New(0, 0)

	return &DocsView{
		list:     l,
		viewport: vp,
		store:    store,
		renderer: renderer,
	}
}

func (v *DocsView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.list.SetSize(width-4, height-4)
	v.viewport.Width = width - 4
	v.viewport.Height = height - 6
}

func (v *DocsView) Update(msg tea.Msg) (*DocsView, tea.Cmd) {
	var cmd tea.Cmd

	if v.viewing {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc", "q":
				v.viewing = false
				v.currentDoc = nil
				return v, nil
			}
		}
		v.viewport, cmd = v.viewport.Update(msg)
		return v, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if item, ok := v.list.SelectedItem().(docItem); ok {
				v.currentDoc = &item.doc
				v.viewing = true
				if v.renderer != nil {
					rendered, err := v.renderer.Render(item.doc.Content)
					if err == nil {
						v.viewport.SetContent(rendered)
					} else {
						v.viewport.SetContent(item.doc.Content)
					}
				} else {
					v.viewport.SetContent(item.doc.Content)
				}
				return v, nil
			}
		}
	}

	v.list, cmd = v.list.Update(msg)
	return v, cmd
}

func (v *DocsView) View() string {
	if len(v.store.GetDocs()) == 0 && !v.viewing {
		var b strings.Builder
		b.WriteString(ui.TitleStyle.Render("Documentation"))
		b.WriteString("\n\n")
		b.WriteString(ui.MutedStyle.Render("No documents found"))
		b.WriteString("\n\n")
		b.WriteString(ui.HelpStyle.Render("Add .md files to docs/procedures/"))
		return b.String()
	}

	if v.viewing {
		var b strings.Builder
		b.WriteString(ui.TitleStyle.Render(v.currentDoc.Name))
		b.WriteString("\n")
		b.WriteString(v.viewport.View())
		b.WriteString("\n")
		b.WriteString(ui.HelpStyle.Render("↑/↓: scroll • Esc: back"))
		return b.String()
	}

	return v.list.View()
}

func (v *DocsView) IsViewing() bool {
	return v.viewing
}

func (v *DocsView) Refresh() {
	items := make([]list.Item, 0)
	for _, doc := range v.store.GetDocs() {
		items = append(items, docItem{doc: doc})
	}
	v.list.SetItems(items)
}
