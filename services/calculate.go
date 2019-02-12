package services

import (
	"fmt"
	models "snift-backend/models"
	"strconv"
	"strings"
)

// CalculateProtocolScore returns a score based on whether the protocol is http/https
func CalculateProtocolScore(protocol string) (score int, message string) {
	score = -1
	if strings.Compare(protocol, "http") == 0 {
		score = 0
		message = "Website is unencrypted and hence subjective to Man-in-the-Middle attacks(MITM) and Eavesdropping Attacks."
	} else if strings.Compare(protocol, "https") == 0 {
		score = 5
		message = "From the protocol level, Website is secure."
	} else {
		message = "Protocol Not Found"
	}
	return
}

// CalculateOverallScore returns the overall score for the incoming request
func CalculateOverallScore(url string) *models.ScoreResponse {
	s := strings.Split(url, "://")
	var messages []string
	var score int
	protocol, host := s[0], s[1]
	protocolScore, protocolMessage := CalculateProtocolScore(protocol)
	messages = append(messages, protocolMessage)
	score = score + protocolScore
	fmt.Println("Protocol Score is " + strconv.Itoa(protocolScore))
	fmt.Println("Message: " + protocolMessage)
	fmt.Println("Final Score for: " + url + " is " + strconv.Itoa(score))
	response := models.GetScoresResponse(models.GetScores(url, score, messages), models.GetCertificate(host))
	return response
}
