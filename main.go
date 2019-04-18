package main

import (
	"log"
	"snift-backend/controllers"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	controllers.HandleRequests()
}
