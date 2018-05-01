package clccam

import (
	uuid "github.com/satori/go.uuid"
)

// Event represents a single box event
type Event struct {
	// URI path of this event
	URL URI `json:"url"`

	// MIME type of the file
	ContentType string `json:"content_type"`

	// Target/destination path of the file
	DestinationPath string `json:"destination_path"`

	// Size of the file (in bytes?)
	Length int64 `json:"length"`

	// Time/date of upload
	UploadDate Timestamp `json:"upload_date"`
}

// BoxVariable specifies a single variable associated with a Box
type BoxVariable struct {
	BasicBoxVariable
	Type             string     `json:"type"`
	Options          string     `json:"options"`
	Required         bool       `json:"required"`
	Scope            string     `json:"scope"`
	Visibility       Visibility `json:"visibility"`
	AutomaticUpdates string     `json:"automatic_updates"`
}

// BasicBoxVariable is used e.g. inside a ServiceBox
type BasicBoxVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (b BasicBoxVariable) String() string { return b.Name + `="` + b.Value + `"` }

// Profile contains profile data associated with a Box.
type Profile struct {
	Cloud          string             `json:"cloud"`
	Flavor         string             `json:"flavor"`
	Image          string             `json:"image"`
	Instances      int64              `json:"instances"`
	Keypair        string             `json:"keypair"`
	Location       string             `json:"location"`
	ManagedOs      bool               `json:"managed_os"`
	PricingInfo    PricingInformation `json:"pricing_info"`
	Schema         URI                `json:"schema"`
	SecurityGroups []string           `json:"security_groups"`
	Subnet         string             `json:"subnet"`
	Volumes        []Volume           `json:"volumes"`
}

// PricingInformation aggregates the pricing information associated with a Profile
type PricingInformation struct {
	EstimatedMonthly int64  `json:"estimated_monthly"`
	Factor           int64  `json:"factor"`
	HourlyPrice      int64  `json:"hourly_price"`
	ProviderType     string `json:"provider_type"`
}

// Volume represents a single disk volme
type Volume struct {
	DeleteOnTermination bool   `json:"delete_on_termination"`
	Device              string `json:"device"`
	Size                int64  `json:"size"`
	Type                string `json:"type"`
}

// Readme represents a READ.ME file
type Readme struct {
	ContentType string `json:"content_type"`
	Length      int64  `json:"length"`
	UploadDate  string `json:"upload_date"`
	URL         URI    `json:"url"`
}

func (r Readme) String() string { return r.URL.String() }

// Service lists a service associated with a Box
type Service struct {
	Name   string     `json:"name"`
	Box    ServiceBox `json:"box"`
	Policy struct {
		Requirements []string      `json:"requirements"`
		Variables    []interface{} `json:"variables"`
	} `json:"policy"`
	AutomaticReconfiguration bool     `json:"automatic_reconfiguration"`
	AutomaticUpdates         string   `json:"automatic_updates"`
	Tags                     []string `json:"tags"`
}

// A ServiceBox is a Box inside a Service
type ServiceBox struct {
	ID        uuid.UUID          `json:"id"`
	Latest    bool               `json:"latest"`
	Variables []BasicBoxVariable `json:"variables"`
}

// Box represents a single box
type Box struct {
	// Box unique identificator.
	ID uuid.UUID `json:"id"`

	// Human readable version of @ID
	FriendlyID string `json:"friendly_id"`

	// Box name
	Name string `json:"name"`

	// Indicates at what level the box is visible.
	Visibility Visibility `json:"visibility"`

	AutomaticUpdates string      `json:"automatic_updates"`
	Categories       []string    `json:"categories"`
	Claims           []string    `json:"claims"`
	Deleted          interface{} `json:"deleted"`

	// Box description.
	Description string `json:"description"`

	// Box requirements.
	Requirements []string `json:"requirements"`

	// List of box variables, each variable object contains
	// the parameters: type, name and value (plus a few more).
	Variables []BoxVariable `json:"variables"`

	// Creation date
	Created Timestamp `json:"created"`

	// Date of the last update.
	Updated Timestamp `json:"updated"`

	// Box URI
	URI URI `json:"uri"`

	// Box schema uri.
	Schema URI `json:"schema"`

	// List of Box members.
	Members []struct {
		Role      string `json:"role"`
		Workspace string `json:"workspace"`
	} `json:"members"`

	// Organization to which the box belongs.
	Organization string `json:"organization"`

	// Box owner.
	Owner string `json:"owner"`

	DraftFrom string `json:"draft_from"`

	// Map of box events
	Events map[BoxEvent]Event

	// Box icon URI
	Icon string `json:"icon"`

	IconMetadata struct {
		Border string `json:"border"`
		Fill   string `json:"fill"`
		Image  URI    `json:"image"`
	} `json:"icon_metadata"`

	Profile Profile `json:"profile"`

	ProviderID string `json:"provider_id"`

	Readme Readme `json:"readme"`

	Services []Service `json:"services"`

	Template struct {
		ContentType string `json:"content_type"`
		Length      int64  `json:"length"`
		UploadDate  string `json:"upload_date"`
		URL         string `json:"url"`
	} `json:"template"`

	Type string `json:"type"`

	ActionButton struct {
		Icon  URI    `json:"icon"`
		Label string `json:"label"`
		Ref   URI    `json:"ref"`
	} `json:"action_button"`
}

// GetBoxes lists all boxes that are accessible in the personal workspace of the authenticated user.
func (c *Client) GetBoxes() ([]Box, error) {
	var res []Box

	return res, c.Get("/services/boxes", &res)
}
