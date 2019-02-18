package utils

import (
	"testing"
)

func TestSplitHostPort(t *testing.T) {
	host1, port1, err1 := SplitHostPort("example.com:443")
	if host1 != "example.com" {
		t.Errorf("Incorrect Host")
	}
	if port1 != "443" {
		t.Errorf("Incorrect Port")
	}
	if err1 != nil {
		t.Error(err1)
	}

	host2, port2, err2 := SplitHostPort("hello.com")
	if host2 != "hello.com" {
		t.Errorf("Incorrect Host")
	}
	if port2 != "443" {
		t.Errorf("Incorrect Port")
	}
	if err2 != nil {
		t.Error(err2)
	}
}
