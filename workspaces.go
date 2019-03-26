package clccam

import (
	uuid "github.com/satori/go.uuid"
)

// Thanks to https://mholt.github.io/json-to-go/

type WorkSpace struct {
	// Workspace unique identifier.
	ID string `json:"id"`

	// Workspace URI
	URI URI `json:"uri"`

	// Indicates true if a personal workspace has a provider.
	AddProvider bool      `json:"add_provider"`
	Costcenter  uuid.UUID `json:"costcenter"`

	// Time/date of creation
	Created Timestamp `json:"created"`

	// Time/date of the last update
	Updated   Timestamp `json:"updated"`
	LastLogin Timestamp `json:"last_login"`

	Deleted *Timestamp `json:"deleted"`

	// Shows true when there are deployed instances in the personal workspace.
	DeployInstance bool `json:"deploy_instance"`

	Favorites []interface{} `json:"favorites"`

	// List of fully qualified names of LDAP groups to which a user’s personal workspace belongs.
	// You can’t update this field. Present in Personal Workspaces
	GroupDNS []interface{} `json:"group_dns"`

	// Workspace icon.
	Icon interface{} `json:"icon"`

	// Workspace name
	Name             string    `json:"name"`
	LastName         string    `json:"last_name"`
	Organization     string    `json:"organization"`
	SamlID           string    `json:"saml_id"`
	Email            string    `json:"email"`
	EmailValidatedAt Timestamp `json:"email_validated_at"`

	// Schema URI.
	// Either "http://elasticbox.net/schemas/workspaces/personal"
	// or     "http://elasticbox.net/schemas/workspaces/team"
	Schema URI `json:"schema"`

	SupportUserCreated bool   `json:"support_user_created"`
	TakeTour           bool   `json:"take_tour"`
	Type               string `json:"type"`
}

// GetWorkSpace returns the workspace of @userId, if any.
func (c *Client) GetWorkSpace(userId string) (*WorkSpace, error) {
	var res = new(WorkSpace)

	return res, c.Get("/services/workspaces/"+userId, &res)
}

// GetWorkSpaces returns the list of all accessible workspaces.
// There are two types of workspaces:
// a) personal workspaces for a single user, and
// b) eam workspaces that can have many members and organizations.
func (c *Client) GetWorkSpaces() ([]WorkSpace, error) {
	var res []WorkSpace

	return res, c.Get("/services/workspaces", &res)
}
