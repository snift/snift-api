package models

// Scores holds a valid score, the incoming url and the outgoing message
type Scores struct {
	URL      string   `json:"url"`
	Score    int      `json:"score"`
	Messages []string `json:"messages"`
	Cert     *Cert    `json:"certificate_details"`
}

// GetScores returns a valid Score instance
func GetScores(url string, score int, messages []string, cert *Cert) *Scores {
	response := &Scores{
		URL:      url,
		Score:    score,
		Messages: messages,
		Cert:     cert,
	}
	return response
}
