package models

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetCertificates(t *testing.T) {
	results := GetCertificate("example.com")
	assert.Equal(t, results.DomainName, "example.com", "Domain Names should be equal")
	assert.NotEmpty(t, results.certChain)
	assert.NotEmpty(t, results.CommonName)
	assert.NotEmpty(t, results.Issuer)
	assert.NotEmpty(t, results.SANs)
	assert.Empty(t, results.Error)
}
