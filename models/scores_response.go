package models

// ScoreResponse holds a Score JSON, the Certificate Details JSON for the main Scores API
type ScoreResponse struct {
	Scores       *Scores       `json:"scores"`
	Cert         *Cert         `json:"certificate_details"`
	IncidentList []Incident    `json:"security_incidents"`
	ServerDetail *ServerDetail `json:"web_server"`
}

// GetScoresResponse returns a main Scores Response
func GetScoresResponse(scores *Scores, cert *Cert, IncidentList []Incident, ServerDetail *ServerDetail) *ScoreResponse {
	response := &ScoreResponse{
		Scores:       scores,
		Cert:         cert,
		IncidentList: IncidentList,
		ServerDetail: ServerDetail,
	}
	return response
}
