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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"fifteen/shared/db"
	h "fifteen/shared/helpers"
	"fifteen/shared/rabbitmq"
	"fifteen/shared/structs"
)

const TimeFormat = "2006-01-02T15:04:05-07:00"
const ParalleledProcesses = 2

func main() {
	channel, queue := rabbitmq.ConnectToRabbitMq()
	dbClient := db.ConnectDB()
	collection := getBikesLocationsCollection(dbClient)
	setupCollectionIndexes(collection)

	var wg sync.WaitGroup
	wg.Add(ParalleledProcesses)

	go startRabbitListener(channel, queue, collection, &wg)
	go startServer(collection, &wg)

	wg.Wait()
}

func startServer(collection *mongo.Collection, wg *sync.WaitGroup) {
	defer wg.Done()

	e := echo.New()

	e.GET("/:id", func(c echo.Context) error {
		id, timestamp := c.Param("id"), c.QueryParams().Get("time")
		parsedTime, err := time.Parse(TimeFormat, timestamp)

		if timestamp == "" || id == "" || err != nil {
			return h.SendErrorResponse(c, err, http.StatusBadRequest)
		}

		entry := new(structs.InternalLocationEntry)
		err = collection.FindOne(context.Background(), bson.M{"bikeId": id, "time": bson.M{"$lte": parsedTime.Unix()}}, options.FindOne().SetSort(bson.M{"time": -1})).Decode(&entry)
		if err != nil {
			return h.SendErrorResponse(c, nil, http.StatusNotFound)
		}

		output := structs.APILocationEntry{
			ID:        entry.ID,
			Location:  entry.Location,
			TimeStamp: time.Unix(entry.Time, 0).Format(TimeFormat),
		}

		return c.JSON(http.StatusOK, output)
	})

	e.Logger.Fatal(e.Start(":8080"))
}

func startRabbitListener(channel *amqp.Channel, queue amqp.Queue, collection *mongo.Collection, wg *sync.WaitGroup) {
	deliveries, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Panic("Could not create Rabbit MQ Consumer")
	}

	for d := range deliveries {
		locationEntry := new(structs.InternalLocationEntry)
		err := json.Unmarshal(d.Body, &locationEntry)

		if h.LogIfIsError(err) {
			continue
		}
		_, err = collection.InsertOne(context.Background(), locationEntry)

		if h.LogIfIsError(err) {
			continue
		}

		h.LogIfIsError(d.Ack(false))
	}
}
