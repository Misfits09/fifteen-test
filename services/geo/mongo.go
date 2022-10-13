package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getBikesLocationsCollection(dbClient *mongo.Client) *mongo.Collection {
	return dbClient.Database("fifteen").Collection("bikes_locations")
}

func setupIndexes(collection *mongo.Collection) {
	_, err := collection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{{
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

	if err != nil {
		log.Fatal("Could not create database indexes")
	}
}
