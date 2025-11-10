package storage

import (
	"encoding/json"
	"os"
	"quizz-ssh/models"
)

// LoadQuestions charge les questions depuis un fichier JSON
func LoadQuestions(filepath string) ([]models.Question, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var quizData models.QuizData
	err = json.Unmarshal(data, &quizData)
	if err != nil {
		return nil, err
	}

	return quizData.Questions, nil
}

// GetQuestionsByCategory filtre les questions par catégorie
func GetQuestionsByCategory(questions []models.Question, category string) []models.Question {
	if category == "" || category == "all" {
		return questions
	}

	var filtered []models.Question
	for _, q := range questions {
		if q.Category == category {
			filtered = append(filtered, q)
		}
	}
	return filtered
}

// GetCategories récupère toutes les catégories uniques
func GetUniqueCategories(questions []models.Question) []string {
	categoryMap := make(map[string]bool)
	for _, q := range questions {
		categoryMap[q.Category] = true
	}

	var categories []string
	for cat := range categoryMap {
		categories = append(categories, cat)
	}
	return categories
}
