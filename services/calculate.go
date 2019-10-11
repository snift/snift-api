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

var badges []*models.Badge

// CalculateProtocolScore returns a score based on whether the protocol is http/https
func CalculateProtocolScore(protocol string) (score int) {
	if protocol == "https" {
		score = 5
		badges = append(badges, utils.GetHTTPSBadge())
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

// CalculateOverallScore returns the overall score for the specified URL
/** The following sub-scores are calculated to determine the overall score
 * Protocol Score
 * Response Headers Score
 * Mail Server Configuration Score
 **/
func CalculateOverallScore(scoresURL string) ([]byte, error) {
	var host string
	var port string
	badges = nil
	dbresponse := utils.FindEntry(scoresURL)
	if dbresponse != "" {
		return []byte(dbresponse), nil
	}
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
	var maximumPossibleScore = new(int)
	var calculatedScore = new(int)
	*maximumPossibleScore = 5 // why is this initialized to 5?

	protocolScore := CalculateProtocolScore(protocol)
	*calculatedScore += protocolScore

	responseHeaderScore, ServerDetail, ServerData, err := GetResponseHeaderScore(scoresURL)
	if err != nil {
		return nil, err
	}

	*maximumPossibleScore += responseHeaderScore.maximumValue
	*calculatedScore += responseHeaderScore.value

	mailServerScore, txtRecords, dmarcRecords := GetMailServerConfigurationScore(MailServerConfigParams{host, maximumPossibleScore})
	*calculatedScore += mailServerScore

	overallScore := math.Ceil((float64(float64(*calculatedScore)/float64(*maximumPossibleScore)))*100) / 100
	fmt.Println("Final Score for: " + scoresURL + " is " + strconv.Itoa(*calculatedScore) + " out of " + strconv.Itoa(*maximumPossibleScore))

	certificates, certError := models.GetCertificate(host, port, protocol)
	if certError != nil {
		return nil, certError
	}

	scores := models.GetScores(scoresURL, overallScore, badges)
	response := models.BuildScoresResponse(scores, certificates, nil, ServerDetail)
	responseBody, err := json.Marshal(response)
	serverdataJSON, serverdataJSONerr := json.Marshal(ServerData)
	if serverdataJSONerr != nil {
		fmt.Println("Error Occured while parsing Server Data JSON", serverdataJSONerr)
	}

	entry := &models.Domain{
		Name:         scoresURL,
		ServerData:   string(serverdataJSON),
		TxtRecords:   txtRecords,
		DmarcRecords: dmarcRecords,
		Response:     string(responseBody),
		IncidentList: "",
		Score:        overallScore,
	}
	utils.CreateEntry(entry)
	return responseBody, err
}

// HeaderScore represents a header score with value, meta and maxValue
type HeaderScore struct {
	value        int
	meta         string
	maximumValue int
}

// ResponseHeader returns a pointer to a the HeaderScore struct
type ResponseHeader func(hScore *HeaderScore) error

// BuildResponseHeaderScore builds header score
func BuildResponseHeaderScore(opts ...ResponseHeader) (*HeaderScore, error) {
	var hScore HeaderScore
	for _, opt := range opts {
		err := opt(&hScore)
		if err != nil {
			return nil, err
		}
		hScore.maximumValue += 5
	}
	return &hScore, nil
}

// GetResponseHeaderScore returns a cumulative score based on the response headers for the specified URL
func GetResponseHeaderScore(url string) (reponseHeaderScore HeaderScore, serverInfo *models.ServerDetail, serverData map[string]string, err error) {
	err = utils.IsValidURL(url)
	if err != nil {
		return reponseHeaderScore, nil, nil, err
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
		return reponseHeaderScore, nil, nil, err
	}
	responseHeaderMap = make(map[string]string)
	// Constructing Response Header Map
	for k, v := range response.Header {
		value := strings.Join(v, ",")
		responseHeaderMap[k] = value
	}
	// Calculating Scores for Individual Headers
	responseHeaderScore, err := BuildResponseHeaderScore(
		GetXSSScore(responseHeaderMap[XSSHeader]),
		GetXFrameScore(responseHeaderMap[XFrameHeader]),
		GetHSTSScore(responseHeaderMap[HSTSHeader]),
		GetCSPScore(responseHeaderMap[CSPHeader]),
		GetPKPScore(responseHeaderMap[PKPHeader]),
		GetReferrerPolicyScore(responseHeaderMap[RPHeader]),
		GetXContentTypeScore(responseHeaderMap[XContentTypeHeader]),
		GetHTTPVersionScore(response.Proto),
		GetTLSVersionScore(response.TLS),
	)

	serverInfo = getServerInformation(responseHeaderMap[Server])
	serverData = responseHeaderMap
	return *responseHeaderScore, serverInfo, serverData, err
}

// GetXSSScore returns the XSS Score of the URL
func GetXSSScore(XSSValue string) ResponseHeader {
	return func(xssHScore *HeaderScore) error {
		if XSSValue != "" {
			XSSValue = strings.TrimSpace(XSSValue)
			if XSSValue == XSSValues[0] {
				xssHScore.value += 0
			} else if strings.HasPrefix(XSSValue, XSSValues[1]) {
				xssHScore.value += 5
			}
			XSSValueReport := strings.Split(XSSValue, "report=")
			if len(XSSValueReport) == 2 {
				xssHScore.meta = XSSValueReport[1]
			}
		}
		return nil
	}
}

// GetXFrameScore returns the HTTP X-Frame-Options Response Header Score of the URL
func GetXFrameScore(XFrameValue string) ResponseHeader {
	return func(xFrameScore *HeaderScore) error {
		if XFrameValue != "" {
			XFrameValue = strings.TrimSpace(strings.ToLower(XFrameValue))
			if XFrameValue == XFrameValues[0] || XFrameValue == XFrameValues[1] {
				badges = append(badges, utils.GetXFrameBadge())
				xFrameScore.value += 5
			} else if strings.HasPrefix(XFrameValue, XFrameValues[2]) {
				xFrameScore.value += 4
			}
		} else {
			xFrameScore.value++
		}
		return nil
	}
}

// GetHSTSScore returns the HTTP Strict-Transport-Security Response Header Score of the URL
func GetHSTSScore(HSTS string) ResponseHeader {
	return func(hstsScore *HeaderScore) error {
		if HSTS != "" {
			if strings.HasPrefix(HSTS, HSTSValues[0]) {
				hstsScore.value += 4
				if strings.Contains(HSTS, HSTSValues[1]) || strings.Contains(HSTS, HSTSValues[2]) {
					badges = append(badges, utils.GetHSTSBadge())
					hstsScore.value++
				}
			}
		} else {
			hstsScore.value += 2
		}
		return nil
	}
}

// GetCSPScore returns the score for Content Security Policy Header
func GetCSPScore(CSP string) ResponseHeader {
	return func(cspScore *HeaderScore) error {
		if CSP != "" {
			badges = append(badges, utils.GetCSPBadge())
			cspScore.value += 5
		} else {
			cspScore.value += 3
		}
		return nil
	}
}

// GetPKPScore returns the score for Public Key Pinning Header
func GetPKPScore(PKP string) ResponseHeader {
	return func(pkpScore *HeaderScore) error {
		if PKP != "" {
			badges = append(badges, utils.GetHPKPBadge())
			pkpScore.value += 5
		} else {
			pkpScore.value += 3
		}
		return nil
	}
}

// GetReferrerPolicyScore returns the HTTP Referrer-Policy Response Header Score of the URL
func GetReferrerPolicyScore(ReferrerPolicy string) ResponseHeader {
	return func(xReferrerPolicyScore *HeaderScore) error {
		if ReferrerPolicy != "" {
			ReferrerPolicy = strings.TrimSpace(strings.ToLower(ReferrerPolicy))
			if score, ok := ReferrerPolicyValues[ReferrerPolicy]; ok {
				xReferrerPolicyScore.value += score
				if score >= 4 {
					badges = append(badges, utils.GetRPBadge())
				}
			}
		}
		return nil
	}
}

// GetXContentTypeScore returns the score for X-Content-Type-Options Header
func GetXContentTypeScore(XContentType string) ResponseHeader {
	return func(xContentTypeScore *HeaderScore) error {
		if XContentType == XContentTypeHeaderValue {
			badges = append(badges, utils.GetXContentTypeBadge())
			xContentTypeScore.value += 5
		}
		return nil
	}

}

// GetHTTPVersionScore returns the score for HTTP Version
func GetHTTPVersionScore(Proto string) ResponseHeader {
	return func(xHTTPVersionScore *HeaderScore) error {
		if Proto == HTTPVersion[0] {
			badges = append(badges, utils.GetHTTPVersionBadge())
			xHTTPVersionScore.value += 5
		} else if Proto == HTTPVersion[1] {
			xHTTPVersionScore.value += 2
		}
		return nil
	}

}

// GetTLSVersionScore returns the score for TLS Version
func GetTLSVersionScore(TLS *tls.ConnectionState) ResponseHeader {
	return func(xTLSVersionScore *HeaderScore) error {
		if TLS != nil {
			if TLS.Version == tls.VersionTLS12 {
				badges = append(badges, utils.GetTLSVersionBadge())
				xTLSVersionScore.value += 5
			} else if TLS.Version == tls.VersionTLS11 {
				xTLSVersionScore.value += 3
			} else if TLS.Version == tls.VersionTLS10 {
				xTLSVersionScore.value++
			}
		}
		return nil
	}
}

//MailServerConfigParams denotes args passed on to GetMailServerConfiguration
type MailServerConfigParams struct {
	host                 string
	maximumPossibleScore *int
}

// GetMailServerConfigurationScore returns the Mail Server Configuration Score of a Domain
func GetMailServerConfigurationScore(params MailServerConfigParams) (mailServerScore int, txtRecords string, dmarcRecord string) {
	mailServerScore = 0
	host := params.host
	maximumPossibleScore := params.maximumPossibleScore

	if strings.HasPrefix(host, "www.") {
		host = strings.Replace(host, "www.", "", -1)
	}

	spfScore, maxSPFScore, txtRecords := GetSPFScore(host)
	mailServerScore += spfScore

	dmarcScore, dmarcRecord := GetDMARCScore(host)
	mailServerScore += dmarcScore

	if maximumPossibleScore != nil {
		*maximumPossibleScore += maxSPFScore
		*maximumPossibleScore += 5
	}

	return
}

// GetSPFScore returns the Sender Policy Framework Score of the Domain
func GetSPFScore(domain string) (spfScore int, maxSPFScore int, txtRecords string) {
	command := strings.Replace(TXTQuery, "domain.com", domain, -1)
	out, err := exec.Command("bash", "-c", command).Output()
	txtRecords = string(out[:])

	if err != nil {
		fmt.Println("Unexpected Error Occured while extracting TXT Records", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(txtRecords))
	scanner.Split(bufio.ScanLines)

	spfRecordCount := 0
	spfScore = 0

	for scanner.Scan() {
		txtRecord := scanner.Text()
		// Removing Surrounding Quotes and trimming spaces
		txtRecord = strings.TrimSpace(txtRecord[1 : len(txtRecord)-1])
		if strings.HasSuffix(txtRecord, "-all") {
			spfScore += 5
			spfRecordCount++

		} else if strings.HasSuffix(txtRecord, "~all") {
			spfScore += 3
			spfRecordCount++

		} else if strings.HasSuffix(txtRecord, "?all") {
			spfScore += 2
			spfRecordCount++

		} else if strings.HasSuffix(txtRecord, "+all") {
			spfRecordCount++
		}
	}
	maxSPFScore = spfRecordCount * 5
	if spfScore == maxSPFScore {
		badges = append(badges, utils.GetSPFBadge())
	}
	return
}

// GetDMARCScore returns the DMARC Score of the Domain
func GetDMARCScore(domain string) (score int, dmarcRecord string) {
	command := strings.Replace(DMARCQuery, "domain.com", domain, -1)
	out, err := exec.Command("bash", "-c", command).Output()
	dmarcRecord = string(out[:])

	score = 0

	if err != nil {
		log.Fatal("Unexpected Error Occured while extracting DMARC Records ", err)
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
