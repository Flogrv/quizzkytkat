package models

import "time"

// Question représente une question du quiz
type Question struct {
	ID              int      `json:"id"`
	Category        string   `json:"category"`
	Text            string   `json:"text"`
	Options         []string `json:"options"`
	Answer          int      `json:"answer"`          // Index de la bonne réponse (original)
	ShuffledOptions []string `json:"-"`               // Options mélangées (pas sauvegardé en JSON)
	ShuffledAnswer  int      `json:"-"`               // Index de la bonne réponse après shuffle
}

// Score représente le score d'un utilisateur
type Score struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Category  string    `json:"category"`
	Score     int       `json:"score"`
	Total     int       `json:"total"`
	CreatedAt time.Time `json:"created_at"`
}

// QuizData contient toutes les questions
type QuizData struct {
	Questions []Question `json:"questions"`
}
