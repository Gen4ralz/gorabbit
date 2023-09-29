PRODUCER_BINARY=producerApp
CONSUMER_BINARY=consumerApp

up_build: build_producer build_consumer
		docker-compose down
		docker-compose up --build -d

build_producer:
		cd producer-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${PRODUCER_BINARY} .

build_consumer:
		cd consumer-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${CONSUMER_BINARY} .