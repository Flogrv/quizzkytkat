package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Couleurs One Dark / Dracula
	primaryColor   = lipgloss.Color("#61afef") // Bleu clair
	secondaryColor = lipgloss.Color("#c678dd") // Violet
	accentColor    = lipgloss.Color("#e06c75") // Rouge/Rose
	errorColor     = lipgloss.Color("#e06c75") // Rouge
	successColor   = lipgloss.Color("#98c379") // Vert
	warningColor   = lipgloss.Color("#e5c07b") // Jaune/Or
	textColor      = lipgloss.Color("#abb2bf") // Gris clair
	dimColor       = lipgloss.Color("#5c6370") // Gris foncé
	bgColor        = lipgloss.Color("#282c34") // Fond dark

	// Styles de base
	TitleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(1, 2).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Italic(true).
			MarginBottom(1)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(bgColor).
			Background(primaryColor).
			Bold(true).
			Padding(0, 2)

	UnselectedStyle = lipgloss.NewStyle().
			Foreground(dimColor).
			Padding(0, 2)

	MenuItemStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Padding(0, 2)

	MenuItemSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#282c34")).
				Background(primaryColor).
				Bold(true).
				Padding(0, 2).
				Width(50)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true).
			Padding(1)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true).
			Padding(1)

	QuestionStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Bold(true).
			Padding(1, 2).
			MarginBottom(1)

	AnswerStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Padding(0, 2)

	AnswerSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#282c34")).
				Background(warningColor).
				Bold(true).
				Padding(0, 2)

	AnswerCorrectStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#282c34")).
				Background(successColor).
				Bold(true).
				Padding(0, 2)

	AnswerWrongStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#abb2bf")).
				Background(errorColor).
				Bold(true).
				Padding(0, 2)

	LeaderboardHeaderStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Padding(0, 1).
				Border(lipgloss.NormalBorder(), false, false, true, false).
				BorderForeground(dimColor)

	LeaderboardRowStyle = lipgloss.NewStyle().
				Foreground(textColor).
				Padding(0, 1)

	LeaderboardTopStyle = lipgloss.NewStyle().
				Foreground(accentColor).
				Bold(true).
				Padding(0, 1)

	HelpStyle = lipgloss.NewStyle().
			Foreground(dimColor).
			Italic(true).
			MarginTop(1)

	StatsStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Italic(true).
			Padding(0, 2)

	CategoryBadgeStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#282c34")).
				Background(secondaryColor).
				Bold(true).
				Padding(0, 1).
				MarginRight(1)

	ScoreBadgeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#282c34")).
			Background(warningColor).
			Bold(true).
			Padding(0, 1)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Border(lipgloss.DoubleBorder(), false, false, true, false).
			BorderForeground(primaryColor).
			Padding(1, 2).
			Width(80).
			Align(lipgloss.Center)
)
