package ui

import (
	"fmt"
	"math/rand"
	"quizz-ssh/models"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type QuizState int

const (
	QuizStateQuestion QuizState = iota
	QuizStateResult
	QuizStateFinished
)

type QuizModel struct {
	username      string
	questions     []models.Question
	currentIndex  int
	cursor        int
	score         int
	state         QuizState
	userAnswer    int
	correctAnswer int
	category      string
	showResult    bool
	resultTime    time.Time
}

func NewQuizModel(username string, questions []models.Question, category string) QuizModel {
	// Shuffle les réponses de chaque question
	shuffledQuestions := make([]models.Question, len(questions))
	for i, q := range questions {
		shuffledQuestions[i] = shuffleQuestion(q)
	}

	return QuizModel{
		username:     username,
		questions:    shuffledQuestions,
		currentIndex: 0,
		cursor:       0,
		score:        0,
		state:        QuizStateQuestion,
		category:     category,
	}
}

// shuffleQuestion mélange les options d'une question et track la bonne réponse
func shuffleQuestion(q models.Question) models.Question {
	// Créer une copie de la question
	shuffled := q

	// Si pas d'options, retourner tel quel
	if len(q.Options) == 0 {
		return shuffled
	}

	// Créer un générateur aléatoire avec seed unique
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Créer un slice avec les indices
	indices := make([]int, len(q.Options))
	for i := range indices {
		indices[i] = i
	}

	// Mélanger les indices avec Fisher-Yates
	rng.Shuffle(len(indices), func(i, j int) {
		indices[i], indices[j] = indices[j], indices[i]
	})

	// Créer les nouvelles options mélangées
	shuffled.ShuffledOptions = make([]string, len(q.Options))
	for newIdx, oldIdx := range indices {
		shuffled.ShuffledOptions[newIdx] = q.Options[oldIdx]
		// Si c'était la bonne réponse, on note sa nouvelle position
		if oldIdx == q.Answer {
			shuffled.ShuffledAnswer = newIdx
		}
	}

	return shuffled
}

func (m QuizModel) Init() tea.Cmd {
	return nil
}

type resultTimeoutMsg struct{}

func waitForResult() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return resultTimeoutMsg{}
	})
}

func (m QuizModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case QuizStateQuestion:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.questions[m.currentIndex].ShuffledOptions)-1 {
					m.cursor++
				}
			case "enter", " ":
				m.userAnswer = m.cursor
				m.correctAnswer = m.questions[m.currentIndex].ShuffledAnswer
				if m.userAnswer == m.correctAnswer {
					m.score++
				}
				m.state = QuizStateResult
				m.showResult = true
				return m, nil
			}

		case QuizStateResult:
			// Appuyer sur Enter pour passer à la question suivante
			switch msg.String() {
			case "enter", " ":
				m.currentIndex++
				if m.currentIndex >= len(m.questions) {
					m.state = QuizStateFinished
				} else {
					m.state = QuizStateQuestion
					m.cursor = 0
					m.showResult = false
				}
				return m, nil
			case "ctrl+c", "q":
				return m, tea.Quit
			}

		case QuizStateFinished:
			switch msg.String() {
			case "enter", " ", "q", "ctrl+c":
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m QuizModel) View() string {
	var b strings.Builder

	// Header
	header := HeaderStyle.Render("🔐 CYBERSEC QUIZ 🔐")
	b.WriteString(header + "\n\n")

	if len(m.questions) == 0 {
		b.WriteString(ErrorStyle.Render("❌ Aucune question disponible") + "\n\n")
		help := HelpStyle.Render("q: retour au menu")
		b.WriteString(help + "\n")
		return lipgloss.NewStyle().Padding(2).Render(b.String())
	}

	switch m.state {
	case QuizStateQuestion, QuizStateResult:
		m.renderQuestion(&b)
	case QuizStateFinished:
		m.renderFinished(&b)
	}

	return lipgloss.NewStyle().Padding(2).Render(b.String())
}

