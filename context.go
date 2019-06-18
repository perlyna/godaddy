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
	XShopperID    string
}

// NewGoDaddy return godaddy context
func NewGoDaddy(key, secret string) *Context {
	authorization := fmt.Sprintf("sso-key %s:%s", key, secret)
	return &Context{key: key, secret: secret, Host: Host, authorization: authorization}
}
