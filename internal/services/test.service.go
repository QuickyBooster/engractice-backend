package services

import (
	"context"
	"engractice/internal/database"
	"engractice/internal/models"
	"log"
	"math/rand"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (t *TestService) FinishTest(test *models.Test) (string, error) {
	// find the test by id : test.ID then update the test.Correct and test.Wrong, keep the rest same as before
	collection := t.collection.Db.Collection(t.collectionTest)
	filter := bson.M{"_id": test.ID}
	update := bson.M{
		"$set": bson.M{
			"correct": test.Correct,
			"wrong":   test.Wrong,
		},
	}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Error updating test: %v", err)
		return "error", err
	}

	return "success", nil
}

func (t *TestService) GetAllTest(date string, tags string, nearestMode string, quantity string, page string) ([]models.Test, error) {
	collection := t.collection.Db.Collection(t.collectionTest)
	// filter & option
	filter := bson.M{}
	option := options.Find()
	if date != "" {
		date, err := time.Parse(time.RFC3339, date)
		if err != nil {
			log.Printf("Error parsing date: %v", err)
			return nil, err
		}
		filter["date"] = bson.M{"$gte": date}
	}
	if tags != "" {
		filter["tag"] = bson.M{"$regex": tags, "$options": "i"}
	}
	if nearestMode != "" {
		if neareastModeBool, err := strconv.ParseBool(nearestMode); err != nil {
			log.Printf("Error parsing nearestMode: %v", err)
			return nil, err
		} else {
			filter["nearestMode"] = neareastModeBool
		}
	}
	if quantity != "" {
		quantityInt, err := strconv.Atoi(quantity)
		if err != nil {
			log.Printf("Error parsing quantity: %v", err)
			return nil, err
		}
		filter["quantity"] = quantityInt
	}
	if page != "" {
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			log.Printf("Error parsing page: %v", err)
			return nil, err
		}
		option.SetSkip(int64(pageInt - 1)).SetLimit(10)
	}
	cursor, err := collection.Find(context.Background(), filter, option)
	if err != nil {
		log.Printf("Error searching vocabulary: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())
	var tests []models.Test
	if err := cursor.All(context.Background(), &tests); err != nil {
		log.Printf("Error decoding vocabulary: %v", err)
		return nil, err
	}

	return tests, nil
}

func (t *TestService) CreateTest(test *models.TestRequest) (models.Test, error) {

	// get the words from database
	collection := t.collection.Db.Collection(t.collectionTest)
	collectionVocab := t.collection.Db.Collection(t.collectionVocab)
	// filter & option
	filter := bson.M{}
	option := options.Find()
	if !test.NearestMode {
		filter = bson.M{
			"tag": bson.M{"$regex": test.Tags, "$options": "i"},
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

	var newTest models.Test
	newTest.Date = time.Now()
	if len(vocabularies) < test.Quantity {
		newTest.Words = vocabularies
	} else {

		rand.Shuffle(len(vocabularies), func(i, j int) {
			vocabularies[i], vocabularies[j] = vocabularies[j], vocabularies[i]
		})
		newTest.Words = vocabularies[:test.Quantity]
	}
	newTest.ID = primitive.NewObjectID()
	newTest.NearestMode = test.NearestMode

	if _, err := collection.InsertOne(context.Background(), newTest); err != nil {
		log.Printf("Error creating vocabulary: %v", err)
		return models.Test{}, err
	}

	return newTest, nil
}
