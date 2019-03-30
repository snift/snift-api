package models

// WebServer contains the information extracted from web_servers.json
type WebServer struct {
	Prefix       string        `json:"starts_with"`
	ServerDetail *ServerDetail `json:"details"`
}

// ServerDetail contains the information about the WebServer
type ServerDetail struct {
	Name         string `json:"name"`
	Developer    string `json:"developed_by"`
	License      string `json:"license"`
	Website      string `json:"website"`
	IsOpenSource bool   `json:"is_open_source"`
}
