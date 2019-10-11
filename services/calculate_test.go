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

func MockBuildResponseHeaderScore(rh ResponseHeader) (*HeaderScore, error) {
	var hScore HeaderScore
	err := rh(&hScore)
	if err != nil {
		return nil, err
	}
	return &hScore, nil
}

func TestCalculateProtocolScore(t *testing.T) {
	protocolScore := CalculateProtocolScore("http")
	assert.Equal(t, protocolScore, 0)

	protocolScore = CalculateProtocolScore("https")
	assert.Equal(t, protocolScore, 5)
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
	xssScore, err := MockBuildResponseHeaderScore(GetXSSScore("0"))
	assert.Equal(t, xssScore.value, 0)
	assert.Nil(t, err)

	xssScore, err = MockBuildResponseHeaderScore(GetXSSScore("1"))
	assert.Equal(t, xssScore.value, 5)
	assert.Equal(t, xssScore.meta, "")
	assert.Nil(t, err)

	xssScore, err = MockBuildResponseHeaderScore(GetXSSScore("1; report=https://www.example.com"))
	assert.Equal(t, xssScore.value, 5)
	assert.Equal(t, xssScore.meta, "https://www.example.com")
	assert.Nil(t, err)

	xssScore, err = MockBuildResponseHeaderScore(GetXSSScore("1; mode=block; report=https://www.example.com"))
	assert.Equal(t, xssScore.value, 5)
	assert.Equal(t, xssScore.meta, "https://www.example.com")
	assert.Nil(t, err)

	xssScore, err = MockBuildResponseHeaderScore(GetXSSScore(""))
	assert.Equal(t, xssScore.value, 0)
	assert.Equal(t, xssScore.meta, "")
	assert.Nil(t, err)
}

func TestGetXFrameScore(t *testing.T) {
	xFrameScore, err := MockBuildResponseHeaderScore(GetXFrameScore("DENY"))
	assert.Equal(t, xFrameScore.value, 5)
	assert.Nil(t, err)

	xFrameScore, err = MockBuildResponseHeaderScore(GetXFrameScore("deny"))
	assert.Equal(t, xFrameScore.value, 5)
	assert.Nil(t, err)

	xFrameScore, err = MockBuildResponseHeaderScore(GetXFrameScore("SAMEORIGIN"))
	assert.Equal(t, xFrameScore.value, 5)
	assert.Nil(t, err)

	xFrameScore, err = MockBuildResponseHeaderScore(GetXFrameScore("sameorigin"))
	assert.Equal(t, xFrameScore.value, 5)

	xFrameScore, err = MockBuildResponseHeaderScore(GetXFrameScore("ALLOW-FROM https://www.example.com"))
	assert.Equal(t, xFrameScore.value, 4)
	assert.Nil(t, err)

	xFrameScore, err = MockBuildResponseHeaderScore(GetXFrameScore("allow-from https://www.example.com"))
	assert.Equal(t, xFrameScore.value, 4)
	assert.Nil(t, err)

	xFrameScore, err = MockBuildResponseHeaderScore(GetXFrameScore(""))
	assert.Equal(t, xFrameScore.value, 1)
	assert.Nil(t, err)

}

func TestGetHSTSScore(t *testing.T) {
	hstsScore, err := MockBuildResponseHeaderScore(GetHSTSScore("max-age=65536"))
	assert.Equal(t, hstsScore.value, 4)
	assert.Nil(t, err)

	hstsScore, err = MockBuildResponseHeaderScore(GetHSTSScore("max-age=31536000; includeSubDomains"))
	assert.Equal(t, hstsScore.value, 5)
	assert.Nil(t, err)

	hstsScore, err = MockBuildResponseHeaderScore(GetHSTSScore("max-age=31536000; includeSubDomains; preload"))
	assert.Equal(t, hstsScore.value, 5)
	assert.Nil(t, err)

	hstsScore, err = MockBuildResponseHeaderScore(GetHSTSScore("max-age=31536;  preload"))
	assert.Equal(t, hstsScore.value, 5)
	assert.Nil(t, err)

	hstsScore, err = MockBuildResponseHeaderScore(GetHSTSScore(""))
	assert.Equal(t, hstsScore.value, 2)
	assert.Nil(t, err)

}

