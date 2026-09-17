package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "8080"

type Application struct {
	RabbitConn *amqp.Connection
}

func main() {
	conn, err := connectToRabbitMQ()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	app := Application{
		RabbitConn: conn,
	}
	log.Printf("Starting broker-service on port: %s\n", webPort)
	server := &http.Server {
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
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

		log.Printf("Backing off for %v...\n", backoff)
		time.Sleep(backoff)
	}
	return  nil, err
}