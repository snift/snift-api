package utils

import (
	"snift-api/models"
)

// CreateBadge creates a Badge struct with given properties
func createBadge(badgeName string, badgeMessage string, badgeCategory string) *models.Badge {
	return &models.Badge{
		Name:     badgeName,
		Message:  badgeMessage,
		Category: badgeCategory,
	}
}

// GetHTTPSBadge returns the HTTP Secure Badge
func GetHTTPSBadge() *models.Badge {
	return createBadge(HTTPSBadge, HTTPSBadgeMessage, "NETWORK_PROTECTION")
}

// GetHSTSBadge returns the HTTP Strict Transport Security Badge
func GetHSTSBadge() *models.Badge {
	return createBadge(HSTSBadge, HSTSBadgeMessage, "NETWORK_PROTECTION")
}

// GetHTTPVersionBadge returns the HyperText Transfer Protocol(HTTP) Version Badge
func GetHTTPVersionBadge() *models.Badge {
	return createBadge(HTTPVersionBadge, HTTPVersionBadgeMessage, "NETWORK_PROTECTION")
}

// GetTLSVersionBadge returns the Transport Layer Security Version Badge
func GetTLSVersionBadge() *models.Badge {
	return createBadge(TLSVersionBadge, TLSVersionBadgeMessage, "NETWORK_PROTECTION")
}

// GetXSSBadge returns the XSS-Protection Badge
func GetXSSBadge() *models.Badge {
	return createBadge(XSSBadge, XSSBadgeMessage, "CONTENT_SECURITY")
}

// GetXFrameBadge returns the X-Frame Options Badge
func GetXFrameBadge() *models.Badge {
	return createBadge(XFrameBadge, XFrameBadgeMessage, "CONTENT_SECURITY")
}

// GetCSPBadge returns the Content Security Policy Badge
func GetCSPBadge() *models.Badge {
	return createBadge(CSPBadge, CSPBadgeMessage, "CONTENT_SECURITY")
}

// GetHPKPBadge returns the HTTP Public Key Pinning Badge
func GetHPKPBadge() *models.Badge {
	return createBadge(HPKPBadge, HPKPBadgeMessage, "EAVESDROPPING_SPOOFING_PROTECTION")
}

// GetSPFBadge returns the Sender Policy Framework Badge
func GetSPFBadge() *models.Badge {
	return createBadge(SPFBadge, SPFBadgeMessage, "EAVESDROPPING_SPOOFING_PROTECTION")
}

// GetRPBadge returns the Referrer Policy Badge
func GetRPBadge() *models.Badge {
	return createBadge(RPBadge, RPBadgeMessage, "USER_PRIVACY")
}

// GetXContentTypeBadge returns the X-Content-Type Options Badge
func GetXContentTypeBadge() *models.Badge {
	return createBadge(XContentTypeBadge, XContentTypeBadgeMessage, "USER_PRIVACY")
}
