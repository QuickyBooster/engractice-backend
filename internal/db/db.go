package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectDB(uri string) *mongo.Database {
	var err error
	Client, err = mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := Client.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")
	return Client.Database("vocabularyDB")
}

func DisconnectDB() {
	if err := Client.Disconnect(context.Background()); err != nil {
		log.Fatal(err)
	}
	log.Println("Disconnected from MongoDB!")
}