package models

// ScoreResponse holds a Score JSON, the Certificate Details JSON for the main Scores API
type ScoreResponse struct {
	Scores *Scores `json:"scores"`
	Cert   *Cert   `json:"certificate_details"`
}

// GetScoresResponse returns a main Scores Response
func GetScoresResponse(scores *Scores, cert *Cert) *ScoreResponse {
	response := &ScoreResponse{
		Scores: scores,
		Cert:   cert,
	}
	return response
}
