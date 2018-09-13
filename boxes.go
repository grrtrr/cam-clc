package clccam

import (
	"fmt"

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

// BasicBoxVariable is used e.g. inside a ServiceBox
type BasicBoxVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (b BasicBoxVariable) String() string {
	return fmt.Sprintf(`%s="%s"`, b.Name, b.Value)
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

// Profile contains profile data associated with a Box.
type Profile struct {
	Cloud          string             `json:"cloud"`           // e.g. "vpc-26ebd840"
	ElasticIP      bool               `json:"elastic_ip"`      // e.g. false
	Flavor         string             `json:"flavor"`          // e.g. "c1.medium"
	Image          string             `json:"image"`           // e.g. "Linux Compute"
	Instances      int64              `json:"instances"`       // e.g. 1
	Keypair        string             `json:"keypair"`         // e.g. "None"
	Location       string             `json:"location"`        // e.g. "us-west-2"
	ManagedOs      bool               `json:"managed_os"`      // e.g. false
	PlacementGroup string             `json:"placement_group"` // e.g. ""
	PricingInfo    PricingInformation `json:"pricing_info"`
	Role           string             `json:"role"`            // e.g. "None"
	Schema         URI                `json:"schema"`          // e.g. "http://elasticbox.net/schemas/aws/ec2/profile"
	SecurityGroups []string           `json:"security_groups"` // e.g. [ "Automatic" ]
	Subnet         string             `json:"subnet"`          // e.g. "subnet-32df1d7a"
	Volumes        []Volume           `json:"volumes"`         // e.g. []
}

// PricingInformation aggregates the pricing information associated with a Profile
type PricingInformation struct {
	EstimatedMonthly int64  `json:"estimated_monthly"` // e.g. 9360000
	Factor           int64  `json:"factor"`            // e.g. 100000
	HourlyPrice      int64  `json:"hourly_price"`      // e.g. 13000
	ProviderType     string `json:"provider_type"`     // e.g. "Amazon Web Services"
}

// Volume represents a single disk volume.
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
	FriendlyID string `json:"friendly_id"` // e.g. "jenkins"

	// Box name
	Name string `json:"name"` // e.g. "Jenkins"

	// Indicates at what level the box is visible.
	Visibility Visibility `json:"visibility"`

	// Presumably "on" or "off":
	AutomaticUpdates string      `json:"automatic_updates"` // e.g. "off"
	Categories       []string    `json:"categories"`        // e.g.  [ "Continuous Integration" ]
	Claims           []string    `json:"claims"`
	Deleted          interface{} `json:"deleted"`

	// Box description.
	Description string `json:"description"` // e.g. "With ElasticBox CI plugin"

	// Box requirements.
	Requirements []string `json:"requirements"` // e.g. [ "linux" ]

	// List of box variables, each variable object contains
	// the parameters: type, name and value (plus a few more).
	Variables []BoxVariable `json:"variables"`

	// Creation date
	Created Timestamp `json:"created"`

	// Date of the last update.
	Updated Timestamp `json:"updated"`

	Lifespan struct {
		Operation string `json:"operation"` // e.g. "none"
	} `json:"lifespan"`

	// Box URI
	URI URI `json:"uri"` // e.g. "/services/boxes/e0715702-cf5c-4c88-bfa1-2e5e3808e597"

	// Box schema uri.
	Schema URI `json:"schema"` // e.g. "http://elasticbox.net/schemas/boxes/script"

	// List of Box members.
	Members []struct {
		Role      string `json:"role"`      // e.g. "collaborator"
		Workspace string `json:"workspace"` // e.g. "cf"
	} `json:"members"`

	// Organization to which the box belongs.
	Organization string `json:"organization"` // e.g. "elasticbox"

	// Box owner.
	Owner string `json:"owner"` // e.g. "elasticbox"

	// References the ID of another box
	DraftFrom uuid.UUID `json:"draft_from"`

	// Map of box events
	Events map[BoxEvent]Event

	// Profile contains cloud-specific details.
	Profile Profile `json:"profile"`

	ProviderID uuid.UUID `json:"provider_id"`

	Services []Service `json:"services"`

	Template struct {
		ContentType string `json:"content_type"` // e.g. "text/x-shellscript"
		Length      int64  `json:"length"`       // e.g. 4395
		UploadDate  string `json:"upload_date"`  // e.g. "2017-05-23 17:05:43.923946"
		URL         URI    `json:"url"`          // e.g. "/services/blobs/download/59246be776d194287289646c/template.json"
	} `json:"template"`

	// Type seems to be often empty
	Type string `json:"type"` // e.g. "CloudFormation Service"

	/*
	 * Html Section
	 */
	Readme Readme `json:"readme"`

	// Box icon URI
	Icon string `json:"icon"`

	IconMetadata struct {
		Border string `json:"border"`
		Fill   string `json:"fill"`
		Image  URI    `json:"image"`
	} `json:"icon_metadata"`

	// More html ...
	ActionButton struct {
		Icon  URI    `json:"icon"`
		Label string `json:"label"`
		Ref   URI    `json:"ref"`
	} `json:"action_button"`
}

// GetBox returns the details of box @boxId.
func (c *Client) GetBox(boxId string) (res Box, err error) {
	return res, c.Get("/services/boxes/"+boxId, &res)
}

// GetBoxes lists all boxes that are accessible in the personal workspace of the authenticated user.
func (c *Client) GetBoxes() (res []Box, err error) {
	return res, c.Get("/services/boxes", &res)
}
