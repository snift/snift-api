package service

import (
    "fmt"
    "strings"
    "strconv"
    models "snift-backend/models"
)

func CalculateScore(url string ) *models.Response {
	s := strings.Split(url, "://")
	var messages []string
	var score int = 0
	protocol := s[0]
	protocol_score,protocolMessage := CalculateProtocolScore(protocol)
	messages= append(messages,protocolMessage)
	score = score + protocol_score
	fmt.Println("Protocol Score is "+strconv.Itoa(protocol_score))
	fmt.Println("Conclusion: "+protocolMessage)
	fmt.Println("Final Score for: "+url+" is "+strconv.Itoa(score))
	response := models.NewResponse(score,messages)
	return response
}