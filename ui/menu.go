package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuChoice int

const (
	MenuQuiz MenuChoice = iota
	MenuLeaderboard
	MenuQuit
)

type MenuModel struct {
	choices  []string
	cursor   int
	username string
	done     bool
}

func NewMenuModel(username string) MenuModel {
	return MenuModel{
		choices: []string{
			"🎯 Jouer au Quiz",
			"🏆 Leaderboard",
			"🚪 Quitter",
		},
		cursor:   0,
		username: username,
	}
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			// Si on appuie sur q, c'est pour quitter
			m.cursor = int(MenuQuit)
			m.done = true
			return m, nil
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			m.done = true
			return m, nil
		}
	}
	return m, nil
}

func (m MenuModel) View() string {
	var b strings.Builder

	// Header
	header := HeaderStyle.Render("🔐 CYBERSEC QUIZ 🔐")
	b.WriteString(header + "\n\n")

	// User greeting
	greeting := TitleStyle.Render(fmt.Sprintf("Salut, %s ! 👋", m.username))
	b.WriteString(greeting + "\n")

	subtitle := SubtitleStyle.Render("Choisis une option pour continuer")
	b.WriteString(subtitle + "\n\n")

	// Menu items
	for i, choice := range m.choices {
		if i == m.cursor {
			b.WriteString(MenuItemSelectedStyle.Render("▶ "+choice) + "\n")
		} else {
			b.WriteString(MenuItemStyle.Render("  "+choice) + "\n")
		}
		b.WriteString("\n")
	}

	// Help
	b.WriteString("\n")
	help := HelpStyle.Render("↑/↓ ou j/k: naviguer • enter: sélectionner • q: quitter")
	b.WriteString(help + "\n")

	return lipgloss.NewStyle().Padding(2).Render(b.String())
}

func (m MenuModel) GetChoice() MenuChoice {
	return MenuChoice(m.cursor)
}

func (m MenuModel) IsDone() bool {
	return m.done
}
