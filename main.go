package main

import (
	"log"
	"net/http"
	"snift-backend/controllers"

	"github.com/gorilla/mux"
)

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", controllers.HomePage).Methods("GET")
	myRouter.HandleFunc("/scores", controllers.GetScore).Methods("GET")
	log.Fatal(http.ListenAndServe(":9700", myRouter))
}

func main() {
	handleRequests()
}