func TestGetReferrerPolicyScore(t *testing.T) {
	xRPScore, err := MockBuildResponseHeaderScore(GetReferrerPolicyScore("no-referrer"))
	assert.Equal(t, xRPScore.value, 5)
	assert.Nil(t, err)

	xRPScore, err = MockBuildResponseHeaderScore(GetReferrerPolicyScore("no-referrer-when-downgrade"))
	assert.Equal(t, xRPScore.value, 4)
	assert.Nil(t, err)

	xRPScore, err = MockBuildResponseHeaderScore(GetReferrerPolicyScore("origin"))
	assert.Equal(t, xRPScore.value, 4)
	assert.Nil(t, err)

	xRPScore, err = MockBuildResponseHeaderScore(GetReferrerPolicyScore("origin-when-cross-origin"))
	assert.Equal(t, xRPScore.value, 4)
	assert.Nil(t, err)

	xRPScore, err = MockBuildResponseHeaderScore(GetReferrerPolicyScore("same-origin"))
	assert.Equal(t, xRPScore.value, 4)
	assert.Nil(t, err)

	xRPScore, err = MockBuildResponseHeaderScore(GetReferrerPolicyScore("strict-origin"))
	assert.Equal(t, xRPScore.value, 4)
	assert.Nil(t, err)

	xRPScore, err = MockBuildResponseHeaderScore(GetReferrerPolicyScore("strict-origin-when-cross-origin"))
	assert.Equal(t, xRPScore.value, 4)
	assert.Nil(t, err)

	xRPScore, err = MockBuildResponseHeaderScore(GetReferrerPolicyScore("unsafe-url"))
	assert.Equal(t, xRPScore.value, 2)
	assert.Nil(t, err)

	xRPScore, err = MockBuildResponseHeaderScore(GetReferrerPolicyScore(""))
	assert.Equal(t, xRPScore.value, 0)
	assert.Nil(t, err)

}

func TestGetCSPScore(t *testing.T) {
	cspScore, err := MockBuildResponseHeaderScore(GetCSPScore("default-src https:; https://example.com/report-uri https://example.com/csp-violation-report-endpoint/"))
	assert.Equal(t, cspScore.value, 5)
	assert.Nil(t, err)

	cspScore, err = MockBuildResponseHeaderScore(GetCSPScore(""))
	assert.Equal(t, cspScore.value, 3)
	assert.Nil(t, err)

}

func TestGetPKPScore(t *testing.T) {
	pkpScore, err := MockBuildResponseHeaderScore(GetPKPScore("pin - sha256 = \"cUPcTAZWKaASuYWhhneDttWpY3oBAkE3h2+soZS7sWs=\" pin - sha256 = \"M8HztCzM3elUxkcjR2S5P4hhyBNf6lHkmjAHKhpGPWE=\" max - age = 5184000 includeSubDomains report - uri = \"https://www.example.org/hpkp-report\""))
	assert.Equal(t, pkpScore.value, 5)
	assert.Nil(t, err)

	pkpScore, err = MockBuildResponseHeaderScore(GetPKPScore(""))
	assert.Equal(t, pkpScore.value, 3)
	assert.Nil(t, err)
}

func TestGetXContentTypeScore(t *testing.T) {
	xContentTypeScore, err := MockBuildResponseHeaderScore(GetXContentTypeScore("nosniff"))
	assert.Equal(t, xContentTypeScore.value, 5)
	assert.Nil(t, err)

	xContentTypeScore, err = MockBuildResponseHeaderScore(GetXContentTypeScore(""))
	assert.Equal(t, xContentTypeScore.value, 0)
	assert.Nil(t, err)
}

func TestGetHTTPVersionScore(t *testing.T) {
	httpVersionScore, err := MockBuildResponseHeaderScore(GetHTTPVersionScore("HTTP/2.0"))
	assert.Equal(t, httpVersionScore.value, 5)
	assert.Nil(t, err)

	httpVersionScore, err = MockBuildResponseHeaderScore(GetHTTPVersionScore("HTTP/1.1"))
	assert.Equal(t, httpVersionScore.value, 2)
	assert.Nil(t, err)

	httpVersionScore, err = MockBuildResponseHeaderScore(GetHTTPVersionScore(""))
	assert.Equal(t, httpVersionScore.value, 0)
	assert.Nil(t, err)
}

func TestGetDMARCScore(t *testing.T) {
	dmarcScore, _ := GetDMARCScore("google.com")
	assert.Equal(t, dmarcScore, 5)

	dmarcScore, _ = GetDMARCScore("www.google.com")
	assert.Equal(t, dmarcScore, 0)
}

func TestGetMailServerConfigurationScore(t *testing.T) {
	mailServerScore, _, _ := GetMailServerConfigurationScore(MailServerConfigParams{host: "google.com"})
	assert.Equal(t, mailServerScore, 8)

	mailServerScore, _, _ = GetMailServerConfigurationScore(MailServerConfigParams{host: "www.google.com"})
	assert.Equal(t, mailServerScore, 8)
}
