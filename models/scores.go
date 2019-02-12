package models

// Scores holds a valid score, the incoming url and the outgoing message
type Scores struct {
	URL      string   `json:"url"`
	Score    int      `json:"score"`
	Messages []string `json:"messages"`
}

// GetScores returns a valid Score instance
func GetScores(url string, score int, messages []string) *Scores {
	response := &Scores{
		URL:      url,
		Score:    score,
		Messages: messages,
	}
	return response
}
