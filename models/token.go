package models

// Token stores the JWT Token and Expiry Time
type Token struct {
	Token      string `json:"token"`
	ExpiryTime int64  `json:"expiry_time"`
}
