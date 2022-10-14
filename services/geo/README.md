# Geo service

This service is responsible for recieving events about bikes locations. It stores all known locations values about a bike and can respond to request to know where was a bike at a given time.

## Usage

When running the docker exposes the 8080 port.

- `GET /` given a bike id and a time retrieves the last known location at that time. It should be called with two query params :

  - `id` a string being the id of the bike
  - `time` a timestamp following this format : `2006-01-02T15:04:05-07:00`

  The response is as such:

```json
{
  "id": "someBikeId",
  "location": {
    "type": "Point",
    "coordinates": [1.23, 4.56]
  },
  "time": "2021-10-13T13:02:21+02:00"
}
```

## Deployment

### Prefered

To make sure the service has all its depencies (including the mongoDB and the RabbitMQ queuing) it is prefered to use the provided `docker-compose.yml` provided at the root of the project. Refer to the [root README](../../README.md) for how to use it.

### Manual

If you want to deploy this service in a custom context you can use the provided Dockerfile directly. You should check that you have the following variables :

- DB_URL: The full URL to access the MongoDB (ex: `mongodb://root:example@mongo:27017`)
- RABBITMQ_URL: The full URL to access the RabbitMQ service (ex: `amqp://guest:guest@rabbitmq:5672/`)

And from the parent `services` directory run `docker build . -f geo/Dockerfile -e DB_URL=<Your URL> -e RABBITMQ_URL=<Your URL>`
