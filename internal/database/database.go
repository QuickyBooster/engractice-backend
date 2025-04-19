package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



type ServiceDb struct {
	Db *mongo.Client
}

var (
	host = os.Getenv("BLUEPRINT_DB_HOST")
	port = os.Getenv("BLUEPRINT_DB_PORT")
	//database = os.Getenv("BLUEPRINT_DB_DATABASE")
	clientInstance      *mongo.Client
	clientInstanceError error
	mongoOnce           sync.Once
)

func New() ServiceDb {
	client, err := GetMongoClient()
	if err != nil {
		log.Fatal(err)
	}
	return ServiceDb{
		Db: client,
	}
}

func (s *ServiceDb) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.Db.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("db down: %v", err)
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func GetMongoClient() (*mongo.Client, error) {
	mongoOnce.Do(func() {
		clientInstance, clientInstanceError = mongo.Connect(context.Background(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port)))
		if clientInstanceError != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", clientInstanceError)
		}
	})
	return clientInstance, clientInstanceError
}