func (m QuizModel) renderQuestion(b *strings.Builder) {
	question := m.questions[m.currentIndex]

	// Progress bar
	progress := fmt.Sprintf("Question %d/%d", m.currentIndex+1, len(m.questions))
	progressBar := m.renderProgressBar()

	catBadge := CategoryBadgeStyle.Render(question.Category)
	scoreBadge := ScoreBadgeStyle.Render(fmt.Sprintf("Score: %d/%d", m.score, m.currentIndex))

	info := lipgloss.JoinHorizontal(lipgloss.Left, catBadge, " ", scoreBadge, "  ", StatsStyle.Render(progress))
	b.WriteString(info + "\n")
	b.WriteString(progressBar + "\n\n")

	// Question
	questionBox := BoxStyle.Render(QuestionStyle.Render("❓ " + question.Text))
	b.WriteString(questionBox + "\n\n")

	// Options (utiliser les options shufflées)
	options := question.ShuffledOptions
	if len(options) == 0 {
		options = question.Options // Fallback si pas shufflé
	}

	for i, option := range options {
		var line string
		prefix := fmt.Sprintf("%c) ", 'A'+i)

		if m.state == QuizStateResult {
			// Afficher le résultat
			if i == m.correctAnswer {
				line = AnswerCorrectStyle.Render(prefix + option + " ✓")
			} else if i == m.userAnswer {
				line = AnswerWrongStyle.Render(prefix + option + " ✗")
			} else {
				line = AnswerStyle.Render(prefix + option)
			}
		} else {
			// Mode sélection
			if i == m.cursor {
				line = AnswerSelectedStyle.Render("▶ " + prefix + option)
			} else {
				line = AnswerStyle.Render("  " + prefix + option)
			}
		}
		b.WriteString(line + "\n")
	}

	b.WriteString("\n")

	// Result message
	if m.state == QuizStateResult {
		if m.userAnswer == m.correctAnswer {
			msg := SuccessStyle.Render("🎉 Bonne réponse !")
			b.WriteString(msg + "\n\n")
		} else {
			msg := ErrorStyle.Render("❌ Mauvaise réponse !")
			b.WriteString(msg + "\n\n")
		}
		help := HelpStyle.Render("enter: question suivante • q: quitter")
		b.WriteString(help + "\n")
	} else {
		// Help
		help := HelpStyle.Render("↑/↓ ou j/k: naviguer • enter: valider • q: quitter")
		b.WriteString(help + "\n")
	}
}

func (m QuizModel) renderFinished(b *strings.Builder) {
	// Title
	title := TitleStyle.Render("🎊 Quiz Terminé ! 🎊")
	b.WriteString(title + "\n\n")

	// Score
	percentage := float64(m.score) / float64(len(m.questions)) * 100
	scoreText := fmt.Sprintf("Score Final: %d/%d (%.1f%%)", m.score, len(m.questions), percentage)

	var scoreStyle lipgloss.Style
	if percentage >= 80 {
		scoreStyle = SuccessStyle
	} else if percentage >= 50 {
		scoreStyle = StatsStyle
	} else {
		scoreStyle = ErrorStyle
	}

	scoreBox := BoxStyle.Render(scoreStyle.Render(scoreText))
	b.WriteString(scoreBox + "\n\n")

	// Category
	catInfo := SubtitleStyle.Render(fmt.Sprintf("Catégorie: %s", m.category))
	b.WriteString(catInfo + "\n\n")

	// Encouragement
	var encouragement string
	if percentage == 100 {
		encouragement = "🏆 Parfait ! Tu es un(e) expert(e) !"
	} else if percentage >= 80 {
		encouragement = "🌟 Excellent travail !"
	} else if percentage >= 50 {
		encouragement = "👍 Pas mal, continue comme ça !"
	} else {
		encouragement = "💪 Continue à t'entraîner !"
	}

	b.WriteString(TitleStyle.Render(encouragement) + "\n\n")

	// Help
	help := HelpStyle.Render("enter ou q: retour au menu")
	b.WriteString(help + "\n")
}

func (m QuizModel) renderProgressBar() string {
	width := 60
	filled := int(float64(m.currentIndex) / float64(len(m.questions)) * float64(width))

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)

	style := lipgloss.NewStyle().Foreground(primaryColor)
	return style.Render(bar)
}

func (m QuizModel) GetScore() models.Score {
	return models.Score{
		Username: m.username,
		Category: m.category,
		Score:    m.score,
		Total:    len(m.questions),
	}
}
