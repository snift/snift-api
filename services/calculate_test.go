package services

import (
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

func TestCalculateProtocolScore(t *testing.T) {
	score := CalculateProtocolScore("http")
	assert.Equal(t, score, 0)
	score = CalculateProtocolScore("https")
	assert.Equal(t, score, 5)
}

func TestGetDefaultPort(t *testing.T) {
	assert.Equal(t, getDefaultPort("http"), "80")
	assert.Equal(t, getDefaultPort("https"), "443")
}

func TestCalculateOverallScore(t *testing.T) {
	_, err := CalculateOverallScore("http://example.com")
	assert.NoError(t, err)
	_, err = CalculateOverallScore("example.com")
	assert.Error(t, err)
}

func TestGetXSSScore(t *testing.T) {
	totalScore, maxScore, XSSReportURL := GetXSSScore("0", 0, 0)
	assert.Equal(t, totalScore, 0)
	assert.Equal(t, maxScore, 5)
	assert.Equal(t, XSSReportURL, "")
	totalScore, maxScore, XSSReportURL = GetXSSScore("1", 0, 0)
	assert.Equal(t, totalScore, 5)
	assert.Equal(t, maxScore, 5)
	assert.Equal(t, XSSReportURL, "")
	totalScore, maxScore, XSSReportURL = GetXSSScore("1; report=https://www.example.com", 0, 0)
	assert.Equal(t, totalScore, 5)
	assert.Equal(t, maxScore, 5)
	assert.Equal(t, XSSReportURL, "https://www.example.com")
	totalScore, maxScore, XSSReportURL = GetXSSScore("1; mode=block; report=https://www.example.com", 0, 0)
	assert.Equal(t, totalScore, 5)
	assert.Equal(t, maxScore, 5)
	assert.Equal(t, XSSReportURL, "https://www.example.com")
	totalScore, maxScore, XSSReportURL = GetXSSScore("", 0, 0)
	assert.Equal(t, totalScore, 0)
	assert.Equal(t, maxScore, 5)
	assert.Equal(t, XSSReportURL, "")
}

func TestGetXFrameScore(t *testing.T) {
	totalScore, maxScore := GetXFrameScore("DENY", 0, 0)
	assert.Equal(t, totalScore, 5)
	assert.Equal(t, maxScore, 5)
	totalScore, maxScore = GetXFrameScore("deny", 0, 0)
	assert.Equal(t, totalScore, 5)
	assert.Equal(t, maxScore, 5)
	totalScore, maxScore = GetXFrameScore("SAMEORIGIN", 0, 0)
	assert.Equal(t, totalScore, 5)
	assert.Equal(t, maxScore, 5)
	totalScore, maxScore = GetXFrameScore("sameorigin", 0, 0)
	assert.Equal(t, totalScore, 5)
	assert.Equal(t, maxScore, 5)
	totalScore, maxScore = GetXFrameScore("ALLOW-FROM https://www.example.com", 0, 0)
	assert.Equal(t, totalScore, 4)
	assert.Equal(t, maxScore, 5)
	totalScore, maxScore = GetXFrameScore("allow-from https://www.example.com", 0, 0)
	assert.Equal(t, totalScore, 4)
	assert.Equal(t, maxScore, 5)
	totalScore, maxScore = GetXFrameScore("", 0, 0)
	assert.Equal(t, totalScore, 1)
	assert.Equal(t, maxScore, 5)
}

func TestGetHSTSScore(t *testing.T) {
	totalScore, maxScore := GetHSTSScore("max-age=65536", 0, 0)
	assert.Equal(t, totalScore, 4)
	assert.Equal(t, maxScore, 5)
	totalScore, maxScore = GetHSTSScore("max-age=31536000; includeSubDomains", 0, 0)
	assert.Equal(t, totalScore, 5)
	assert.Equal(t, maxScore, 5)
	totalScore, maxScore = GetHSTSScore("max-age=31536000; includeSubDomains; preload", 0, 0)
	assert.Equal(t, totalScore, 5)
	assert.Equal(t, maxScore, 5)
	totalScore, maxScore = GetHSTSScore("max-age=31536;  preload", 0, 0)
	assert.Equal(t, totalScore, 5)
	assert.Equal(t, maxScore, 5)
	totalScore, maxScore = GetHSTSScore("", 0, 0)
	assert.Equal(t, totalScore, 2)
	assert.Equal(t, maxScore, 5)
}

func TestGetReferrerPolicyScore(t *testing.T) {
	totalScore, maxScore := GetReferrerPolicyScore("no-referrer", 0, 0)
	assert.Equal(t, totalScore, 5)
	assert.Equal(t, maxScore, 5)
	totalScore, maxScore = GetReferrerPolicyScore("no-referrer-when-downgrade", 0, 0)
	assert.Equal(t, totalScore, 4)
	assert.Equal(t, maxScore, 5)
	totalScore, maxScore = GetReferrerPolicyScore("origin", 0, 0)
	assert.Equal(t, totalScore, 4)
	assert.Equal(t, maxScore, 5)
	totalScore, maxScore = GetReferrerPolicyScore("origin-when-cross-origin", 0, 0)
	assert.Equal(t, totalScore, 4)
	assert.Equal(t, maxScore, 5)
	totalScore, maxScore = GetReferrerPolicyScore("same-origin", 0, 0)
	assert.Equal(t, totalScore, 4)
	assert.Equal(t, maxScore, 5)
	totalScore, maxScore = GetReferrerPolicyScore("strict-origin", 0, 0)
	assert.Equal(t, totalScore, 4)
	assert.Equal(t, maxScore, 5)
	totalScore, maxScore = GetReferrerPolicyScore("strict-origin-when-cross-origin", 0, 0)
	assert.Equal(t, totalScore, 4)
	assert.Equal(t, maxScore, 5)
	totalScore, maxScore = GetReferrerPolicyScore("unsafe-url", 0, 0)
	assert.Equal(t, totalScore, 2)
	assert.Equal(t, maxScore, 5)
	totalScore, maxScore = GetReferrerPolicyScore("", 0, 0)
	assert.Equal(t, totalScore, 0)
	assert.Equal(t, maxScore, 5)
}

func TestGetCSPScore(t *testing.T) {
	totalScore, maxScore := GetCSPScore("default-src https:; https://example.com/report-uri https://example.com/csp-violation-report-endpoint/", 0, 0)
	assert.Equal(t, totalScore, 5)
	assert.Equal(t, maxScore, 5)
	totalScore, maxScore = GetCSPScore("", 0, 0)
	assert.Equal(t, totalScore, 3)
	assert.Equal(t, maxScore, 5)
}

func TestGetPKPScore(t *testing.T) {
	totalScore, maxScore := GetPKPScore("pin - sha256 = \"cUPcTAZWKaASuYWhhneDttWpY3oBAkE3h2+soZS7sWs=\" pin - sha256 = \"M8HztCzM3elUxkcjR2S5P4hhyBNf6lHkmjAHKhpGPWE=\" max - age = 5184000 includeSubDomains report - uri = \"https://www.example.org/hpkp-report\"", 0, 0)
	assert.Equal(t, totalScore, 5)
	assert.Equal(t, maxScore, 5)
	totalScore, maxScore = GetPKPScore("", 0, 0)
	assert.Equal(t, totalScore, 3)
	assert.Equal(t, maxScore, 5)
}

func TestGetXContentTypeScore(t *testing.T) {
	totalScore, maxScore := GetXContentTypeScore("nosniff", 0, 0)
	assert.Equal(t, totalScore, 5)
	assert.Equal(t, maxScore, 5)
	totalScore, maxScore = GetXContentTypeScore("", 0, 0)
	assert.Equal(t, totalScore, 0)
	assert.Equal(t, maxScore, 5)
}

func TestGetHTTPVersionScore(t *testing.T) {
	totalScore, maxScore := GetHTTPVersionScore("HTTP/2.0", 0, 0)
	assert.Equal(t, totalScore, 5)
	assert.Equal(t, maxScore, 5)
	totalScore, maxScore = GetHTTPVersionScore("HTTP/1.1", 0, 0)
	assert.Equal(t, totalScore, 2)
	assert.Equal(t, maxScore, 5)
	totalScore, maxScore = GetHTTPVersionScore("", 0, 0)
	assert.Equal(t, totalScore, 0)
	assert.Equal(t, maxScore, 5)
}

func TestGetDMARCScore(t *testing.T) {
	score, _ := GetDMARCScore("google.com")
	assert.Equal(t, score, 5)
	score, _ = GetDMARCScore("www.google.com")
	assert.Equal(t, score, 0)
}

func TestGetMailServerConfigurationScore(t *testing.T) {
	score, maxScore, _, _ := GetMailServerConfigurationScore("google.com")
	assert.Equal(t, score, 8)
	assert.Equal(t, maxScore, 10)
	score, maxScore, _, _ = GetMailServerConfigurationScore("www.google.com")
	assert.Equal(t, score, 8)
	assert.Equal(t, maxScore, 10)
}
