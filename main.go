package main

import (
	"log"
	"net/http"
	"os"
	"snift-backend/controllers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func handleRequests() {
	port := os.Getenv("PORT")
	log.Print("Server starting at PORT ", port)
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", controllers.HomePage).Methods("GET")
	myRouter.HandleFunc("/scores", controllers.GetScore).Methods("GET")
	log.Fatal(http.ListenAndServe(port, myRouter))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	handleRequests()
}
