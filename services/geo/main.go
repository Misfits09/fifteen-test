package main

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	mongodb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"fifteen/shared/helpers"
	"fifteen/shared/mongo"
	"fifteen/shared/rabbitmq"
	"fifteen/shared/structs"
)

const TIME_FORMAT = "2006-01-02T15:04:05-07:00"

func main() {
	channel, queue := rabbitmq.ConnectToRabbitMq()
	dbClient := mongo.ConnectDB()
	collection := getBikesLocationsCollection(dbClient)
	setupIndexes(collection)

	var wg sync.WaitGroup
	wg.Add(2)

	go startRabbitListener(channel, queue, collection, &wg)
	go startServer(collection, &wg)

	wg.Wait()
}

func startServer(collection *mongodb.Collection, wg *sync.WaitGroup) {
	defer wg.Done()

	e := echo.New()

	e.GET("/:id", func(c echo.Context) error {
		id, timestamp := c.Param("id"), c.QueryParams().Get("time")
		parsedTime, err := time.Parse(TIME_FORMAT, timestamp)

		if timestamp == "" || id == "" || err != nil {
			return helpers.SendErrorResponse(c, err, http.StatusBadRequest)
		}

		entry := new(structs.InternalLocationEntry)
		err = collection.FindOne(context.Background(), bson.M{"bikeId": id, "time": bson.M{"$lte": parsedTime.Unix()}}, options.FindOne().SetSort(bson.M{"time": -1})).Decode(&entry)

		if err != nil {
			return helpers.SendErrorResponse(c, nil, http.StatusNotFound)
		}

		output := structs.APILocationEntry{
			Id:        entry.Id,
			Location:  entry.Location,
			TimeStamp: time.Unix(entry.Time, 0).Format(TIME_FORMAT),
		}

		return c.JSON(http.StatusOK, output)
	})

	e.Logger.Fatal(e.Start(":8080"))
}

func startRabbitListener(channel *amqp.Channel, queue amqp.Queue, collection *mongodb.Collection, wg *sync.WaitGroup) {
	deliveries, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Panic("Could not create Rabbit MQ Consumer")
		panic(1)
	}

	for d := range deliveries {
		locationEntry := new(structs.InternalLocationEntry)
		json.Unmarshal(d.Body, &locationEntry)
		_, err := collection.InsertOne(context.Background(), locationEntry)

		if err != nil {
			log.Error(err)
			continue
		}

		d.Ack(false)
	}
}
