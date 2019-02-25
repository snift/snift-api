package models

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"time"
)

// TimeoutSeconds references the total Time Out duration for the Handshake
var TimeoutSeconds = 3

//Cert holds the certificate details
type Cert struct {
	DomainName string   `json:"domain_name"`
	IP         string   `json:"ip_address"`
	Issuer     string   `json:"issuer"`
	CommonName string   `json:"common_name"`
	SANs       []string `json:"sans"`
	NotBefore  string   `json:"not_before"`
	NotAfter   string   `json:"not_after"`
	certChain  []*x509.Certificate
}

var serverCert = func(host string, port string) ([]*x509.Certificate, string, error) {
	d := &net.Dialer{
		Timeout: time.Duration(TimeoutSeconds) * time.Second,
	}
	fmt.Println(host + " " + port)
	conn, err := tls.DialWithDialer(d, "tcp", host+":"+port, &tls.Config{
		InsecureSkipVerify: false,
	})
	if err != nil {
		return []*x509.Certificate{&x509.Certificate{}}, "", err
	}
	defer conn.Close()

	addr := conn.RemoteAddr()
	ip, _, _ := net.SplitHostPort(addr.String())
	cert := conn.ConnectionState().PeerCertificates

	return cert, ip, nil
}

// GetCertificate returns the Certificate associated with a host-port
func GetCertificate(host string, port string, protocol string) (*Cert, error) {
	// dont get certificates for non-https protocols, and when port number is 80
	// trying to fetch certs with port:80 causes tls overload
	if protocol != "https" || (protocol == "https" && port == "80") {
		return nil, nil
	}
	certChain, ip, err := serverCert(host, port)
	if err != nil {
		return &Cert{DomainName: host}, err
	}
	cert := certChain[0]

	var loc = time.UTC // Setting UTC as Standard Time

	return &Cert{
		DomainName: host,
		IP:         ip,
		Issuer:     cert.Issuer.CommonName,
		CommonName: cert.Subject.CommonName,
		SANs:       cert.DNSNames, // Subject Alternative Name
		NotBefore:  cert.NotBefore.In(loc).String(),
		NotAfter:   cert.NotAfter.In(loc).String(),
		certChain:  certChain,
	}, nil
}
