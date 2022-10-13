package main

import (
	"context"
	"errors"
	"fifteen/shared/structs"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func listBikes(c echo.Context, bikeCollection mongo.Collection) ([]structs.Bike, error) {
	bikeList := make([]structs.Bike, 0)

	cursor, err := bikeCollection.Find(context.Background(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &bikeList)
	if err != nil {
		return nil, err
	}

	return bikeList, nil
}

func getBike(c echo.Context, bikeCollection mongo.Collection, id string) (*structs.Bike, error) {
	result := bikeCollection.FindOne(context.Background(), bson.D{{Key: "_id", Value: id}})
	if result.Err() != nil {
		return nil, result.Err()
	}

	var bike structs.Bike

	err := result.Decode(&bike)
	if err != nil {
		return nil, err
	}

	return &bike, nil
}

func insertOrUpdateBike(c echo.Context, bikeCollection mongo.Collection, bike *structs.Bike) error {
	updateStatus, updateError := bikeCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": bike.ID},
		bson.M{"$set": bike},
		options.Update().SetUpsert(true),
	)

	if updateError != nil {
		return updateError
	}

	if updateStatus.UpsertedCount > 0 || updateStatus.MatchedCount > 0 {
		return nil
	}
	return errors.New("not inserted")
}
