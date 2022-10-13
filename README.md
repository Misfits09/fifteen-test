<img src="https://storage.googleapis.com/zoov-assets/zoov-logo.png" alt="zoov log" width="200px">

# Backend technical test

In this test, you are given to develop two services working together to achieve something similar to what we are developing here at Zoov.

We use mostly `Go` and `Mongo` but feel free to use both the language(s) and the tool(s) (databases, queues, etc.) you think that match the best the needs of this test.

> :information_source: Most of the cloud providers have free tiers, do not hesitate to use them!

You should have one folder per service, each coming with a `README.md` file explaining how to build and run the service. Finally, you must provide a main `README.md` file containing an explanation of your design choices. Do not hesitate to create another folder for common dependencies if revelant.

Once you are done, store the code and documentation in a private `Github`, `Bitbucket` or `GitLab` repository and share it with us by providing read access to our accounts `MathieuZoov`, `schmurfy` and `nhoway`.

If you have any question about the technical choices, the architecture or the subject of the test itself, feel free to ask!

### Bike

This service should listen on port `8081` and expose three endpoints to:

- Return the list of all the bikes from database.
- Given a bike id, return the corresponding bike.
- Given a bike id and a location (RFC 7946 GeoJSON), update the bike location and send an event to the message broker of your choice (Google Pub/Sub, Amazon SQS, NSQ, ...).

Here is a JSON representation of a bike:

    {
      "id": "bb2cdchl52n4orsopmtg",
      "location": {
        "type": "Point",
        "coordinates": [2.2861460, 48.8268020],
      }
    }

You will find a `bikes.json` file in this repository containing a list of bikes (located in Paris) you can use to populate your database.

### Geo

This service will listen the events sent from the bike service on location updates and will store them. It should also listen on port `8082` and expose one endpoint to:

- Given a bike id and a timestamp, retrieve the bike location at that time

Here is a JSON representation of a record:

    {
      "id": "bb2cdb1l52n4oiuufrig",
      "location": {
        "type": "Point",
        "coordinates": [2.2861460, 48.8268020],
      }
      "time": "2018-04-04T14:40:05+02:00",
    }

### Bonus

- Auth
- Docker
- Logging
- Monitoring
- Unit tests
- CI/CD
- UI
- ... Anything else you think would be relevant !

> :warning: Please note that we will mostly judge you on the quality of the mandatory part of this test!
