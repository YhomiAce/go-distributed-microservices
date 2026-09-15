package main

import (
	"context"
	"fmt"
	"log"
	"logger/data"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const webPort = "8082"

type Application struct {
	DB *mongo.Client
	Models data.Models
}

func main() {
	db := connectToDB()
	if db == nil {
		log.Fatal("Could not connect to the database")
	}
	log.Println("Connected to database")
	app := Application{
		DB: db,
		Models: data.New(db),
	}
	log.Printf("Starting logger-service on port: %s\n", webPort)
	server := &http.Server {
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

func connectToDB() *mongo.Client {
	uri := os.Getenv("MONGO_URI")

	for range 10 {
		ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
		client, err := mongo.Connect(options.Client().ApplyURI(uri))
		if err != nil {
			log.Println("Error opening database:", err)
			continue
		}
		if err := client.Ping(ctx, nil); err != nil {
			log.Println("Error pinging database:", err)
			cancel()
		}
		log.Println("Connected to Mongo")
		return client
	}
	log.Fatal("failed to connect to database after 10 attempts")
	return nil
}