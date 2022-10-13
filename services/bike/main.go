package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"fifteen/shared/db"
	h "fifteen/shared/helpers"
	"fifteen/shared/rabbitmq"
	"fifteen/shared/structs"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	e := echo.New()
	dbClient := db.ConnectDB()
	channel, queue := rabbitmq.ConnectToRabbitMq()

	defer func() {
		if err := dbClient.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	bikeCollection := getBikeCollection(dbClient)

	e.GET("/", func(c echo.Context) error {
		bikeList, err := listBikes(c, *bikeCollection)
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

		bike, err := getBike(c, *bikeCollection, id)
		if err != nil {
			return h.SendErrorResponse(c, err, http.StatusNotFound)
		}

		return c.JSON(http.StatusOK, *bike)
	})

	e.POST("/", func(c echo.Context) error {
		bike := new(structs.Bike)

		err := c.Bind(bike)
		if err != nil || bike.LocationType != "Point" {
			return h.SendErrorResponse(c, nil, http.StatusBadRequest)
		}

		err = insertOrUpdateBike(c, *bikeCollection, bike)
		if err != nil {
			return h.SendErrorResponse(c, err, http.StatusInternalServerError)
		}

		err = fireBikeLocationEvent(bike, channel, queue.Name)
		if err != nil {
			return h.SendErrorResponse(c, err, http.StatusInternalServerError)
		}

		return c.JSON(http.StatusOK, &structs.SuccessResponse{
			Success: true,
		})
	})

	e.Logger.Fatal(e.Start(":8080"))
}

func fireBikeLocationEvent(bike *structs.Bike, channel *amqp.Channel, queueName string) error {
	jsonLocationEntry, _ := json.Marshal(structs.InternalLocationEntry{
		ID:   bike.ID,
		Time: time.Now().Unix(),
		Location: structs.Location{
			LocationType: bike.LocationType,
			LocationData: bike.LocationData,
		},
	})

	err := channel.PublishWithContext(context.Background(), "", queueName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonLocationEntry,
	})

	return err
}
