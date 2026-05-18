package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwr-tony/popeye/internal/cmdb"
	"github.com/pwr-tony/popeye/internal/ui"
)

type serverItem struct {
	server cmdb.Server
}

func (i serverItem) Title() string {
	return fmt.Sprintf("%s (%s)", i.server.Hostname, i.server.IP)
}

func (i serverItem) Description() string {
	return fmt.Sprintf("%s • %s", i.server.Environment, i.server.OS)
}

func (i serverItem) FilterValue() string {
	return i.server.Hostname + " " + i.server.IP
}

type CMDBView struct {
	list   list.Model
	store  *cmdb.Store
	width  int
	height int
}

func NewCMDBView(store *cmdb.Store) *CMDBView {
	items := make([]list.Item, 0)

	for _, srv := range store.GetServers() {
		items = append(items, serverItem{server: srv})
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = ui.SelectedStyle
	delegate.Styles.SelectedDesc = ui.MutedStyle

	l := list.New(items, delegate, 0, 0)
	l.Title = "CMDB - Servers"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = ui.TitleStyle

	return &CMDBView{
		list:  l,
		store: store,
	}
}

func (v *CMDBView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.list.SetSize(width-4, height-4)
}

func (v *CMDBView) Update(msg tea.Msg) (*CMDBView, tea.Cmd) {
	var cmd tea.Cmd
	v.list, cmd = v.list.Update(msg)
	return v, cmd
}

func (v *CMDBView) View() string {
	if len(v.store.GetServers()) == 0 {
		var b strings.Builder
		b.WriteString(ui.TitleStyle.Render("CMDB - Servers"))
		b.WriteString("\n\n")
		b.WriteString(ui.MutedStyle.Render("No servers in CMDB"))
		b.WriteString("\n\n")
		b.WriteString(ui.HelpStyle.Render("Add servers to data/cmdb.json"))
		return b.String()
	}
	return v.list.View()
}

func (v *CMDBView) SelectedServer() *cmdb.Server {
	if item, ok := v.list.SelectedItem().(serverItem); ok {
		return &item.server
	}
	return nil
}

func (v *CMDBView) GetServerDetails(srv *cmdb.Server) string {
	var b strings.Builder

	b.WriteString(ui.TitleStyle.Render(srv.Hostname))
	b.WriteString("\n\n")

	b.WriteString(ui.SubtitleStyle.Render("IP: "))
	b.WriteString(srv.IP)
	b.WriteString("\n")

	b.WriteString(ui.SubtitleStyle.Render("OS: "))
	b.WriteString(srv.OS)
	b.WriteString("\n")

	b.WriteString(ui.SubtitleStyle.Render("Environment: "))
	b.WriteString(srv.Environment)
	b.WriteString("\n")

	if len(srv.Tags) > 0 {
		b.WriteString(ui.SubtitleStyle.Render("Tags: "))
		b.WriteString(strings.Join(srv.Tags, ", "))
		b.WriteString("\n")
	}

	if srv.Notes != "" {
		b.WriteString("\n")
		b.WriteString(ui.SubtitleStyle.Render("Notes:"))
		b.WriteString("\n")
		b.WriteString(srv.Notes)
	}

	return b.String()
}

func (v *CMDBView) Refresh() {
	items := make([]list.Item, 0)
	for _, srv := range v.store.GetServers() {
		items = append(items, serverItem{server: srv})
	}
	v.list.SetItems(items)
}
