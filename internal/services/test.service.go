package services

import (
	"engractice/internal/database"
	"engractice/internal/models"
	"log"
	"math/rand"
	"sort"
	"time"
)

type TestService struct {
	db *database.Database
}

// CreateTest
// FinishTest (upload the result)

func NewTestService(db *database.Database) *TestService {
	return &TestService{
		db: db,
	}
}

func (t *TestService) CreateTest(numberOfWord int, tag string) (models.Test, error) {
	words, err := t.db.GetSpreadsheetData()
	if err != nil {
		log.Printf("ERROR %v", err)
		return models.Test{}, err
	}
	test := models.Test{}
	if tag != "" {
		rand.Shuffle(len(words), func(i, j int) {
			words[i], words[j] = words[j], words[i]
		})
		sort.Slice(words, func(i, j int) bool {
			return words[i].Point < words[j].Point
		})

		if numberOfWord > len(words) {
			numberOfWord = len(words)
		}
		test = models.Test{
			Words:       words[:numberOfWord],
			Date:        time.Now(),
			Correct:     make([]models.Vocabulary, 0),
			Wrong:       make([]models.Vocabulary, 0),
			NearestMode: true,
		}
		log.Printf("INFO Created test (no tags) with %d words", len(test.Words))
	} else {
		// if tag is provided, filter the words by tag
		filteredWords := []models.Vocabulary{}
		for _, word := range words {
			if word.Tag == tag {
				filteredWords = append(filteredWords, word)
			}
		}
		if len(filteredWords) == 0 {
			return test, nil // No words found with the given tag
		}
		rand.Shuffle(len(filteredWords), func(i, j int) {
			filteredWords[i], filteredWords[j] = filteredWords[j], filteredWords[i]
		})
		sort.Slice(filteredWords, func(i, j int) bool {
			return filteredWords[i].Point < filteredWords[j].Point
		})
		if numberOfWord > len(filteredWords) {
			numberOfWord = len(filteredWords)
		}
		test = models.Test{
			Words:       filteredWords[:numberOfWord],
			Date:        time.Now(),
			Correct:     make([]models.Vocabulary, 0),
			Wrong:       make([]models.Vocabulary, 0),
			NearestMode: false,
		}
		log.Printf("INFO Created test (yes tags) with %d words", len(test.Words))
	}
	return test, err
}
