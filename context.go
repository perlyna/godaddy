package godaddy

import (
	"fmt"
)

const (
	// DefaultTTL default ttl
	DefaultTTL = 600
	// Host godaddy api host
	Host = "api.godaddy.com"
	// OTEHost godaddy test api host
	OTEHost = "api.ote-godaddy.com"
	// request size limit
	limit = 100
)

// Context godaddy context
type Context struct {
	key           string
	secret        string
	Host          string
	authorization string
}

// NewGoDaddy return godaddy context
func NewGoDaddy(key, secret string) *Context {
	authorization := fmt.Sprintf("sso-key %s:%s", key, secret)
	return &Context{key: key, secret: secret, Host: Host, authorization: authorization}
}

// RecordType is an enumeration of possible DNS record types
type RecordType int

const (
	// A is an address record type
	A RecordType = iota
	// AAAA is an IPv6 address record type
	AAAA
	// CNAME is a Canonical record name (alias) type
	CNAME
	// MX is a mail exchange record type
	MX
	// NS is a name server record type
	NS
	// SOA is a start of authority record type
	SOA
	// SRV is a service locator type
	SRV
	// TXT is a text record type
	TXT
)
