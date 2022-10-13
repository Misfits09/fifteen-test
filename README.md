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

Services are all defined in `services` folder and share as much as possible of their code through the `shared` module.

A lightweight http server has been use called `Echo`. This one has been chosen as there is not big complex logic to be done in either services and provided all necessary features while being very popular/maintained.

From the start both services have been run with `docker-compose` to allow better DX and overall quicker development. Some commands have been also added in the root Makefile for simplicity.
`RabbitMQ` has been chosen as the queuing service to allow it to be easily deployed locally along other services in a dockerised environment.

Both services are developped around the `shared` module which mainly allows to abstract the logic to connect to `RabbitMQ` or `MongoDB` and have common structure types which is very important to make sure communications run smoothly.
For example using shared types allow to have safe JSON (De)Marshalling through the queuing system.
Finally some helpers have been defined there for other common logics (like in our case generating an error responses).

As we only have few routes per service, all the routing and request handling has been centralized in both `main.go` despite having some of the logic splitted in other files for visibility.

Finally some simple linting has been added to the gitflow in the Github Actions CI and triggered on Pull Requests.
