package utils

import (
	"fmt"
	"os"
	"snift-api/models"

	"github.com/jinzhu/gorm"
	// Import for Postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func initConnection() (db *gorm.DB, err error) {
	db, err = gorm.Open("postgres", getDBURL())
	return
}

func getDBHost() string {
	return os.Getenv("DB.HOST")
}

func getDBPort() string {
	return os.Getenv("DB.PORT")
}

func getDBUser() string {
	return os.Getenv("DB.USER")
}

func getDBName() string {
	return os.Getenv("DB.NAME")
}

func getDBPassword() string {
	return os.Getenv("DB.PASSWORD")
}

func getDBURL() string {
	return "host=" + getDBHost() + " port=" + getDBPort() + " user=" + getDBUser() + " dbname=" + getDBName() + " password=" + getDBPassword()

}

// CreateEntry is used to create entry to the table
func CreateEntry(domain *models.Domain) {
	db, err := initConnection()
	if err != nil {
		fmt.Println("Error Occured while initializing connection", err)
		return
	}
	db.AutoMigrate(&models.Domain{})
	db.Create(domain)
	defer db.Close()

}

// FindEntry is used to find entry from the table
func FindEntry(url string) (response string) {
	var domain models.Domain
	db, err := initConnection()
	if err != nil {
		fmt.Println("Error Occured while initializing connection", err)
		return
	}
	db.Where("name = ?", url).First(&domain)
	response = domain.Response
	defer db.Close()
	return
}
