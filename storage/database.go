package storage

import (
	"database/sql"
	"fmt"
	"quizz-ssh/models"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

// NewDatabase crée une nouvelle connexion à la base de données
func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Créer les tables si elles n'existent pas
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS scores (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			category TEXT NOT NULL,
			score INTEGER NOT NULL,
			total INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_category ON scores(category);
		CREATE INDEX IF NOT EXISTS idx_username ON scores(username);
	`)
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

// SaveScore enregistre un score
func (d *Database) SaveScore(score models.Score) error {
	_, err := d.db.Exec(
		"INSERT INTO scores (username, category, score, total, created_at) VALUES (?, ?, ?, ?, ?)",
		score.Username, score.Category, score.Score, score.Total, time.Now(),
	)
	return err
}

// GetLeaderboard récupère le top 10 pour une catégorie (ou global si category == "")
func (d *Database) GetLeaderboard(category string, limit int) ([]models.Score, error) {
	var rows *sql.Rows
	var err error

	if category == "" || category == "global" {
		// Leaderboard global: meilleur score par utilisateur toutes catégories confondues
		rows, err = d.db.Query(`
			SELECT username, 'global' as category, SUM(score) as total_score, SUM(total) as total_questions, MAX(created_at) as created_at
			FROM scores
			GROUP BY username
			ORDER BY total_score DESC
			LIMIT ?
		`, limit)
	} else {
		// Leaderboard par catégorie: meilleur score par utilisateur pour cette catégorie
		rows, err = d.db.Query(`
			SELECT username, category, MAX(score) as best_score, total, created_at
			FROM scores
			WHERE category = ?
			GROUP BY username
			ORDER BY best_score DESC
			LIMIT ?
		`, category, limit)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []models.Score
	for rows.Next() {
		var score models.Score
		err := rows.Scan(&score.Username, &score.Category, &score.Score, &score.Total, &score.CreatedAt)
		if err != nil {
			return nil, err
		}
		scores = append(scores, score)
	}

	return scores, nil
}

// GetUserBestScore récupère le meilleur score d'un utilisateur pour une catégorie
func (d *Database) GetUserBestScore(username, category string) (int, error) {
	var bestScore int
	err := d.db.QueryRow(
		"SELECT IFNULL(MAX(score), 0) FROM scores WHERE username = ? AND category = ?",
		username, category,
	).Scan(&bestScore)
	return bestScore, err
}

// Close ferme la connexion à la base de données
func (d *Database) Close() error {
	return d.db.Close()
}

// GetCategories récupère toutes les catégories distinctes
func (d *Database) GetCategories() ([]string, error) {
	rows, err := d.db.Query("SELECT DISTINCT category FROM scores ORDER BY category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var cat string
		if err := rows.Scan(&cat); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

func (d *Database) GetStats() (string, error) {
	var totalScores, uniqueUsers int
	err := d.db.QueryRow("SELECT COUNT(*), COUNT(DISTINCT username) FROM scores").Scan(&totalScores, &uniqueUsers)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Total parties: %d | Joueurs uniques: %d", totalScores, uniqueUsers), nil
}
