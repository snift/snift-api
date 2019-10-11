package models

// ScoresResponse holds a Score JSON, the Certificate Details JSON for the main Scores API
type ScoresResponse struct {
	Scores       *Scores       `json:"scores"`
	Cert         *Cert         `json:"certificate_details,omitempty"`
	IncidentList []Incident    `json:"security_incidents,omitempty"`
	ServerDetail *ServerDetail `json:"web_server,omitempty"`
}

// BuildScoresResponse builds the final api response for /score
func BuildScoresResponse(scores *Scores, cert *Cert, IncidentList []Incident, ServerDetail *ServerDetail) *ScoresResponse {
	response := &ScoresResponse{
		Scores:       scores,
		Cert:         cert,
		IncidentList: IncidentList,
		ServerDetail: ServerDetail,
	}
	return response
}
