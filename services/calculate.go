package services

import (
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"net/url"
	models "snift-backend/models"
	"snift-backend/utils"
	"strconv"
	"strings"
)

// XSSHeader has the XSS Header Name
const XSSHeader = "X-Xss-Protection"

// XFrameHeader has the XFrame Header Name
const XFrameHeader = "X-Frame-Options"

// HSTSHeader has the HSTS Header Name
const HSTSHeader = "Strict-Transport-Security"

// CSPHeader has the CSP Header Name
const CSPHeader = "Content-Security-Policy"

// PKPHeader has the PKP Header Name
const PKPHeader = "Public-Key-Pins"

// RPHeader has the RP Header Name
const RPHeader = "Referrer-Policy"

// XSSValues is used to store the X-Xss-Protection Header values
var XSSValues = [...]string{"0", "1"}

// XFrameValues is used to store the X-Frame-Options Header values
var XFrameValues = [...]string{"deny", "sameorigin", "allow-from"}

// HSTSValues used to store the X-Frame-Options Header values
var HSTSValues = [...]string{"max-age", "includeSubDomains", "preload"}

// ReferrerPolicyValues used to store the Referrer-Policy Header values
var ReferrerPolicyValues = [...]string{"no-referrer", "no-referrer-when-downgrade", "origin", "origin-when-cross-origin", "same-origin", "strict-origin", "strict-origin-when-cross-origin", "unsafe-url"}

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
	protocolScore, protocolMessage := CalculateProtocolScore(protocol)
	messages = append(messages, protocolMessage)
	score = score + protocolScore
	var maximumScore = 5
	headerScore, _, maxScore, err := GetResponseHeaderScore(scoresURL)
	if err != nil {
		return nil, err
	}
	maximumScore = maximumScore + maxScore
	score = score + headerScore
	totalScore := math.Ceil((float64(float64(score)/float64(maximumScore)))*100) / 100
	fmt.Println("Protocol Score is " + strconv.Itoa(protocolScore))
	fmt.Println("Message: " + protocolMessage)
	fmt.Println("Final Score for: " + scoresURL + " is " + strconv.Itoa(score) + " out of " + strconv.Itoa(maximumScore))
	certificates, certError := models.GetCertificate(host, port, protocol)
	if certError != nil {
		return nil, certError
	}
	scores := models.GetScores(scoresURL, totalScore, messages)
	response := models.GetScoresResponse(scores, certificates)
	return response, nil
}

// GetResponseHeaderScore returns the Response Header Score for the HTTP Request
func GetResponseHeaderScore(url string) (totalScore int, XSSReportURL string, maxScore int, err error) {
	err = utils.IsValidURL(url)
	if err != nil {
		return 0, "", 0, err
	}
	var responseHeaderMap map[string]string
	response, err := http.Head(url)
	log.Print(err)
	responseHeaderMap = make(map[string]string)
	for k, v := range response.Header {
		value := strings.Join(v, ",")
		responseHeaderMap[k] = value
	}
	totalScore = 0
	var score = 0
	maxScore = 0
	if val, ok := responseHeaderMap[XSSHeader]; ok {
		score, XSSReportURL = GetXSSScore(val)
	}
	maxScore = maxScore + 5
	totalScore = totalScore + score
	score = 1
	if val, ok := responseHeaderMap[XFrameHeader]; ok {
		score = GetXFrameScore(val)
	}
	maxScore = maxScore + 5
	totalScore = totalScore + score
	score = 2
	if val, ok := responseHeaderMap[HSTSHeader]; ok {
		score = GetHSTSScore(val)
	}
	maxScore = maxScore + 5
	totalScore = totalScore + score
	score = 3
	if _, ok := responseHeaderMap[CSPHeader]; ok {
		score = 5
	}
	maxScore = maxScore + 5
	totalScore = totalScore + score
	score = 3
	if _, ok := responseHeaderMap[PKPHeader]; ok {
		score = 5
	}
	maxScore = maxScore + 5
	totalScore = totalScore + score
	score = 2
	if val, ok := responseHeaderMap[RPHeader]; ok {
		score = GetReferrerPolicyScore(val)
	}
	maxScore = maxScore + 5
	totalScore = totalScore + score
	return
}

// GetXSSScore returns the XSS Score of the URL
func GetXSSScore(XSSValue string) (score int, XSSReportURL string) {
	XSSValue = strings.TrimSpace(XSSValue)
	if strings.Compare(XSSValue, XSSValues[0]) == 0 {
		score = 0
	} else if strings.HasPrefix(XSSValue, XSSValues[1]) {
		score = 5
	}
	XSSValueReport := strings.Split(XSSValue, "report=")
	if len(XSSValueReport) == 2 {
		XSSReportURL = XSSValueReport[1]
	}
	return
}

// GetXFrameScore returns the HTTP X-Frame-Options Response Header Score of the URL
func GetXFrameScore(XFrameValue string) (score int) {
	XFrameValue = strings.TrimSpace(strings.ToLower(XFrameValue))
	if strings.Compare(XFrameValue, XFrameValues[0]) == 0 || strings.Compare(XFrameValue, XFrameValues[1]) == 0 {
		score = 5
	} else if strings.HasPrefix(XFrameValue, XFrameValues[2]) {
		score = 4
	}
	return
}

// GetHSTSScore returns the HTTP Strict-Transport-Security Response Header Score of the URL
func GetHSTSScore(HSTS string) (score int) {
	if strings.HasPrefix(HSTS, HSTSValues[0]) {
		score = 4
		if strings.Contains(HSTS, HSTSValues[1]) || strings.Contains(HSTS, HSTSValues[2]) {
			score = 5
		}
	}
	return
}

// GetReferrerPolicyScore returns the HTTP Referrer-Policy Response Header Score of the URL
func GetReferrerPolicyScore(ReferrerPolicy string) (score int) {
	ReferrerPolicy = strings.TrimSpace(strings.ToLower(ReferrerPolicy))
	if strings.Compare(ReferrerPolicy, ReferrerPolicyValues[0]) == 0 {
		score = 5
	} else if strings.Compare(ReferrerPolicy, ReferrerPolicyValues[1]) == 0 || strings.Compare(ReferrerPolicy, ReferrerPolicyValues[2]) == 0 || strings.Compare(ReferrerPolicy, ReferrerPolicyValues[3]) == 0 || strings.Compare(ReferrerPolicy, ReferrerPolicyValues[4]) == 0 || strings.Compare(ReferrerPolicy, ReferrerPolicyValues[5]) == 0 || strings.Compare(ReferrerPolicy, ReferrerPolicyValues[6]) == 0 {
		score = 4
	} else if strings.Compare(ReferrerPolicy, ReferrerPolicyValues[7]) == 0 {
		score = 2
	}
	return
}
