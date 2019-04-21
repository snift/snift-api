package models

import (
	"github.com/jinzhu/gorm"
)

// Domain has domain details that is stored in the database
type Domain struct {
	gorm.Model
	Name         string `gorm:"size:255;unique"`
	ServerData   string `gorm:"size:4095"`
	Response     string `gorm:"size:8191"`
	TxtRecords   string `gorm:"size:2047"`
	DmarcRecords string `gorm:"size:255"`
	IncidentList string `gorm:"size:4095"`
	Score        float64
}
