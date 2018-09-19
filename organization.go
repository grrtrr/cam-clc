package clccam

import (
	"github.com/coreos/go-semver/semver"
	uuid "github.com/satori/go.uuid"
)

// Organization contains the data associated with a CAM organization.
// See https://www.ctl.io/api-docs/cam/#cam-platform-organizations-api
type Organization struct {
	// Organization name
	Name string `json:"name"`

	// Billing account alias
	ClcAlias string `json:"clc_alias"`

	// Status can e.g. be "Active"
	AccountStatus        string      `json:"account_status"`
	AccountType          string      `json:"account_type"`
	BillingAccountNumber interface{} `json:"billing_account_number"`
	RemedyAccountID      interface{} `json:"remedy_account_id"`

	DefaultCostcenter uuid.UUID `json:"default_costcenter"`
	DisplayName       *string   `json:"display_name"`
	Domains           []string  `json:"domains"`
	FederatedTo       []string  `json:"federated_to"`

	// Release version, e.g. "4.0"
	Release semver.Version `json:"release"`

	// Organization schema URI
	Schema URI `json:"schema"`

	Theme struct {
		Accent interface{} `json:"accent"`
		CSS    URI         `json:"css"`
		Logo   URI         `json:"logo"`
	} `json:"theme"`

	// Organization icon URI
	Icon URI `json:"icon"`

	VantiveID interface{} `json:"vantive_id"`
}

// GetOrganization gets the organization schema of @orgName
func (c *Client) GetOrganization(orgName string) (*Organization, error) {
	var res = new(Organization)

	return res, c.Get("/services/organizations/"+orgName, &res)
}
