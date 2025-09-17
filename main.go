package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"goida/internal/app"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found")
	}

	application, err := app.New()
	if err != nil {
		log.Fatal("Failed to create application:", err)
	}
	defer application.Close()

	if err := application.Run(); err != nil {
		log.Fatal("Failed to run application:", err)
	}
}
