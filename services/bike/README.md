# Bike service

This service is responsible for recieving updated about bike, it stores the current value of a bike and fires an event to a queuing service

## Usage

When running the docker exposes the 8080 port.

- `GET /` will yield the list of all known bikes with a JSON Array response of the following objects :

```json
{
  "id": "someBikeId",
  "location": {
    "type": "Point",
    "coordinates": [1.23, 4.56]
  }
}
```

- `GET /<someBikeId>` will get the location of bike given its id with the same format as above

- `POST /` allows to update a bike location. If the bike did not exist in database it will be created. The Payload should be formatted accordingly to the `Content-Type` header and have tha same keys as the payload above.

## Deployment

### Prefered

To make sure the service has all its depencies (including the mongoDB and the RabbitMQ queuing) it is prefered to use the provided `docker-compose.yml` provided at the root of the project. Refer to the [root README](../../README.md) for how to use it.

### Manual

If you want to deploy this service in a custom context you can use the provided Dockerfile directly. You should check that you have the following variables :

- DB_URL: The full URL to access the MongoDB (ex: `mongodb://root:example@mongo:27017`)
- RABBITMQ_URL: The full URL to access the RabbitMQ service (ex: `amqp://guest:guest@rabbitmq:5672/`)

And from the parent `services` directory run `docker build . -f bike/Dockerfile -e DB_URL=<Your URL> -e RABBITMQ_URL=<Your URL>`
