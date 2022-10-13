package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mongodb "go.mongodb.org/mongo-driver/mongo"
)

func getBikesLocationsCollection(dbClient *mongo.Client) *mongo.Collection {
	return dbClient.Database("fifteen").Collection("bikes_locations")
}

func setupIndexes(collection *mongo.Collection) {
	collection.Indexes().CreateMany(context.Background(), []mongodb.IndexModel{{
		Keys: bson.M{
			"bikeId": 1,
		},
		Options: nil,
	}, {
		Keys: bson.M{
			"time": -1,
		},
		Options: nil,
	}})
}
