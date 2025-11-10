package ui

import (
	"fmt"
	"quizz-ssh/models"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LeaderboardModel struct {
	username string
	category string
	scores   []models.Score
	stats    string
	done     bool
}

func NewLeaderboardModel(username, category string, scores []models.Score, stats string) LeaderboardModel {
	return LeaderboardModel{
		username: username,
		category: category,
		scores:   scores,
		stats:    stats,
	}
}

func (m LeaderboardModel) Init() tea.Cmd {
	return nil
}

func (m LeaderboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ", "esc":
			m.done = true
			return m, nil
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m LeaderboardModel) IsDone() bool {
	return m.done
}

func (m LeaderboardModel) View() string {
	var b strings.Builder

	// Header
	header := HeaderStyle.Render("🔐 CYBERSEC QUIZ 🔐")
	b.WriteString(header + "\n\n")

	// Title
	var title string
	if m.category == "" || m.category == "global" {
		title = "🏆 Leaderboard Global"
	} else {
		title = fmt.Sprintf("📊 Leaderboard - %s", m.category)
	}

	b.WriteString(TitleStyle.Render(title) + "\n")

	subtitle := SubtitleStyle.Render(fmt.Sprintf("Connecté en tant que: %s", m.username))
	b.WriteString(subtitle + "\n\n")

	if len(m.scores) == 0 {
		b.WriteString(ErrorStyle.Render("❌ Aucun score enregistré pour l'instant") + "\n\n")
		b.WriteString(SubtitleStyle.Render("Sois le premier à jouer ! 🎮") + "\n\n")
	} else {
		// Stats
		if m.stats != "" {
			b.WriteString(StatsStyle.Render(m.stats) + "\n\n")
		}

		// Table header
		headerRow := fmt.Sprintf("%-5s %-20s %-15s %-10s", "Rank", "Pseudo", "Score", "Réussite")
		b.WriteString(LeaderboardHeaderStyle.Render(headerRow) + "\n\n")

		// Scores
		for i, score := range m.scores {
			rank := fmt.Sprintf("#%d", i+1)
			percentage := float64(score.Score) / float64(score.Total) * 100
			scoreText := fmt.Sprintf("%d/%d", score.Score, score.Total)
			successRate := fmt.Sprintf("%.1f%%", percentage)

			row := fmt.Sprintf("%-5s %-20s %-15s %-10s", rank, score.Username, scoreText, successRate)

			var style lipgloss.Style
			if i == 0 {
				// Premier place
				row = "🥇 " + row
				style = LeaderboardTopStyle
			} else if i == 1 {
				// Deuxième place
				row = "🥈 " + row
				style = LeaderboardTopStyle
			} else if i == 2 {
				// Troisième place
				row = "🥉 " + row
				style = LeaderboardTopStyle
			} else {
				row = "   " + row
				style = LeaderboardRowStyle
			}

			// Highlight current user
			if score.Username == m.username {
				style = style.Copy().Foreground(accentColor).Bold(true)
			}

			b.WriteString(style.Render(row) + "\n")
		}
	}

	b.WriteString("\n")
	help := HelpStyle.Render("q ou enter: retour au menu")
	b.WriteString(help + "\n")

	return lipgloss.NewStyle().Padding(2).Render(b.String())
}
