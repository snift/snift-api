package models

type Response struct {
	Score int `json:"score"`
	Messages []string `json:"messages"`
}

func NewResponse(score int, messages []string) *Response {
  response := &Response{
    Score: score,
    Messages:  messages,
  }
  return response
}