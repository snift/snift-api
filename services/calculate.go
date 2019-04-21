package services

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"snift-api/models"
	"snift-api/utils"
	"strconv"
	"strings"
	"time"
)

// CalculateProtocolScore returns a score based on whether the protocol is http/https
func CalculateProtocolScore(protocol string) (score int, message string) {
	score = -1
	if protocol == "http" {
		score = 0
		message = "Website is unencrypted and hence subjective to Man-in-the-Middle attacks(MITM) and Eavesdropping Attacks."
	} else if protocol == "https" {
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
		fmt.Println(err)
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
	headerScore, _, maxScore, ServerDetail, err := GetResponseHeaderScore(scoresURL)
	if err != nil {
		fmt.Println("hello", err)
		return nil, err
	}
	maximumScore = maximumScore + maxScore
	score = score + headerScore
	mailServerScore, maxScore := GetMailServerConfigurationScore(host)
	score += mailServerScore
	maximumScore += maxScore
	totalScore := math.Ceil((float64(float64(score)/float64(maximumScore)))*100) / 100
	fmt.Println("Protocol Score is " + strconv.Itoa(protocolScore))
	fmt.Println("Message: " + protocolMessage)
	fmt.Println("Final Score for: " + scoresURL + " is " + strconv.Itoa(score) + " out of " + strconv.Itoa(maximumScore))
	certificates, certError := models.GetCertificate(host, port, protocol)
	if certError != nil {
		return nil, certError
	}
	scores := models.GetScores(scoresURL, totalScore, messages)
	response := models.GetScoresResponse(scores, certificates, nil, ServerDetail)
	return response, nil
}

// GetResponseHeaderScore returns the Response Header Score for the HTTP Request
func GetResponseHeaderScore(url string) (totalScore int, XSSReportURL string, maxScore int, serverInfo *models.ServerDetail, err error) {
	err = utils.IsValidURL(url)
	if err != nil {
		return 0, "", 0, nil, err
	}
	var responseHeaderMap map[string]string
	// Initializing client to avoid Redirects
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}
	response, err := client.Head(url)
	if err != nil {
		fmt.Println(err)
		return 0, "", 0, nil, err
	}
	responseHeaderMap = make(map[string]string)
	// Constructing Response Header Map
	for k, v := range response.Header {
		value := strings.Join(v, ",")
		responseHeaderMap[k] = value
	}
	// Calculating Scores for Individual Headers
	totalScore, maxScore, XSSReportURL = GetXSSScore(responseHeaderMap[XSSHeader], totalScore, maxScore)
	totalScore, maxScore = GetXFrameScore(responseHeaderMap[XFrameHeader], totalScore, maxScore)
	totalScore, maxScore = GetHSTSScore(responseHeaderMap[HSTSHeader], totalScore, maxScore)
	totalScore, maxScore = GetCSPScore(responseHeaderMap[CSPHeader], totalScore, maxScore)
	totalScore, maxScore = GetPKPScore(responseHeaderMap[PKPHeader], totalScore, maxScore)
	totalScore, maxScore = GetReferrerPolicyScore(responseHeaderMap[RPHeader], totalScore, maxScore)
	totalScore, maxScore = GetXContentTypeScore(responseHeaderMap[XContentTypeHeader], totalScore, maxScore)
	totalScore, maxScore = GetHTTPVersionScore(response.Proto, totalScore, maxScore)
	totalScore, maxScore = GetTLSVersionScore(response.TLS, totalScore, maxScore)
	serverInfo = getServerInformation(responseHeaderMap[Server])
	return
}

// GetXSSScore returns the XSS Score of the URL
func GetXSSScore(XSSValue string, totalScore int, maxScore int) (int, int, string) {
	var XSSReportURL string
	maxScore += 5
	if XSSValue != "" {
		XSSValue = strings.TrimSpace(XSSValue)
		if XSSValue == XSSValues[0] {
			totalScore += 0
		} else if strings.HasPrefix(XSSValue, XSSValues[1]) {
			totalScore += 5
		}
		XSSValueReport := strings.Split(XSSValue, "report=")
		if len(XSSValueReport) == 2 {
			XSSReportURL = XSSValueReport[1]
		}
	}

	return totalScore, maxScore, XSSReportURL
}

