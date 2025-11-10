package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type UsernameModel struct {
	textInput textinput.Model
	err       error
	username  string
	done      bool // Indique si l'utilisateur a validé son pseudo
}

func NewUsernameModel() UsernameModel {
	ti := textinput.New()
	ti.Placeholder = "Ton pseudo..."
	ti.Focus()
	ti.CharLimit = 20
	ti.Width = 30

	return UsernameModel{
		textInput: ti,
	}
}

func (m UsernameModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m UsernameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			username := strings.TrimSpace(m.textInput.Value())
			if len(username) >= 3 {
				m.username = username
				m.done = true
				return m, nil // Ne pas quitter, juste marquer comme done
			}
			m.err = fmt.Errorf("le pseudo doit faire au moins 3 caractères")
			return m, nil
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit // Vraiment quitter uniquement si Ctrl+C
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m UsernameModel) View() string {
	var b strings.Builder

	// Header avec ASCII art
	header := HeaderStyle.Render("🔐 CYBERSEC QUIZ 🔐")
	b.WriteString(header + "\n\n")

	// Welcome message
	welcome := TitleStyle.Render("Bienvenue sur le quiz de cybersécurité !")
	b.WriteString(welcome + "\n")

	subtitle := SubtitleStyle.Render("Entre ton pseudo pour commencer")
	b.WriteString(subtitle + "\n\n")

	// Input box
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(1, 2).
		Width(50).
		Align(lipgloss.Center)

	b.WriteString(box.Render(m.textInput.View()) + "\n\n")

	// Error message
	if m.err != nil {
		b.WriteString(ErrorStyle.Render(fmt.Sprintf("❌ %s", m.err.Error())) + "\n\n")
	}

	// Help text
	help := HelpStyle.Render("enter: valider • ctrl+c: quitter")
	b.WriteString(help + "\n")

	return b.String()
}

func (m UsernameModel) GetUsername() string {
	return m.username
}

func (m UsernameModel) IsDone() bool {
	return m.done
}
