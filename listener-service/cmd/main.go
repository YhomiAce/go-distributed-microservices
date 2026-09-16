package main

import (
	"listener/lib/events"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	log.Println("Starting listener-service")
	// try to connect to rabbitmq
	conn, err := connectToRabbitMQ()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	// start listening to messages
	log.Println("Listening and consuming for RabbitMQ messages")
	consumer, err := events.NewConsumer(conn)
	if err != nil {
		panic(err)
	}


	// create consumer

	// watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
		return
	}
}


func connectToRabbitMQV2() (*amqp.Connection, error) {
	var count int
	var backoff  = 2 * time.Second

	var connection *amqp.Connection
	var url = os.Getenv("RABBITMQ_URL")
	for {
		conn, err := amqp.Dial(url)
		if err != nil {
			log.Println("Failed to connect to Rabbitmq, RabbitMq not ready yet...")
			count++
		}else {
			log.Println("Connected to RabbitMQ")
			connection = conn
			break
		}

		if count > 10 {
			log.Println("Could not connect to Rabbitmq, exiting...")
			return nil, err
		}

		backoff = time.Duration(math.Pow(float64(count), 2)) * time.Second

		log.Println("Backing off...")
		time.Sleep(backoff)
		continue
	}
	return  connection, nil
}


func connectToRabbitMQ() (*amqp.Connection, error) {
	var err error
	var conn *amqp.Connection
	var url = os.Getenv("RABBITMQ_URL")
	maxTries := 10
	for i:= 1; i < maxTries; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			log.Println("Connected to RabbitMQ")
			return conn, nil
		}

		log.Printf("RabbitMQ not ready yet (attempt %d/%d): %v", i, maxTries, err)

		backoff := time.Duration(i*i) * time.Second

		log.Println("Backing off for %v...", backoff)
		time.Sleep(backoff)
	}
	return  nil, err
}