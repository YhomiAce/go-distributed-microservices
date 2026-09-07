package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const webPort = "8081"

type Application struct {
	DB *sql.DB
	Models data.Models
}

func main() {
	db := connectToDB()
	if db == nil {
		log.Fatal("Could not connect to the database")
	}
	defer db.Close()
	log.Println("Connected to database")
	app := Application{
		DB: db,
		Models: data.New(db),
	}
	log.Printf("Starting auth-service on port: %s\n", webPort)
	server := &http.Server {
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for range 10 {
		db, err := sql.Open("pgx", dsn)
		if err != nil {
			log.Println("Error opening database:", err)
			continue
		}
		if err := db.Ping(); err != nil {
			log.Println("Error pinging database:", err)
			continue
		}
		return db
	}
	log.Fatal("failed to connect to database after 10 attempts")
	return nil
}