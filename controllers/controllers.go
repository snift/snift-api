package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	services "snift-backend/services"
)

// HomePage - the default root endpoint of Snift Backend
func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /")
	fmt.Fprintf(w, "Welcome to the Snift!")
}

// GetScore - GET /scores handler
func GetScore(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	log.Print("GET /scores")
	response := services.CalculateOverallScore(url)
	json.NewEncoder(w).Encode(response)
}
