package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"quizz-ssh/models"
	"quizz-ssh/storage"
	"quizz-ssh/ui"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

const (
	host         = "0.0.0.0"
	port         = 2222
	dbPath       = "./data/quiz.db"
	questionsPath = "./questions.json"
)

var (
	db        *storage.Database
	questions []models.Question
)

func main() {
	// Créer le dossier data si nécessaire
	if err := os.MkdirAll("./data", 0755); err != nil {
		log.Fatalf("Erreur création dossier data: %v", err)
	}

	// Initialiser la base de données
	var err error
	db, err = storage.NewDatabase(dbPath)
	if err != nil {
		log.Fatalf("Erreur connexion DB: %v", err)
	}
	defer db.Close()

	// Charger les questions
	questions, err = storage.LoadQuestions(questionsPath)
	if err != nil {
		log.Printf("⚠️  Erreur chargement questions: %v", err)
		log.Printf("ℹ️  Utilisation de questions par défaut")
		questions = getDefaultQuestions()
	}

	log.Printf("✅ %d questions chargées", len(questions))

	// Configuration du serveur SSH
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Fatalf("Erreur création serveur: %v", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("🚀 Serveur SSH démarré sur %s:%d", host, port)
	log.Printf("📝 Connectez-vous avec: ssh -p %d %s", port, host)

	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Fatalf("Erreur serveur: %v", err)
		}
	}()

	<-done
	log.Println("🛑 Arrêt du serveur...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Fatalf("Erreur arrêt serveur: %v", err)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, active := s.Pty()
	if !active {
		wish.Fatalln(s, "no active terminal")
		return nil, nil
	}

	// Modèle initial: demander le pseudo
	m := &appModel{
		session: s,
		width:   pty.Window.Width,
		height:  pty.Window.Height,
		state:   stateUsername,
	}

	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

type appState int

const (
	stateUsername appState = iota
	stateMenu
	stateQuiz
	stateLeaderboard
)

type appModel struct {
	session  ssh.Session
	width    int
	height   int
	state    appState
	username string
	category string
	subModel tea.Model
}

func (m *appModel) Init() tea.Cmd {
	return nil
}

func (m *appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	// State machine
	switch m.state {
	case stateUsername:
		return m.updateUsername(msg)
	case stateMenu:
		return m.updateMenu(msg)
	case stateQuiz:
		return m.updateQuiz(msg)
	case stateLeaderboard:
		return m.updateLeaderboard(msg)
	}

	return m, nil
}

func (m *appModel) updateUsername(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.subModel == nil {
		m.subModel = ui.NewUsernameModel()
		return m, m.subModel.Init()
	}

	var cmd tea.Cmd
	m.subModel, cmd = m.subModel.Update(msg)

	// Vérifier si l'utilisateur a validé son pseudo
	if usernameModel, ok := m.subModel.(ui.UsernameModel); ok {
		if usernameModel.IsDone() {
			// Pseudo validé, passer au menu
			m.username = usernameModel.GetUsername()
			m.state = stateMenu
			m.subModel = nil
			return m, nil
		}
	}

	return m, cmd
}

func (m *appModel) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.subModel == nil {
		m.subModel = ui.NewMenuModel(m.username)
		return m.updateMenu(msg)
	}

	var cmd tea.Cmd
	m.subModel, cmd = m.subModel.Update(msg)

	// Vérifier si l'utilisateur a choisi
	if menuModel, ok := m.subModel.(ui.MenuModel); ok {
		if menuModel.IsDone() {
			choice := menuModel.GetChoice()
			m.subModel = nil

			log.Printf("DEBUG: Menu choice = %d", choice)

			switch choice {
			case ui.MenuQuiz:
				log.Printf("DEBUG: Switching to quiz state")
				m.category = "Cybersecurity"
				m.state = stateQuiz
			case ui.MenuLeaderboard:
				log.Printf("DEBUG: Switching to leaderboard state")
				m.category = "global"
				m.state = stateLeaderboard
			case ui.MenuQuit:
				log.Printf("DEBUG: Quitting")
				return m, tea.Quit
			}
		}
	}

	return m, cmd
}

func (m *appModel) updateQuiz(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.subModel == nil {
		log.Printf("DEBUG: Création quiz model avec %d questions pour %s", len(questions), m.username)
		// Toutes les questions (pas de filtre par catégorie)
		m.subModel = ui.NewQuizModel(m.username, questions, m.category)
		log.Printf("DEBUG: Quiz model créé, initialisation...")
		return m, m.subModel.Init()
	}

	var cmd tea.Cmd
	m.subModel, cmd = m.subModel.Update(msg)

	// Vérifier si le quiz est terminé
	if quizModel, ok := m.subModel.(ui.QuizModel); ok {
		if quizModel.IsDone() {
			// Sauvegarder le score
			score := quizModel.GetScore()
			if err := db.SaveScore(score); err != nil {
				log.Printf("Erreur sauvegarde score: %v", err)
			}
			log.Printf("DEBUG: Score sauvegardé, retour au menu")
			m.subModel = nil
			m.state = stateMenu
			return m, nil
		}
	}

	return m, cmd
}

func (m *appModel) updateLeaderboard(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.subModel == nil {
		scores, err := db.GetLeaderboard(m.category, 10)
		if err != nil {
			log.Printf("Erreur récupération leaderboard: %v", err)
			scores = []models.Score{}
		}

		stats, _ := db.GetStats()
		m.subModel = ui.NewLeaderboardModel(m.username, m.category, scores, stats)
		return m.updateLeaderboard(msg)
	}

	var cmd tea.Cmd
	m.subModel, cmd = m.subModel.Update(msg)

	// Vérifier si l'utilisateur veut retourner
	if lbModel, ok := m.subModel.(ui.LeaderboardModel); ok {
		if lbModel.IsDone() {
			m.subModel = nil
			m.state = stateMenu
			return m, nil
		}
	}

	return m, cmd
}

func (m *appModel) View() string {
	if m.subModel != nil {
		return m.subModel.View()
	}

	// Vue par défaut
	style := lipgloss.NewStyle().
		Padding(2).
		Foreground(lipgloss.Color("#00ff9f"))

	return style.Render("Chargement...")
}

// Questions par défaut pour démarrer
func getDefaultQuestions() []models.Question {
	return []models.Question{
		{
			ID:       1,
			Category: "Réseau",
			Text:     "Quel protocole est utilisé pour sécuriser HTTP ?",
			Options:  []string{"SSL/TLS", "FTP", "SMTP", "DNS"},
			Answer:   0,
		},
		{
			ID:       2,
			Category: "Cryptographie",
			Text:     "Qu'est-ce que AES ?",
			Options:  []string{"Un hash", "Un chiffrement symétrique", "Un chiffrement asymétrique", "Un protocole réseau"},
			Answer:   1,
		},
		{
			ID:       3,
			Category: "Réseau",
			Text:     "Quel port utilise SSH par défaut ?",
			Options:  []string{"21", "22", "23", "25"},
			Answer:   1,
		},
		{
			ID:       4,
			Category: "Web",
			Text:     "Qu'est-ce qu'une attaque XSS ?",
			Options:  []string{"Cross-Site Scripting", "Cross-Site Security", "eXtreme Site Security", "eXternal Script Source"},
			Answer:   0,
		},
		{
			ID:       5,
			Category: "Cryptographie",
			Text:     "Quelle est la taille d'un hash SHA-256 en bits ?",
			Options:  []string{"128", "192", "256", "512"},
			Answer:   2,
		},
	}
}