// GetXFrameScore returns the HTTP X-Frame-Options Response Header Score of the URL
func GetXFrameScore(XFrameValue string, totalScore int, maxScore int) (int, int) {
	maxScore += 5
	if XFrameValue != "" {
		XFrameValue = strings.TrimSpace(strings.ToLower(XFrameValue))
		if XFrameValue == XFrameValues[0] || XFrameValue == XFrameValues[1] {
			totalScore += 5
		} else if strings.HasPrefix(XFrameValue, XFrameValues[2]) {
			totalScore += 4
		}
	} else {
		totalScore++
	}
	return totalScore, maxScore
}

// GetHSTSScore returns the HTTP Strict-Transport-Security Response Header Score of the URL
func GetHSTSScore(HSTS string, totalScore int, maxScore int) (int, int) {
	maxScore += 5
	if HSTS != "" {
		if strings.HasPrefix(HSTS, HSTSValues[0]) {
			totalScore += 4
			if strings.Contains(HSTS, HSTSValues[1]) || strings.Contains(HSTS, HSTSValues[2]) {
				totalScore++
			}
		}
	} else {
		totalScore += 2
	}
	return totalScore, maxScore
}

// GetCSPScore returns the score for Content Security Policy Header
func GetCSPScore(CSP string, totalScore int, maxScore int) (int, int) {
	maxScore += 5
	if CSP != "" {
		totalScore += 5
	} else {
		totalScore += 3
	}
	return totalScore, maxScore
}

// GetPKPScore returns the score for Public Key Pinning Header
func GetPKPScore(PKP string, totalScore int, maxScore int) (int, int) {
	maxScore += 5
	if PKP != "" {
		totalScore += 5
	} else {
		totalScore += 3
	}
	return totalScore, maxScore
}

// GetReferrerPolicyScore returns the HTTP Referrer-Policy Response Header Score of the URL
func GetReferrerPolicyScore(ReferrerPolicy string, totalScore int, maxScore int) (int, int) {
	maxScore += 5
	if ReferrerPolicy != "" {
		ReferrerPolicy = strings.TrimSpace(strings.ToLower(ReferrerPolicy))
		if score, ok := ReferrerPolicyValues[ReferrerPolicy]; ok {
			totalScore += score
		}
	}
	return totalScore, maxScore
}

// GetXContentTypeScore returns the score for X-Content-Type-Options Header
func GetXContentTypeScore(XContentType string, totalScore int, maxScore int) (int, int) {
	maxScore += 5
	if XContentType == XContentTypeHeaderValue {
		totalScore += 5
	}
	return totalScore, maxScore
}

// GetHTTPVersionScore returns the score for HTTP Version
func GetHTTPVersionScore(Proto string, totalScore int, maxScore int) (int, int) {
	maxScore += 5
	if Proto == HTTPVersion[0] {
		totalScore += 5
	} else if Proto == HTTPVersion[1] {
		totalScore += 2
	}
	return totalScore, maxScore
}

// GetTLSVersionScore returns the score for TLS Version
func GetTLSVersionScore(TLS *tls.ConnectionState, totalScore int, maxScore int) (int, int) {
	maxScore += 5
	if TLS != nil {
		if TLS.Version == tls.VersionTLS12 {
			totalScore += 5
		} else if TLS.Version == tls.VersionTLS11 {
			totalScore += 3
		} else if TLS.Version == tls.VersionTLS10 {
			totalScore++
		}
	}
	return totalScore, maxScore
}

// GetMailServerConfigurationScore returns the Mail Server Configuration Score of a Domain
func GetMailServerConfigurationScore(host string) (totalScore int, maximumScore int) {
	maximumScore = 0
	totalScore = 0
	if strings.HasPrefix(host, "www.") {
		host = strings.Replace(host, "www.", "", -1)
	}
	spfScore, maxScore := GetSPFScore(host)
	maximumScore = maximumScore + maxScore
	totalScore = totalScore + spfScore
	totalScore += GetDMARCScore(host)
	maximumScore += 5
	return
}

