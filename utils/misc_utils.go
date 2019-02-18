package utils

import (
	"net"
	"strings"
)

var defaultPort = "443"

// TODO: We should replace this with native net/url.SplitHostPort instead of using our own method

// SplitHostPort returns a Host and Port seperately in a host-port URL
func SplitHostPort(hostport string) (string, string, error) {
	if !strings.Contains(hostport, ":") {
		return hostport, defaultPort, nil
	}

	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		return "", "", err
	}
	if port == "" || len(port) == 0 {
		port = defaultPort
	}

	return host, port, nil
}
