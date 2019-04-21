package models

// Scores holds a valid score, the incoming url and the outgoing message
type Scores struct {
	URL    string   `json:"url"`
	Score  float64  `json:"score"`
	Badges []*Badge `json:"badges"`
}

// ScoresRequest holds the structure for Scores API Request Body
type ScoresRequest struct {
	URL string `json:"url"`
}

// GetScores returns a valid Score instance
func GetScores(url string, score float64, badges []*Badge) *Scores {
	response := &Scores{
		URL:    url,
		Score:  score,
		Badges: badges,
	}
	return response
}
