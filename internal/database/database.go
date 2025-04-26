package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ServiceDb struct {
	Db     *mongo.Database
	DbName string
}

var (
	host     = os.Getenv("BLUEPRINT_DB_HOST")
	port     = os.Getenv("BLUEPRINT_DB_PORT")
	user     = os.Getenv("BLUEPRINT_DB_USERNAME")
	password = os.Getenv("BLUEPRINT_DB_ROOT_PASSWORD")
	//database = os.Getenv("BLUEPRINT_DB_DATABASE")
	clientInstance      *mongo.Database
	clientInstanceError error
	mongoOnce           sync.Once
)

func New(dbName string) ServiceDb {
	client, err := getMongoClient(dbName)
	if err != nil {
		log.Fatal(err)
	}
	return ServiceDb{
		Db:     client,
		DbName: dbName,
	}
}

func (s *ServiceDb) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.Db.Client().Ping(ctx, nil)
	if err != nil {
		log.Fatalf("db down: %v", err)
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func getMongoClient(dbName string) (*mongo.Database, error) {
	var instance *mongo.Client
	mongoOnce.Do(func() {
		connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%s", user, password, host, port)
		instance, clientInstanceError = mongo.Connect(context.Background(), options.Client().ApplyURI(connectionString))
		if clientInstanceError != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", clientInstanceError)
		}
		ensureDatabaseAndCollection(instance, dbName)
	})
	return instance.Database(dbName), clientInstanceError
}

func ensureDatabaseAndCollection(client *mongo.Client, dbName string) {

	if dbName == "" {
		log.Fatal("Database name is not set in the environment variable BLUEPRINT_DB_DATABASE")
	}

	ctx := context.Background()
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Failed to list databases: %v", err)
	}

	dbExists := false
	for _, db := range databases {
		if db == dbName {
			dbExists = true
			break
		}
	}

	if !dbExists {
		log.Printf("Database %s does not exist. Creating it...", dbName)
		client.Database(dbName).CreateCollection(ctx, "vocabulary")
	}

	collections, err := client.Database(dbName).ListCollectionNames(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Failed to list collections: %v", err)
	}

	vocabCollection := false
	testCollection := false
	for _, collection := range collections {
		if collection == "vocabulary" {
			vocabCollection = true
		}
		if collection == "test" {
			testCollection = true
		}
	}

	if !vocabCollection {
		log.Printf("Collection 'vocabulary' does not exist. Creating it...")
		err = client.Database(dbName).CreateCollection(ctx, "vocabulary")
		if err != nil {
			log.Fatalf("Failed to create collection 'vocabulary': %v", err)
		}
	}
	if !testCollection {
		log.Printf("Collection 'test' does not exist. Creating it...")
		err = client.Database(dbName).CreateCollection(ctx, "test")
		if err != nil {
			log.Fatalf("Failed to create collection 'test': %v", err)
		}
	}
}
