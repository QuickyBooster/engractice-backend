package services

import (
	"engractice/internal/database"
	"engractice/internal/models"
)

type VocabularyService struct {
	db *database.Database
}

func NewVocabularyService(db *database.Database) *VocabularyService {
	return &VocabularyService{
		db: db,
	}
}
func (v *VocabularyService) GetAllWords() ([]models.Vocabulary, error) {
	words, err := v.db.GetSpreadsheetData()
	if err != nil {
		return nil, err
	}
	return words, nil
}
func (v *VocabularyService) UpdateWords(words []models.Vocabulary) error {
	err := v.db.UpdateSpreadsheetData(words)
	if err != nil {
		return err
	}
	return nil
}
