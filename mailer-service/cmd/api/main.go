package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const webPort = "8083"

type Application struct {
	Mailer Mail
}

func main() {
	app := Application{
		Mailer: createMail(),
	}
	log.Printf("Starting mailer-service on port: %s\n", webPort)
	server := &http.Server {
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

func createMail() Mail  {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	m := Mail{
		Domain: os.Getenv("MAIL_DOMAIN"),
		Host: os.Getenv("MAIL_HOST"),
		Username: os.Getenv("MAIL_USERNAME"),
		Password: os.Getenv("MAIL_PASSWORD"),
		Encryption: os.Getenv("MAIL_ENCRYPTION"),
		FromName: os.Getenv("MAIL_FROM_NAME"),
		FromAddress: os.Getenv("MAIL_FROM_ADDRESS"),
		Port: port,
	}
	return m
}