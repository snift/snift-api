package utils

import (
	"net"
	"net/http"
	"os"
	"snift-api/models"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GetToken creates a JWT Token and returns it
func GetToken(r *http.Request) (response *models.Token, err error) {
	expiryTime := time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ip":  r.RemoteAddr,
		"ua":  r.Header.Get("User-Agent"),
		"exp": expiryTime,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	response = &models.Token{
		Token:      tokenString,
		ExpiryTime: expiryTime,
	}
	return
}

// ValidateToken is used to check if the token is valid or not
func ValidateToken(r *http.Request) bool {
	var host string
	var actualHost string
	AuthToken := r.Header.Get("X-Auth-Token")

	if AuthToken == "" {
		return false
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(AuthToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})

	if strings.Contains(r.RemoteAddr, ":") {
		host, _, _ = net.SplitHostPort(r.RemoteAddr)
	} else {
		host = r.RemoteAddr
	}

	if strings.Contains(claims["ip"].(string), ":") {
		actualHost, _, _ = net.SplitHostPort(claims["ip"].(string))
	} else {
		actualHost = claims["ip"].(string)
	}

	if err == nil && token.Valid && host == actualHost && claims["exp"].(float64) > float64(time.Now().Unix()) {
		return true
	}

	return false
}
