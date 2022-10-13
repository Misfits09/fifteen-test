### Builder image
FROM golang:alpine AS builder

WORKDIR /app/build

# Setup shared module
COPY ./shared ./shared
WORKDIR /app/build/shared
RUN go mod download

# Setup Bike module
COPY ./bike /app/build/bike
WORKDIR /app/build/bike
RUN go mod download

RUN go build -o ../output



### Runtime image
FROM golang:alpine

WORKDIR /app
COPY --from=builder /app/build/output .

CMD [ "./output"]
EXPOSE 8080
