package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateProtocolScore(t *testing.T) {
	score, message := CalculateProtocolScore("http")
	assert.Equal(t, score, 0)
	assert.Equal(t, message, "Website is unencrypted and hence subjective to Man-in-the-Middle attacks(MITM) and Eavesdropping Attacks.")
	score, message = CalculateProtocolScore("https")
	assert.Equal(t, score, 5)
	assert.Equal(t, message, "From the protocol level, Website is secure.")
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
	score, XSSReportURL := GetXSSScore("0")
	assert.Equal(t, score, 0)
	assert.Equal(t, XSSReportURL, "")
	score, XSSReportURL = GetXSSScore("1")
	assert.Equal(t, score, 5)
	assert.Equal(t, XSSReportURL, "")
	score, XSSReportURL = GetXSSScore("1; report=https://www.example.com")
	assert.Equal(t, score, 5)
	assert.Equal(t, XSSReportURL, "https://www.example.com")
	score, XSSReportURL = GetXSSScore("1; mode=block; report=https://www.example.com")
	assert.Equal(t, score, 5)
	assert.Equal(t, XSSReportURL, "https://www.example.com")
}

func TestGetXFrameScore(t *testing.T) {
	assert.Equal(t, GetXFrameScore("DENY"), 5)
	assert.Equal(t, GetXFrameScore("deny"), 5)
	assert.Equal(t, GetXFrameScore("SAMEORIGIN"), 5)
	assert.Equal(t, GetXFrameScore("sameorigin"), 5)
	assert.Equal(t, GetXFrameScore("ALLOW-FROM https://www.example.com"), 4)
	assert.Equal(t, GetXFrameScore("allow-from https://www.example.com"), 4)
}

func TestGetHSTSScore(t *testing.T) {
	assert.Equal(t, GetHSTSScore("max-age=65536"), 4)
	assert.Equal(t, GetHSTSScore("max-age=31536000; includeSubDomains"), 5)
	assert.Equal(t, GetHSTSScore("max-age=31536000; includeSubDomains; preload"), 5)
	assert.Equal(t, GetHSTSScore("max-age=31536;  preload"), 5)
}

func TestGetReferrerPolicyScore(t *testing.T) {
	assert.Equal(t, GetReferrerPolicyScore("no-referrer"), 5)
	assert.Equal(t, GetReferrerPolicyScore("no-referrer-when-downgrade"), 4)
	assert.Equal(t, GetReferrerPolicyScore("origin"), 4)
	assert.Equal(t, GetReferrerPolicyScore("origin-when-cross-origin"), 4)
	assert.Equal(t, GetReferrerPolicyScore("same-origin"), 4)
	assert.Equal(t, GetReferrerPolicyScore("strict-origin"), 4)
	assert.Equal(t, GetReferrerPolicyScore("strict-origin-when-cross-origin"), 4)
	assert.Equal(t, GetReferrerPolicyScore("unsafe-url"), 2)
}
