package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"snift-backend/controllers"
)

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", controllers.HomePage).Methods("GET")
	myRouter.HandleFunc("/get_score", controllers.GetScore).Methods("GET")
	log.Fatal(http.ListenAndServe(":9700", myRouter))
}

func main() {
	handleRequests()
}