package handlers

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Vocabulary struct {
	ID         string    `json:"id" bson:"_id,omitempty"`
	English    string    `json:"english" bson:"english"`
	Vietnamese string    `json:"vietnamese" bson:"vietnamese"`
	AudioLink  string    `json:"audioLink" bson:"audioLink"`
	Tags       []string  `json:"tags" bson:"tags"`
	CreatedAt  time.Time `json:"createdAt" bson:"createdAt"`
}

var collection *mongo.Collection

func InitHandlers(db *mongo.Database) {
	collection = db.Collection("vocabularies")
}

func AddVocabulary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var vocab Vocabulary
	if err := json.NewDecoder(r.Body).Decode(&vocab); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	vocab.CreatedAt = time.Now()
	_, err := collection.InsertOne(context.Background(), vocab)
	if err != nil {
		http.Error(w, "Failed to add vocabulary", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(vocab)
}

func GetVocabularies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch vocabularies", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var vocabularies []Vocabulary
	if err := cursor.All(context.Background(), &vocabularies); err != nil {
		http.Error(w, "Failed to parse vocabularies", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(vocabularies)
}

func PracticeVocabularies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	numWords := 10 // Default number of words to practice
	if n := r.URL.Query().Get("num"); n != "" {
		fmt.Sscanf(n, "%d", &numWords)
	}

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch vocabularies", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var vocabularies []Vocabulary
	if err := cursor.All(context.Background(), &vocabularies); err != nil {
		http.Error(w, "Failed to parse vocabularies", http.StatusInternalServerError)
		return
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(vocabularies), func(i, j int) { vocabularies[i], vocabularies[j] = vocabularies[j], vocabularies[i] })

	if numWords > len(vocabularies) {
		numWords = len(vocabularies)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(vocabularies[:numWords])
}