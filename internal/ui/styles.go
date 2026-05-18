package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	PrimaryColor   = lipgloss.Color("#00D7FF")
	SecondaryColor = lipgloss.Color("#FF5F87")
	AccentColor    = lipgloss.Color("#87FF5F")
	MutedColor     = lipgloss.Color("#626262")
	ErrorColor     = lipgloss.Color("#FF0000")
	SuccessColor   = lipgloss.Color("#00FF00")

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(PrimaryColor).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor)

	SelectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(PrimaryColor).
			Padding(0, 1)

	NormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1)

	MutedStyle = lipgloss.NewStyle().
			Foreground(MutedColor)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryColor).
			Padding(1, 2)

	HelpStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			MarginTop(1)

	CategoryStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(AccentColor).
			MarginTop(1)

	OutputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(MutedColor).
			Padding(1, 2).
			MarginTop(1)

	TabActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(PrimaryColor).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(PrimaryColor).
			Padding(0, 2)

	TabInactiveStyle = lipgloss.NewStyle().
				Foreground(MutedColor).
				Padding(0, 2)

	StatusBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#333333")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1)
)
