package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	h "fifteen/shared/helpers"
	"fifteen/shared/mongo"
	"fifteen/shared/rabbitmq"
	"fifteen/shared/structs"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	e := echo.New()
	dbClient := mongo.ConnectDB()
	channel, queue := rabbitmq.ConnectToRabbitMq()

	defer func() {
		if err := dbClient.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	bikeCollection := mongo.GetBikeCollection(dbClient)

	e.GET("/", func(c echo.Context) error {
		bikeList := make([]structs.Bike, 0)

		cursor, err := bikeCollection.Find(context.Background(), bson.D{{}})
		if err != nil {
			return h.SendErrorResponse(c, err, http.StatusInternalServerError)
		}

		err = cursor.All(context.Background(), &bikeList)
		if err != nil {
			return h.SendErrorResponse(c, err, http.StatusInternalServerError)
		}

		return c.JSON(http.StatusOK, bikeList)
	})

	e.GET("/:id", func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return h.SendErrorResponse(c, nil, http.StatusBadRequest)
		}

		result := bikeCollection.FindOne(context.Background(), bson.D{{Key: "_id", Value: id}})
		if result.Err() != nil {
			return h.SendErrorResponse(c, nil, http.StatusNotFound)
		}

		var bike structs.Bike

		err := result.Decode(&bike)
		if err != nil {
			return h.SendErrorResponse(c, err, http.StatusInternalServerError)
		}

		return c.JSON(http.StatusOK, bike)
	})

	e.POST("/", func(c echo.Context) error {
		bike := new(structs.Bike)

		err := c.Bind(bike)
		if err != nil || bike.LocationType != "Point" {
			return h.SendErrorResponse(c, nil, http.StatusBadRequest)
		}

		updateStatus, updateError := bikeCollection.UpdateOne(
			context.Background(),
			bson.M{"_id": bike.Id},
			bson.M{"$set": bike},
			options.Update().SetUpsert(true),
		)

		if updateError != nil {
			return h.SendErrorResponse(c, updateError, http.StatusInternalServerError)
		}

		if updateStatus.UpsertedCount > 0 || updateStatus.MatchedCount > 0 {
			json_locationEntry, _ := json.Marshal(structs.InternalLocationEntry{
				Id:   bike.Id,
				Time: time.Now().Unix(),
				Location: structs.Location{
					LocationType: bike.LocationType,
					LocationData: bike.LocationData,
				},
			})

			channel.PublishWithContext(context.Background(), "", queue.Name, false, false, amqp.Publishing{
				ContentType: "application/json",
				Body:        json_locationEntry,
			})

			return c.JSON(http.StatusOK, &structs.SuccessResponse{
				Success:     true,
				Description: strconv.FormatInt(updateStatus.MatchedCount, 10) + strconv.FormatInt(updateStatus.ModifiedCount, 10),
			})
		} else {
			return h.SendErrorResponse(c, nil, http.StatusInternalServerError)
		}
	})

	e.Logger.Fatal(e.Start(":8080"))
}
