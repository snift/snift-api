package services

// XSSHeader has the XSS Header Name
const XSSHeader = "X-Xss-Protection"

// XFrameHeader has the XFrame Header Name
const XFrameHeader = "X-Frame-Options"

// HSTSHeader has the HSTS Header Name
const HSTSHeader = "Strict-Transport-Security"

// CSPHeader has the CSP Header Name
const CSPHeader = "Content-Security-Policy"

// PKPHeader has the PKP Header Name
const PKPHeader = "Public-Key-Pins"

// RPHeader has the RP Header Name
const RPHeader = "Referrer-Policy"

// XContentTypeHeader has the X-Content-Type Header Name
const XContentTypeHeader = "X-Content-Type-Options"

// Server has the Server Header
const Server = "Server"

// TXTQuery is used to extract all the TXT Records of a Domain
const TXTQuery = "dig @8.8.8.8 +ignore +short +bufsize=1024 domain.com txt"

// DMARCQuery is used to extract all the DMARC Records of a Domain
const DMARCQuery = "dig +short TXT _dmarc.domain.com"

// OpenBugBountyURL is used to query for previous security incidents
const OpenBugBountyURL = "https://www.openbugbounty.org/api/1/search/?domain="

// MaxIncidentResponseTime is the Maximum Incident Response Time taken as 30 days -> 30 * 24 = 720 hours
const MaxIncidentResponseTime = 720

// XSSValues is used to store the X-Xss-Protection Header values
var XSSValues = [...]string{"0", "1"}

// XFrameValues is used to store the X-Frame-Options Header values
var XFrameValues = [...]string{"deny", "sameorigin", "allow-from"}

// HSTSValues used to store the X-Frame-Options Header values
var HSTSValues = [...]string{"max-age", "includeSubDomains", "preload"}

// ReferrerPolicyValues used to store the Referrer-Policy Header values
var ReferrerPolicyValues = [...]string{"no-referrer", "no-referrer-when-downgrade", "origin", "origin-when-cross-origin", "same-origin", "strict-origin", "strict-origin-when-cross-origin", "unsafe-url"}

// XContentTypeHeaderValue is used to store the value for X-Content-Type Options Header
const XContentTypeHeaderValue = "nosniff"

// HTTPVersion is used to store the HTTP Versions
var HTTPVersion = [...]string{"HTTP/2.0", "HTTP/1.1"}

// Stores the Scores for various Parameters
const (
	HTTPScore  = 0
	HTTPSScore = 5
)
