package controllers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestHomePage(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HomePage)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, rr.Body.String(), "Welcome to Snift!")

}

func TestInvalidURL(t *testing.T) {
	req, _ := http.NewRequest("GET", "/scores?url=example", nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetScore)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusBadRequest)
	assert.Equal(t, rr.Body.String(), "{\"error\":\"Invalid URL\"}")

}
func TestValidURL(t *testing.T) {
	req, _ := http.NewRequest("GET", "/scores?url=https://example.com", nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetScore)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)

}
