package clccam

import (
	"fmt"

	"github.com/coreos/go-semver/semver"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// A box having this schema identifies itself as a Script Box.
const ScriptBoxSchema = "http://elasticbox.net/schemas/boxes/script"

// The expected file name for README files.
const ReadmeName = "readme.MD"

// Event represents a single box event
type Event struct {
	BlobResponse

	// Target/destination path of the file
	DestinationPath string `json:"destination_path,omitempty"` // e.g. "scripts"
}

// BasicVariable is used e.g. inside a ServiceBox
type BasicVariable struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (b BasicVariable) String() string {
	return fmt.Sprintf(`%s="%s"`, b.Name, b.Value)
}

// BoxVariable specifies a single variable associated with a Box
type BoxVariable struct {
	BasicVariable
	Options          string     `json:"options,omitempty"`
	Required         bool       `json:"required"`
	Scope            string     `json:"scope,omitempty"`
	Visibility       Visibility `json:"visibility"`
	AutomaticUpdates string     `json:"automatic_updates,omitempty"`
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
	ID        uuid.UUID       `json:"id"`
	Latest    bool            `json:"latest"`
	Variables []BasicVariable `json:"variables"`
}

// Box represents a single box
type Box struct {
	// Box unique identificator.
	ID uuid.UUID `json:"id"`

	// Human readable version of @ID
	FriendlyID string `json:"friendly_id,omitempty"` // e.g. "jenkins"

	// Box name
	Name string `json:"name"` // e.g. "Jenkins"

	// Indicates at what level the box is visible.
	Visibility Visibility `json:"visibility"`

	// Automatic updates: seems to be one of { "off", "major", "minor", "patch" }
	AutomaticUpdates string   `json:"automatic_updates,omitempty"` // e.g. "off"
	Categories       []string `json:"categories,omitempty"`        // e.g.  [ "Continuous Integration" ]
	Claims           []string `json:"claims,omitempty"`            // e.g. [ "safehaven-cms" ]

	// Box description.
	Description string `json:"description"` // e.g. "With ElasticBox CI plugin"

	// Box requirements.
	Requirements []string `json:"requirements"` // e.g. [ "safehaven-cms" ]

	// List of box variables, each variable object contains
	// the parameters: type, name and value (plus a few more).
	Variables []BoxVariable `json:"variables"`

	// Creation date
	Created Timestamp `json:"created"`

	// Date of the last update.
	Updated *Timestamp `json:"updated,omitempty"`

	// Deleted is non-null if the Box has been deleted.
	Deleted *Timestamp `json:"deleted,omitempty"`

	// FIXME: not clear what LifeSpan is, or what it is used for.
	Lifespan *LifeSpan `json:"lifespan,omitempty"`

	// Box URI
	URI *URI `json:"uri,omitempty"` // e.g. "/services/boxes/e0715702-cf5c-4c88-bfa1-2e5e3808e597"

	// Box schema URI.
	Schema URI `json:"schema"` // e.g. "http://elasticbox.net/schemas/boxes/script"

	// List of Box members.
	Members []WorkSpaceMember `json:"members"`

	// Organization to which the box belongs: one of 'public' or 'elasticbox'
	Organization string `json:"organization"` // e.g. "elasticbox"

	// Box owner.
	Owner string `json:"owner"` // e.g. "elasticbox"

	// ID of the box version that this box is a draft from
	DraftFrom *uuid.UUID `json:"draft_from,omitempty"`

	// Map of box events
	Events map[BoxEvent]Event `json:"events,omitempty"`

	// Profile contains cloud-specific details.
	Profile *Profile `json:"profile,omitempty"`

	ProviderID *uuid.UUID `json:"provider_id,omitempty"`

	Services []Service `json:"services,omitempty"`

	Template *BlobResponse `json:"template,omitempty"`

	// Type seems to be often empty
	Type string `json:"type,omitempty"` // e.g. "CloudFormation Service"

	// BoxVersion seems to be only included when making the 'versions' API call
	BoxVersion *BoxVersion `json:"version,omitempty"`

	/*
	 * Html Section
	 */
	Readme BlobResponse `json:"readme"`

	// Box icon URI (these two seem to be mutually exclusive):
	Icon         string        `json:"icon,omitempty"`
	IconMetadata *IconMetadata `json:"icon_metadata,omitempty"`

	// More html ...
	ActionButton *ActionButton `json:"action_button,omitempty"`
}

// Members of a workspace
type WorkSpaceMember struct {
	Role      string `json:"role"`      // e.g. "collaborator"
	Workspace string `json:"workspace"` // e.g. "cf"
}

// Miscellaneous box sub-structs
type LifeSpan struct {
	Operation string `json:"operation"` // e.g. "none"
}

type IconMetadata struct {
	Border string `json:"border"`
	Fill   string `json:"fill"`
	Image  URI    `json:"image"`
}

type ActionButton struct {
	Icon  URI    `json:"icon"`
	Label string `json:"label"`
	Ref   URI    `json:"ref"`
}

// Version returns the version of @b (where defined).
func (b *Box) Version() semver.Version {
	return semver.Version{
		Major: b.BoxVersion.Number.Major,
		Minor: b.BoxVersion.Number.Minor,
		Patch: b.BoxVersion.Number.Patch,
	}
}

// BoxVersion is included in the Box struct when doing the 'versions' API call.
type BoxVersion struct {
	Box         uuid.UUID                           `json:"box"`         // e.g. "04560033-0d5c-47ed-9c77-7b13b096c172"
	Number      struct{ Major, Minor, Patch int64 } `json:"number"`      // e.g.  { "major": 0, "minor": 1, "patch": 0 }
	Workspace   string                              `json:"workspace"`   // e.g. "elasticbox"
	Description string                              `json:"description"` // e.g. "ElasticBox automatic version"
}

// IsZero returns true if @b is not initialized.
func (b BoxVersion) IsZero() bool {
	return uuid.Equal(uuid.Nil, b.Box) && b.Number.Major == 0 && b.Number.Minor == 0 && b.Number.Patch == 0
}

// GetBoxes lists all boxes that are accessible in the personal workspace of the authenticated user.
func (c *Client) GetBoxes() (res []Box, err error) {
	return res, c.Get("/services/boxes", &res)
}

// GetBox returns the details of box @boxId.
func (c *Client) GetBox(boxId string) (res Box, err error) {
	var versions []Box

	if err := c.Get(fmt.Sprintf("/services/boxes/%s/versions", boxId), &versions); err != nil {
		/* Ignore error, try the other URL below. */
	} else if len(versions) > 0 {
		return versions[0], nil
	}
	return res, c.Get("/services/boxes/"+boxId, &res)
}

// GetBoxStack returns the stack of the box @boxId.
func (c *Client) GetBoxStack(boxId string) (res []Box, err error) {
	return res, c.Get(fmt.Sprintf("/services/boxes/%s/stack", boxId), &res)
}

// BoxBinding is returned by the 'bindings' API call.
type BoxBinding struct {
	ID   uuid.UUID `json:"id"`   // e.g. "71c9a7bf-56fc-43b5-973b-0161981f4857"
	Name string    `json:"name"` // e.g. "MySQL"
	Icon URI       `json:"icon"` // e.g. "/icons/boxes/71c9a7bf-56fc-43b5-973b-0161981f4857"
	URL  URI       `json:"uri"`  // e.g. "/services/boxes/71c9a7bf-56fc-43b5-973b-0161981f4857"
}

// GetBoxBindings returns the bindings of @boxId.
func (c *Client) GetBoxBindings(boxId string) (res []BoxBinding, err error) {
	return res, c.Get(fmt.Sprintf("/services/boxes/%s/bindings", boxId), &res)
}

// GetBoxVersions returns the versions of @boxId.
func (c *Client) GetBoxVersions(boxId string) (res []Box, err error) {
	return res, c.Get(fmt.Sprintf("/services/boxes/%s/versions", boxId), &res)
}

// GetBoxDiff returns the differences of @boxId.
// FIXME: no documentation for this method and the call returns 405 (not allowed).
func (c *Client) GetBoxDiff(boxId string) error {
	return c.Get(fmt.Sprintf("/services/boxes/%s/diff", boxId), nil)
}

// UploadBox uploads @box depending on whether @boxId is not empty (create vs update).
// It returns an updated box struct on success, with fields filled in by the server.
func (c *Client) UploadBox(box *Box, boxId string) (*Box, error) {
	var res Box

	if box == nil {
		return nil, errors.Errorf("attempt to upload nil box")
	} else if boxId != "" {
		if err := c.getResponse("/services/boxes/"+boxId, "PUT", box, &res); err != nil {
			return nil, err
		}
	} else if err := c.getResponse("/services/boxes/", "POST", box, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// UploadApplianceBox performs create/update similar to UploadBox.
func (c *Client) UploadApplianceBox(box *Box, boxId string) (*Box, error) {
	var res Box

	if box == nil {
		return nil, errors.Errorf("attempt to upload nil box")
	} else if uuid.Equal(uuid.Nil, box.ID) {
		return nil, errors.Errorf("attempt to upload Appliance Box without ID")
	} else if boxUuid := box.ID.String(); boxId != "" {
		if err := c.getResponse("/services/appliance/boxes/"+boxUuid, "PUT", box, &res); err != nil {
			return nil, err
		}
	} else if err := c.getResponse("/services/appliance/boxes/"+boxUuid, "POST", box, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// DeleteBox attempts to remove box @boxId.
func (c *Client) DeleteBox(boxId string) error {
	return c.getResponse("/services/boxes/"+boxId, "DELETE", nil, nil)
}
