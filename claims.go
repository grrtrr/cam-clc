package clccam

import (
	"fmt"
	"time"

	humanize "github.com/dustin/go-humanize"
)

// Claims contains a subset of the fields contained in a CAM OAuth Bearer Token payload.
type Claims struct {
	/*
	 * 1) Fields common to both types of token
	 */
	// Unix expiration time/date
	Exp int64 `json:"exp"`

	// Unix issue time/date
	Iat int64 `json:"iat"`

	// Token Type: one of "user"  or  "service"
	Type string `json:"type"`

	/*
	 * 2) Fields specific to 'user' type tokens
	 */
	// The unique CAM username
	Subject string `json:"sub,omitempty"`

	// Name field (seems to be internal for full name)
	Name string `json:"name,omitempty"`

	// Organization name
	Organization string `json:"organization,omitempty"`

	/*
	 * 2) Fields specific to 'service' type tokens
	 */
	// Service name (e.g. "eb-1cm83")
	Service string `json:"service,omitempty"`

	// Instance name
	Instance string `json:"instance,omitempty"`

	// Machine name
	Machine string `json:"machine,omitempty"`
}

// Expired returns true if @c is already expired
func (c *Claims) Expired() bool {
	return c.Exp > 0 && time.Since(c.Expires()) > 0
}

// Expires returns the expiration time
func (c *Claims) Expires() time.Time {
	if c == nil || c.Exp == 0 {
		return time.Time{} // ensure that Expires().IsZero() returns true
	}
	return time.Unix(c.Exp, 0)
}

// Issued returns the issue time
func (c *Claims) Issued() time.Time {
	if c == nil {
		return time.Time{}
	}
	return time.Unix(c.Iat, 0)
}

func (c Claims) String() string {
	var s string

	switch c.Type {
	case "user":
		s = fmt.Sprintf("CAM user token for %q (%s at %s)", c.Subject, c.Name, c.Organization)
	case "service":
		s = fmt.Sprintf("CAM access token for service %s on %s/%s", c.Service, c.Instance, c.Machine)
	default:
		s = fmt.Sprintf("CAM %s token", c.Type)
	}

	if exp := c.Expires(); exp.IsZero() {
		s = "Permanent " + s
	} else if c.Expired() {
		s += fmt.Sprintf(", expired on %s", exp.Format(time.UnixDate))
	} else {
		s += fmt.Sprintf(", expires %s", humanize.Time(exp))
	}
	return s
}
