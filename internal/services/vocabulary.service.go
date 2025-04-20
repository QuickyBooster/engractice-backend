package services

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (v *VocabularyService) GetByID(id string) (*models.Vocabulary, error) {
	collection := v.collection.Db.Database("engractice").Collection("vocabulary")

	var vocabulary models.Vocabulary
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&vocabulary)
	if err != nil {
		log.Printf("Error fetching vocabulary by ID: %v", err)
		return nil, err
	}

	return &vocabulary, nil
}
func (v *VocabularyService) Create(vocabulary *models.Vocabulary) (*models.Vocabulary, error) {
	collection := v.collection.Db.Database("engractice").Collection("vocabulary")

	result, err := collection.InsertOne(context.Background(), vocabulary)
	if err != nil {
		log.Printf("Error creating vocabulary: %v", err)
		return nil, err
	}

	vocabulary.ID = result.InsertedID.(primitive.ObjectID)

	return vocabulary, nil
}
func (v *VocabularyService) Update(id *string, vocabulary *models.Vocabulary) (*models.Vocabulary, error) {
	collection := v.collection.Db.Database("engractice").Collection("vocabulary")

	update := bson.M{
		"$set": vocabulary,
	}

	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	if err != nil {
		log.Printf("Error updating vocabulary: %v", err)
		return nil, err
	}

	vocabulary.ID = *id.(primitive.ObjectID)

	return vocabulary, nil
}
func (v *VocabularyService) Delete(id *string) error {
	collection := v.collection.Db.Database("engractice").Collection("vocabulary")

	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		log.Printf("Error deleting vocabulary: %v", err)
		return err
	}

	return nil
}
func (v *VocabularyService) Search(query *string, page *int64) ([]models.Vocabulary, error) {
	collection := v.collection.Db.Database("engractice").Collection("vocabulary")

	findOptions := options.Find().
		SetSkip((*page - 1) * 100).
		SetLimit(100)

	cursor, err := collection.Find(context.Background(), bson.M{"$text": bson.M{"$search": query}}, findOptions)
	if err != nil {
		log.Printf("Error searching vocabulary: %v", err)
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
