package utils

import (
	"snift-api/models"
)

// CreateBadge creates a Badge struct with given properties
func createBadge(badgeName string, badgeMessage string) *models.Badge {
	return &models.Badge{
		Name:    badgeName,
		Message: badgeMessage,
	}
}

// GetHTTPSBadge returns the HTTP Secure Badge
func GetHTTPSBadge() *models.Badge {
	return createBadge(HTTPSBadge, HTTPSBadgeMessage)
}

// GetXSSBadge returns the XSS-Protection Badge
func GetXSSBadge() *models.Badge {
	return createBadge(XSSBadge, XSSBadgeMessage)
}

// GetXFrameBadge returns the X-Frame Options Badge
func GetXFrameBadge() *models.Badge {
	return createBadge(XFrameBadge, XFrameBadgeMessage)
}

// GetHSTSBadge returns the HTTP Strict Transport Security Badge
func GetHSTSBadge() *models.Badge {
	return createBadge(HSTSBadge, HTTPSBadgeMessage)
}

// GetCSPBadge returns the Content Security Policy Badge
func GetCSPBadge() *models.Badge {
	return createBadge(CSPBadge, CSPBadgeMessage)
}

// GetHPKPBadge returns the HTTP Public Key Pinning Badge
func GetHPKPBadge() *models.Badge {
	return createBadge(HPKPBadge, HPKPBadgeMessage)
}

// GetRPBadge returns the Referrer Policy Badge
func GetRPBadge() *models.Badge {
	return createBadge(RPBadge, RPBadgeMessage)
}

// GetXContentTypeBadge returns the X-Content-Type Options Badge
func GetXContentTypeBadge() *models.Badge {
	return createBadge(XContentTypeBadge, XContentTypeBadgeMessage)
}

// GetHTTPVersionBadge returns the HyperText Transfer Protocol(HTTP) Version Badge
func GetHTTPVersionBadge() *models.Badge {
	return createBadge(HTTPVersionBadge, HTTPVersionBadgeMessage)
}

// GetTLSVersionBadge returns the Transport Layer Security Version Badge
func GetTLSVersionBadge() *models.Badge {
	return createBadge(TLSVersionBadge, TLSVersionBadgeMessage)
}

// GetSPFBadge returns the Sender Policy Framework Badge
func GetSPFBadge() *models.Badge {
	return createBadge(SPFBadge, SPFBadgeMessage)
}
