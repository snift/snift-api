package services

import (
	"fmt"
	"net"
	"net/url"
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

var getDefaultPort = func(protocol string) string {
	// default http port
	port := "80"
	// default https port
	if protocol == "https" {
		port = "443"
	}
	return port
}

// CalculateOverallScore returns the overall score for the incoming request
func CalculateOverallScore(scoresURL string) (*models.ScoreResponse, error) {
	var messages []string
	var score int
	var host string
	var port string
	domain, err := url.Parse(scoresURL)
	if err != nil {
		return nil, err
	}
	protocol := domain.Scheme
	if strings.Contains(domain.Host, ":") {
		host, port, _ = net.SplitHostPort(domain.Host)
	} else {
		host = domain.Host
	}
	if port == "" {
		port = getDefaultPort(protocol)
	}
	fmt.Println(host + port)
	protocolScore, protocolMessage := CalculateProtocolScore(protocol)
	messages = append(messages, protocolMessage)
	score = score + protocolScore
	fmt.Println("Protocol Score is " + strconv.Itoa(protocolScore))
	fmt.Println("Message: " + protocolMessage)
	fmt.Println("Final Score for: " + scoresURL + " is " + strconv.Itoa(score))
	certificates, certError := models.GetCertificate(host, port, protocol)
	if certError != nil {
		return nil, certError
	}
	scores := models.GetScores(scoresURL, score, messages)
	response := models.GetScoresResponse(scores, certificates)
	return response, nil
}