// GetSPFScore returns the Sender Policy Framework Score of the Domain
func GetSPFScore(domain string) (totalScore int, maxScore int) {
	command := strings.Replace(TXTQuery, "domain.com", domain, -1)
	out, err := exec.Command("bash", "-c", command).Output()
	txtRecords := string(out[:])

	if err != nil {
		fmt.Println("Unexpected Error Occured while extracting TXT Records", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(txtRecords))
	scanner.Split(bufio.ScanLines)

	spfRecordCount := 0
	totalScore = 0

	for scanner.Scan() {
		txtRecord := scanner.Text()
		// Removing Surrounding Quotes and trimming spaces
		txtRecord = strings.TrimSpace(txtRecord[1 : len(txtRecord)-1])
		if strings.HasSuffix(txtRecord, "-all") {
			totalScore = totalScore + 5
			spfRecordCount++

		} else if strings.HasSuffix(txtRecord, "~all") {
			totalScore = totalScore + 3
			spfRecordCount++

		} else if strings.HasSuffix(txtRecord, "?all") {
			totalScore = totalScore + 2
			spfRecordCount++

		} else if strings.HasSuffix(txtRecord, "+all") {
			spfRecordCount++
		}
	}
	maxScore = spfRecordCount * 5
	return
}

// GetDMARCScore returns the DMARC Score of the Domain
func GetDMARCScore(domain string) (score int) {
	command := strings.Replace(DMARCQuery, "domain.com", domain, -1)
	out, err := exec.Command("bash", "-c", command).Output()
	dmarcRecord := string(out[:])

	score = 0

	if err != nil {
		log.Fatal("Unexpected Error Occured while extracting DMARC Records", err)
	}
	if len(dmarcRecord) > 2 {
		dmarcRecord = strings.TrimSpace(dmarcRecord[1 : len(dmarcRecord)-1])
		if strings.HasPrefix(dmarcRecord, "v=DMARC") {
			score = 5
		}
	}
	return
}

// GetPreviousVulnerabilitiesScore gets the score for Previous Vulnerabilities taken from openbugbounty.org
func GetPreviousVulnerabilitiesScore(host string) (totalScore int, maxScore int, IncidentList []models.Incident) {
	if strings.HasPrefix(host, "www.") {
		host = strings.Replace(host, "www.", "", -1)
	}
	resp, err := http.Get(OpenBugBountyURL + host)
	if err != nil {
		log.Fatalln("Error Occured while sending Request ", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error Occured while reading HTTP Response ", err)
	}

	var incidents models.Incidents
	err = xml.Unmarshal(body, &incidents)
	if err != nil {
		log.Fatalln("Error Occured while Unmarshalling XML Response", err)
	}
	maxScore = len(incidents.IncidentList) * 10
	totalScore = 0
	for _, incident := range incidents.IncidentList {
		if incident.Fixed {
			ReportedDate, _ := time.Parse(time.RFC1123Z, incident.ReportedDate)
			FixedDate, _ := time.Parse(time.RFC1123Z, incident.FixedDate)
			diff := FixedDate.Sub(ReportedDate)
			if diff.Hours() > MaxIncidentResponseTime {
				totalScore += 5
			} else {
				totalScore += 10
			}
		}
	}
	return totalScore, maxScore, incidents.IncidentList
}

func getServerInformation(server string) (serverInfo *models.ServerDetail) {
	if server == "" {
		return
	}
	jsonFile, err := os.Open("resources/web_servers.json")
	if err != nil {
		log.Fatal("Error Occured while opening JSON File ", err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	jsonValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		log.Fatal("Error Occured while reading JSON File ", err)
	}

	values := make([]models.WebServer, 0)
	err = json.Unmarshal(jsonValue, &values)
	if err != nil {
		log.Fatal("Error Occured while parsing JSON ", err)
	}

	for _, serverValue := range values {
		if strings.HasPrefix(server, serverValue.Prefix) {
			serverInfo = serverValue.ServerDetail
			break
		}
	}
	return
}
