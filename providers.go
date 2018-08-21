package clccam

import uuid "github.com/satori/go.uuid"

type Provider struct {
	// ID of the provider account in Cloud Application Manager
	ID uuid.UUID `json:"id"`

	// Name used to identify the provider account in Cloud Application Manager.
	Name string `json:"name"`

	// Provider account owner in Cloud Application Manager
	Owner string `json:"owner"`

	// Users with whom the provider is shared
	Members []interface{} `json:"members"`

	// Date the provider was added
	Created Timestamp `json:"created"`

	// Date the provider was updated
	Updated Timestamp `json:"updated"`

	// Description of the provider
	Description string `json:"description"`

	Credentials struct{} `json:"credentials"`

	// Icon used for the provider account.
	Icon URI `json:"icon"`

	Services []struct {
		Locations []struct {
			Clusters []interface{} `json:"clusters"`
		} `json:"locations"`
		Name string `json:"name"`
	} `json:"services"`

	// E.g. "ready"
	State string `json:"state"`

	// Identifies the provider as one of the following:
	// - Amazon Web Services,
	// - Rackspace, Openstack,
	// - VMware vSphere,
	// - Google Compute,
	// - Microsoft Azure,
	// - Cloudstack,
	// - SoftLayer.
	Type string `json:"type"`

	// Unique resource identifier path to the provider account
	URI URI `json:"uri"`

	// Schema URi
	Schema URI `json:"schema"`
}

// GetProviders lists all providers.
func (c *Client) GetProviders() ([]Provider, error) {
	var res []Provider

	return res, c.Get("/services/providers", &res)
}
