package utils

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

// Writer checks and validates the response
func Writer(n int, err error) {
	if err != nil {
		log.Fatal("Error Occur while writing response: ", err)
	}
}

// BadRequest returns error JSON for 400 Bad Request
func BadRequest(w http.ResponseWriter, isJSON bool, err string) {
	if !isJSON {
		http.Error(w, err, http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, `{"error":%q}`, err)
}

// InternalServerError returns error JSON for InternalServerError
func InternalServerError(w http.ResponseWriter, isJSON bool, err string) {
	if !isJSON {
		http.Error(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, `{"error":%q}`, err)
}

// Unauthorized returns error JSON for Unauthorized Error
func Unauthorized(w http.ResponseWriter, isJSON bool, err string) {
	if !isJSON {
		http.Error(w, err, http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprintf(w, `{"error":%q}`, err)
}

// IsValidURL tests a string to determine if it is a url or not.
func IsValidURL(rawURL string) error {
	_, err := url.ParseRequestURI(rawURL)
	return err
}

// GetAccessControlAllowOrigin returns the value of Access-Control-Allow-Origin Header
func GetAccessControlAllowOrigin() string {
	return os.Getenv("ACCESS-CONTROL-ALLOW-ORIGIN")
}
