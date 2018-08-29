package clccam

import (
	"fmt"
	"os"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/grrtrr/clccam/logger"
	"github.com/olekukonko/tablewriter"
	uuid "github.com/satori/go.uuid"
)

// Claims contains a subset of the fields contained in a CAM OAuth Bearer Token payload.
type Claims struct {
	/*
	 * 1) Fields common to both types of token:
	 */
	// Token Type: one of "user"  or  "service"
	Type string `json:"type"`

	// Unix expiration time/date
	Exp int64 `json:"exp"`

	// Unix time/date of issuing the token.
	Iat int64 `json:"iat"`

	// Unique token ID
	Jti uuid.UUID `json:"jti"`

	/*
	 * 2) Fields specific to 'user' type tokens:
	 */
	// The unique CAM username
	Subject string `json:"sub,omitempty"`

	// Name field (seems to be internal for full name)
	Name string `json:"name,omitempty"`

	// Organization name (e.g. "centurylink")
	Organization string `json:"organization,omitempty"`

	/*
	 * 2) Fields specific to 'service' type tokens:
	 */
	// Instance Id (e.g. "i-z48wub")
	InstanceId string `json:"instance,omitempty"`

	// Service Id (e.g. "eb-1cm83")
	ServiceId string `json:"service,omitempty"`

	// Machine Id (e.g. "cms1-eb-e775t-1")
	MachineId string `json:"machine,omitempty"`
}

// Expired returns true if @c is already expired.
func (c *Claims) Expired() bool {
	return !c.IsPermanent() && time.Since(c.Expires()) > 0
}

// IsPermanent returns true if @c has a zero expiry date.
func (c *Claims) IsPermanent() bool {
	return c.Exp == 0
}

// Expires returns the expiration time.
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
		s = fmt.Sprintf("CAM token for service %s on %s/%s", c.ServiceId, c.InstanceId, c.MachineId)
	default:
		s = fmt.Sprintf("CAM %s token", c.Type)
	}

	if exp := c.Expires(); c.IsPermanent() {
		s = "Permanent " + s
	} else {
		if c.Expired() {
			s += fmt.Sprintf(", expired %s", humanize.Time(exp))
		} else {
			s += fmt.Sprintf(", expires %s", humanize.Time(exp))
		}
	}
	return s
}

// DumpToStdout prints a representation of @cl to stdout.
func (c Claims) DumpToStdout() {
	const timeFmt = `Mon Jan _2 15:04:05 MST 2006`
	var (
		table = tablewriter.NewWriter(os.Stdout)
		exp   = time.Unix(c.Exp, 0).Format(timeFmt)
	)

	if c.IsPermanent() {
		exp = "never (permanent token)"
	}

	table.SetAutoFormatHeaders(false)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"Field", "Token Value"})

	table.AppendBulk([][]string{
		[]string{"exp", exp},
		[]string{"iat", time.Unix(c.Iat, 0).Format(timeFmt)},
		[]string{"jti", c.Jti.String()},
	})

	switch c.Type {
	case "user":
		table.AppendBulk([][]string{
			[]string{"sub", c.Subject},
			[]string{"name", c.Name},
			[]string{"organization", c.Organization},
		})
	case "service":
		table.AppendBulk([][]string{
			[]string{"instance", c.InstanceId},
			[]string{"machine", c.MachineId},
			[]string{"service", c.ServiceId},
		})
	default:
		logger.Fatalf("unexpected token type %q", c.Type)
	}
	fmt.Printf("%s:\n", c)
	table.Render()
}
