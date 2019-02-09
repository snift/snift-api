package controllers

import (
	"net/http"
	"fmt"
	"log"
	"encoding/json"
	service "snift-backend/service"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Snift!")
	fmt.Println("Endpoint Hit: homePage")
}

func GetScore(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Query().Get("url")
	log.Print("Endpoint Hit: return Score")
	//fmt.Fprintf(w, "The URL is "+url )
	response := service.CalculateScore(url)
	json.NewEncoder(w).Encode(response)

}