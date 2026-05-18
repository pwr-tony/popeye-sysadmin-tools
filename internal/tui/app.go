package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pwr-tony/popeye/internal/cmdb"
	"github.com/pwr-tony/popeye/internal/commands"
	"github.com/pwr-tony/popeye/internal/docs"
	"github.com/pwr-tony/popeye/internal/tui/views"
	"github.com/pwr-tony/popeye/internal/ui"
)

type View int

const (
	ViewHome View = iota
	ViewCommands
	ViewExecutor
	ViewCMDB
	ViewDocs
)

type Model struct {
	currentView  View
	previousView View
	width        int
	height       int
	help         help.Model
	showHelp     bool

	commandStore *commands.CommandStore
	cmdbStore    *cmdb.Store
	docStore     *docs.DocStore
	docRenderer  *docs.Renderer

	homeView     *views.HomeView
	commandsView *views.CommandsView
	executorView *views.ExecutorView
	cmdbView     *views.CMDBView
	docsView     *views.DocsView
}

func NewModel(cmdStore *commands.CommandStore, cmdbStore *cmdb.Store, docStore *docs.DocStore) Model {
	h := help.New()

	var docRenderer *docs.Renderer
	docRenderer, _ = docs.NewRenderer()

	return Model{
		currentView:  ViewHome,
		help:         h,
		commandStore: cmdStore,
		cmdbStore:    cmdbStore,
		docStore:     docStore,
		docRenderer:  docRenderer,
		homeView:     views.NewHomeView(),
		commandsView: views.NewCommandsView(cmdStore),
		cmdbView:     views.NewCMDBView(cmdbStore),
		docsView:     views.NewDocsView(docStore, docRenderer),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width

		contentHeight := msg.Height - 4
		m.homeView.SetSize(msg.Width, contentHeight)
		m.commandsView.SetSize(msg.Width, contentHeight)
		m.cmdbView.SetSize(msg.Width, contentHeight)
		m.docsView.SetSize(msg.Width, contentHeight)
		if m.executorView != nil {
			m.executorView.SetSize(msg.Width, contentHeight)
		}
		return m, nil

	case tea.KeyMsg:
		if m.showHelp {
			if msg.String() == "?" || msg.String() == "esc" {
				m.showHelp = false
				return m, nil
			}
		}

		if m.currentView == ViewExecutor {
			if key.Matches(msg, Keys.Back) {
				m.currentView = m.previousView
				m.executorView = nil
				return m, nil
			}
			var cmd tea.Cmd
			m.executorView, cmd = m.executorView.Update(msg)
			return m, cmd
		}

		if m.currentView == ViewDocs && m.docsView.IsViewing() {
			var cmd tea.Cmd
			m.docsView, cmd = m.docsView.Update(msg)
			return m, cmd
		}

		switch {
		case key.Matches(msg, Keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, Keys.Help):
			m.showHelp = !m.showHelp
			return m, nil

		case key.Matches(msg, Keys.Commands):
			m.currentView = ViewCommands
			return m, nil

		case key.Matches(msg, Keys.CMDB):
			m.currentView = ViewCMDB
			return m, nil

		case key.Matches(msg, Keys.Docs):
			m.currentView = ViewDocs
			return m, nil

		case key.Matches(msg, Keys.Back):
			if m.currentView != ViewHome {
				m.currentView = ViewHome
			}
			return m, nil

		case key.Matches(msg, Keys.Enter):
			if m.currentView == ViewCommands {
				if item := m.commandsView.SelectedCommand(); item != nil {
					m.previousView = m.currentView
					m.executorView = views.NewExecutorView(&item.Command, item.Category, item.Icon)
					m.executorView.SetSize(m.width, m.height-4)
					m.currentView = ViewExecutor

					if m.executorView.IsReady() {
						return m, m.executorView.StartExecution()
					}
				}
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	switch m.currentView {
	case ViewCommands:
		m.commandsView, cmd = m.commandsView.Update(msg)
	case ViewCMDB:
		m.cmdbView, cmd = m.cmdbView.Update(msg)
	case ViewDocs:
		m.docsView, cmd = m.docsView.Update(msg)
	case ViewExecutor:
		if m.executorView != nil {
			m.executorView, cmd = m.executorView.Update(msg)
		}
	}

	return m, cmd
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var content string

	if m.showHelp {
		content = m.renderHelp()
	} else {
		switch m.currentView {
		case ViewHome:
			content = m.homeView.View()
		case ViewCommands:
			content = m.commandsView.View()
		case ViewExecutor:
			if m.executorView != nil {
				content = m.executorView.View()
			}
		case ViewCMDB:
			content = m.cmdbView.View()
		case ViewDocs:
			content = m.docsView.View()
		}
	}

	tabs := m.renderTabs()
	statusBar := m.renderStatusBar()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		tabs,
		content,
		statusBar,
	)
}

func (m Model) renderTabs() string {
	tabs := []struct {
		name string
		view View
	}{
		{"Home", ViewHome},
		{"Commands", ViewCommands},
		{"CMDB", ViewCMDB},
		{"Docs", ViewDocs},
	}

	var rendered []string
	for _, t := range tabs {
		style := ui.TabInactiveStyle
		if m.currentView == t.view || (m.currentView == ViewExecutor && t.view == ViewCommands) {
			style = ui.TabActiveStyle
		}
		rendered = append(rendered, style.Render(t.name))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
}

func (m Model) renderStatusBar() string {
	left := "q: quit • ?: help"
	right := "Popeye v0.1.0"

	gap := m.width - lipgloss.Width(left) - lipgloss.Width(right) - 2
	if gap < 0 {
		gap = 0
	}

	bar := left + lipgloss.NewStyle().Width(gap).Render("") + right
	return ui.StatusBarStyle.Width(m.width).Render(bar)
}

func (m Model) renderHelp() string {
	return ui.BoxStyle.Render(m.help.View(Keys))
}
