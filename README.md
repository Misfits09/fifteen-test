# Fifteen technical test project

## Usage

Refer to both `services/bike` and `services/geo` for technical usage directions.

A postman collection example is available at `data/postman_collection.json` for ease of use of both APIs.

## Deployment

### Prefered

To make sure the service all services have their depencies (including the mongoDB and the RabbitMQ queuing) it is prefered to use the provided `docker-compose.yml` provided at the root of this project.

_NB: Do not hesitate to update the `docker-compose.yml` file to use actual secrets on both mongodb and rabbitmq services._

To deploy simply make sure Docker is installed on your machine and run `make deploy` or alternativelly `docker-compose up -d`

### Manual

Have a MongoDB and RabbitMQ instances setup. Refer to `services/bike` and `services/geo` READMEs to run both services.

## Design Choices

**Coming Soon**
