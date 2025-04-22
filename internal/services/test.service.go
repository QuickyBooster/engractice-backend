package services

import (
	"context"
	"engractice/internal/database"
	"engractice/internal/models"
	"log"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TestService struct {
	collection      *database.ServiceDb
	dbName          string
	collectionTest  string
	collectionVocab string
}

func NewTestService(db *database.ServiceDb, database string) *TestService {
	return &TestService{
		collection:      db,
		dbName:          database,
		collectionTest:  "test",
		collectionVocab: "vocabulary",
	}
}

func (t *TestService) FinishTest(test *models.Test) (models.Test, error) {
	panic("unimplemented")
}

func (t *TestService) GetAllTest(date string, tags string, nearestMode string, quantity string, page string) ([]models.Test, error) {
	panic("unimplemented")
}

func (t *TestService) CreateTest(test *models.TestRequest) (models.Test, error) {
	var newTest models.Test

	// get the words from database
	collection := t.collection.Db.Collection(t.collectionTest)
	collectionVocab := t.collection.Db.Collection(t.collectionVocab)
	// filter & option
	filter := bson.M{}
	option := options.Find()
	if !test.NearestMode {
		filter = bson.M{
			"tags": bson.M{"$regex": test.Tags, "$options": "i"},
		}
	} else {
		option.SetLimit(int64(test.Quantity)).SetSort(bson.M{"created_at": 1})
	}

	cursor, err := collectionVocab.Find(context.Background(), filter, option)
	if err != nil {
		log.Printf("Error searching vocabulary and creating test: %v", err)
		return models.Test{}, err
	}
	defer cursor.Close(context.Background())
	var vocabularies []models.Vocabulary
	if err := cursor.All(context.Background(), &vocabularies); err != nil {
		log.Printf("Error decoding vocabulary: %v", err)
		return models.Test{}, err
	}

	newTest.Date = time.Now()
	if len(vocabularies) < test.Quantity {
		newTest.Words = vocabularies
	} else {

		rand.Shuffle(len(vocabularies), func(i, j int) {
			vocabularies[i], vocabularies[j] = vocabularies[j], vocabularies[i]
		})
		newTest.Words = vocabularies[:test.Quantity]
	}

	if _, err := collection.InsertOne(context.Background(), newTest); err != nil {
		log.Printf("Error creating vocabulary: %v", err)
		return models.Test{}, err
	}

	return newTest, nil
}
