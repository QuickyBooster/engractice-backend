package services

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"engractice/internal/database"
	"engractice/internal/models"
)

type VocabularyService struct {
	collection *database.ServiceDb
}

func NewVocabularyService(db *database.ServiceDb) *VocabularyService {
	return &VocabularyService{
		collection: db,
	}
}

func (v *VocabularyService) GetAll(page int64) ([]models.Vocabulary, error) {
	const pageSize int64 = 100
	skip := (page - 1) * pageSize

	collection := v.collection.Db.Database("engractice").Collection("vocabulary")

	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(pageSize)

	cursor, err := collection.Find(context.Background(), bson.M{}, findOptions)
	if err != nil {
		log.Printf("Error fetching vocabulary: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var vocabularies []models.Vocabulary
	if err := cursor.All(context.Background(), &vocabularies); err != nil {
		log.Printf("Error decoding vocabulary: %v", err)
		return nil, err
	}

	return vocabularies, nil
}
