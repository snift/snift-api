package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidURL(t *testing.T) {
	assert.NoError(t, IsValidURL("https://www.example.com"))
	assert.NoError(t, IsValidURL("http://www.example.com"))
	assert.NoError(t, IsValidURL("http://example.com"))
	assert.NoError(t, IsValidURL("https://example.com"))
	assert.Error(t, IsValidURL("example.com"))
	assert.Error(t, IsValidURL("example-domain"))
	assert.Error(t, IsValidURL("example"))
}
