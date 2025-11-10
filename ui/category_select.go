package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CategorySelectModel struct {
	categories []string
	cursor     int
	username   string
	title      string
}

func NewCategorySelectModel(username string, categories []string, title string) CategorySelectModel {
	return CategorySelectModel{
		categories: categories,
		cursor:     0,
		username:   username,
		title:      title,
	}
}

func (m CategorySelectModel) Init() tea.Cmd {
	return nil
}

func (m CategorySelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.categories)-1 {
				m.cursor++
			}
		case "enter", " ":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m CategorySelectModel) View() string {
	var b strings.Builder

	// Header
	header := HeaderStyle.Render("🔐 CYBERSEC QUIZ 🔐")
	b.WriteString(header + "\n\n")

	// Title
	title := TitleStyle.Render(m.title)
	b.WriteString(title + "\n")

	subtitle := SubtitleStyle.Render(fmt.Sprintf("Connecté en tant que: %s", m.username))
	b.WriteString(subtitle + "\n\n")

	if len(m.categories) == 0 {
		b.WriteString(ErrorStyle.Render("❌ Aucune catégorie disponible") + "\n\n")
		help := HelpStyle.Render("q: retour au menu")
		b.WriteString(help + "\n")
		return lipgloss.NewStyle().Padding(2).Render(b.String())
	}

	// Categories
	for i, cat := range m.categories {
		badge := CategoryBadgeStyle.Render(cat)
		if i == m.cursor {
			line := MenuItemSelectedStyle.Render(fmt.Sprintf("▶ %s", badge))
			b.WriteString(line + "\n")
		} else {
			line := MenuItemStyle.Render(fmt.Sprintf("  %s", badge))
			b.WriteString(line + "\n")
		}
		b.WriteString("\n")
	}

	// Help
	b.WriteString("\n")
	help := HelpStyle.Render("↑/↓ ou j/k: naviguer • enter: sélectionner • q: retour")
	b.WriteString(help + "\n")

	return lipgloss.NewStyle().Padding(2).Render(b.String())
}

func (m CategorySelectModel) GetSelectedCategory() string {
	if m.cursor < 0 || m.cursor >= len(m.categories) {
		return ""
	}
	return m.categories[m.cursor]
}
