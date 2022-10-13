package main

import "go.mongodb.org/mongo-driver/mongo"

func getBikeCollection(dbClient *mongo.Client) *mongo.Collection {
	return dbClient.Database("fifteen").Collection("bikes")
}
