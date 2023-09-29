FROM golang:1.20-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o producerApp .

RUN chmod +x /app/producerApp

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/producerApp /app

CMD [ "/app/producerApp" ]