package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCertificates(t *testing.T) {
	// returns certificates for valid https urls
	results, error := GetCertificate("example.com", "443", "https")
	assert.Nil(t, error)
	assert.Equal(t, results.DomainName, "example.com", "Domain Names should be equal")
	assert.NotEmpty(t, results.certChain)
	assert.NotEmpty(t, results.CommonName)
	assert.NotEmpty(t, results.Issuer)
	assert.NotEmpty(t, results.SANs)
	// returns nil for non-https endpoints
	results, error = GetCertificate("example.com", "80", "http")
	assert.Nil(t, results)
	assert.Nil(t, error)
}
