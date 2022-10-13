package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("DB_URL")))

	if err != nil {
		log.Panic("Could not connect to database", err)
	}

	return client
}

func GetBikeCollection(dbClient *mongo.Client) *mongo.Collection {
	return dbClient.Database("fifteen").Collection("bikes")
}
