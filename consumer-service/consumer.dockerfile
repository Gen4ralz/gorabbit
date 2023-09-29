FROM golang:1.20-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o consumerApp .

RUN chmod +x /app/consumerApp

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/consumerApp /app

CMD [ "/app/consumerApp" ]