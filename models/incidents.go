package models

import "encoding/xml"

// Incidents contains the list of Security Vulnerabilities Reported
type Incidents struct {
	XMLName      xml.Name   `xml:"incidents" json:"-"`
	IncidentList []Incident `xml:"item" json:"incidents"`
}

// Incident contains the properties of a individual Security Vulnerability
type Incident struct {
	XMLName       xml.Name `xml:"item" json:"-"`
	URL           string   `xml:"url" json:"url"`
	Host          string   `xml:"host" json:"host"`
	Type          string   `xml:"type" json:"type"`
	ReportedDate  string   `xml:"reporteddate" json:"reported_date"`
	Researcher    string   `xml:"researcher" json:"researcher"`
	ResearcherURL string   `xml:"researcherurl" json:"researcher_url"`
	Fixed         bool     `xml:"fixed" json:"fixed"`
	FixedDate     string   `xml:"fixeddate" json:"fixed_date"`
}
