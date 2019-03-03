package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	services "snift-backend/services"
	utils "snift-backend/utils"
)

// HomePage - the default root endpoint of Snift Backend
func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /")
	fmt.Fprintf(w, "Welcome to Snift!")
}

// GetScore - GET /scores handler
func GetScore(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	log.Print("GET /scores")

	err := utils.IsValidURL(url)
	if err != nil {
		utils.BadRequest(w, true, "Invalid URL")
		return
	}
	response, _ := services.CalculateOverallScore(url)
	body, jsonError := json.Marshal(response)
	if jsonError != nil {
		utils.InternalServerError(w, true, "Unexpected Error Occured")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	utils.Writer(w.Write(body))
}
