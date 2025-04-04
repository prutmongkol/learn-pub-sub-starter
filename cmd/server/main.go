package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Connected to RabbitMQ")

	messageChannel, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}
	// defer messageChannel.Close()
	
	exchange := routing.ExchangePerilDirect
	key := routing.PauseKey
	msg := routing.PlayingState{IsPaused: true}

	err = pubsub.PublishJSON(messageChannel, exchange, key, msg)
	if err != nil {
		log.Fatalf("failed to publish message: %v", err)
	}

	// wait for ctrl+c
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	<-signalChannel
	fmt.Println("Shutting down...")
}
