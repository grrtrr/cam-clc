package clccam

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

/*
 * Instances API
 * See https://www.ctl.io/api-docs/cam/#application-lifecycle-management-instances-api
 */

// Instance represents a single instance of a service.
type Instance struct {
	// Creation date
	Created Timestamp `json:"created"`

	// Date of termination (not always present)
	Terminated Timestamp `json:"terminated"`

	// Date of last update
	Updated Timestamp `json:"updated"`

	// List of members that are sharing this instance
	Members []interface{} `json:"members"`

	// Instance ID
	ID string `json:"id"` // e.g. "i-z48wub"

	// Instance URI corresponding to @ID
	URI string `json:"uri"` // e.g. "/services/instances/i-z48wub"

	// Instance name
	Name string `json:"name"`

	// Instance owner
	Owner string `json:"owner"` // e.g. "gerritrenker"

	// Last operation.
	Operation struct {
		Created   Timestamp     `json:"created"`   // e.g. "2018-08-29 14:52:47.423508"
		Event     InstanceEvent `json:"event"`     // e.g. "deploy"
		Workspace string        `json:"workspace"` // e.g. "gerritrenker"
	} `json:"operation"`

	// Instance service:
	Service struct {
		// Service unique identifier
		ID string `json:"id"` // e.g. "eb-e775t"

		// List of service machines:
		Machines []Machine `json:"machines"`

		// Type is a required field, which can be one of
		// - Linux Compute
		// - Windows Compute
		// - CloudFormation Service
		Type string `json:"type"`
	} `json:"service"`

	// List of instance tags
	Tags []string `json:"tags"` // e.g. [ "myTag" ]

	// Instance state
	State InstanceState `json:"state"` // e.g. "done"

	AutomaticReconfiguration bool `json:"automatic_reconfiguration"`

	AutomaticUpdates string `json:"automatic_updates"` // e.g. "off"

	// List of instance bindings
	Bindings []struct {
		Instance string `json:"instance"`
		Name     string `json:"name"`
	} `json:"bindings"`

	Box uuid.UUID `json:"box"` // e.g. "37cb9262-d04f-4bf0-97d7-4429d2bad6c3"

	// List of boxes
	Boxes     []Box `json:"boxes"`
	PolicyBox Box   `json:"policy_box"`

	Deleted     interface{} `json:"deleted"`     // e.g. null
	Description string      `json:"description"` // Descriptive text of this deployment

	IsDeployOnly bool `json:"is_deploy_only"`

	PricingHistory []struct {
		From        Timestamp          `json:"from"`
		PricingInfo PricingInformation `json:"pricing_info"`
	} `json:"pricing_history"`

	// Instance schema URI
	Schema URI `json:"schema"` // e.g. "http://elasticbox.net/schemas/instance"

	Variables []interface{} `json:"variables"` // e.g. []
}

// Machine is used to describe a VM within a service
type Machine struct {
	// Machine name
	Name string `json:"name"` // e.g. "cms1-eb-e775t-1"

	// Machine state
	State InstanceState `json:"state"` // e.g. "done"

	// List of workflow actions:
	Workflow []struct {
		// Workflow action box
		Box Box `json:"box"`

		// Workflow action event
		Event string `json:"event"`

		// Workflow action script URI
		Script URI `json:"script"`
	} `json:"workflow"`
}

// GetInstance retrieves details of @instanceId
func (c *Client) GetInstance(instanceId string) (res Instance, err error) {
	return res, c.Get("/services/instances/"+instanceId, &res)
}

// GetInstances returns a list of all instances owned by the token user.
func (c *Client) GetInstances() (res []Instance, err error) {
	return res, c.Get("/services/instances", &res)
}

// Service represents the service associated with an instance.
type InstanceService struct {
	ClcAlias string      `json:"clc_alias"` // e.g. "AVCR"
	Created  Timestamp   `json:"created"`   // e.g. "2018-08-29 14:52:47.357210"
	Deleted  interface{} `json:"deleted"`   // e.g. null
	Icon     string      `json:"icon"`      // e.g. "images/platform/linux.png"
	ID       string      `json:"id"`        // e.g. "eb-e775t"
	Machines []struct {
		Address struct {
			Private string  `json:"private"` // e.g.  "172.31.1.161"
			Public  *string `json:"public"`  // e.g. null
		} `json:"address"`
		AgentVersion   string        `json:"agent_version"`    // e.g. "6.11"
		ExternalID     string        `json:"external_id"`      // e.g. "i-000727071fbfcec16"
		Hostname       string        `json:"hostname"`         // e.g. ""
		LastAgentClose Timestamp     `json:"last_agent_close"` // e.g. "2018-08-29 14:55:20.850348"
		LastAgentPing  Timestamp     `json:"last_agent_ping"`  // e.g. "2018-08-29 14:53:39.801080"
		Name           string        `json:"name"`             // e.g. "cms1-eb-e775t-1"
		Schema         string        `json:"schema"`           // e.g. "http://elasticbox.net/schemas/aws/service-machine"
		State          InstanceState `json:"state"`            // e.g. "done"
		SupportID      string        `json:"support_id"`       // e.g. "AWSUSW29417"
		Token          uuid.UUID     `json:"token"`            // e.g. "b7445ed9-a4ba-4e93-9462-703ab2bd700e"
	} `json:"machines"`
	Operation    string `json:"operation"`    // e.g. "deploy"
	Organization string `json:"organization"` // e.g. "centurylink"
	Profile      struct {
		Cloud          string `json:"cloud"`           // e.g. "vpc-0e77a66a"
		ElasticIP      bool   `json:"elastic_ip"`      // e.g. false
		Flavor         string `json:"flavor"`          // e.g. "c1.medium"
		Image          string `json:"image"`           // e.g. "ami-96d129ee"
		Instances      int    `json:"instances"`       // e.g. 1
		Keypair        string `json:"keypair"`         // e.g. "None"
		Location       string `json:"location"`        // e.g. "us-west-2"
		ManagedOs      bool   `json:"managed_os"`      // e.g. false
		PlacementGroup string `json:"placement_group"` // e.g. ""
		PricingInfo    struct {
			EstimatedMonthly int    `json:"estimated_monthly"` // e.g. 9360000
			Factor           int    `json:"factor"`            // e.g. 100000
			HourlyPrice      int    `json:"hourly_price"`      // e.g. 13000
			ProviderType     string `json:"provider_type"`     // e.g. "Amazon Web Services"
		} `json:"pricing_info"`
		Role           string        `json:"role"`            // e.g. "None"
		Schema         string        `json:"schema"`          // e.g. "http://elasticbox.net/schemas/aws/ec2/profile"
		SecurityGroups []string      `json:"security_groups"` // e.g. [ "sg-7d528c04" ]
		Subnet         string        `json:"subnet"`          // e.g. "subnet-e3336a95"
		Volumes        []interface{} `json:"volumes"`         // e.g. []
	} `json:"profile"`
	ProviderID   uuid.UUID     `json:"provider_id"` // e.g. "8c50965d-4fd0-481a-b161-eff9fab52e51"
	Schema       string        `json:"schema"`      // e.g. "http://elasticbox.net/schemas/service"
	State        InstanceState `json:"state"`       // e.g. "done"
	StateHistory []struct {
		Started Timestamp `json:"started"` // e.g. "2018-08-29 14:53:31.469591"
		State   string    `json:"state"`
	} `json:"state_history"`
	Tags      []interface{} `json:"tags"`      // e.g. []
	Token     uuid.UUID     `json:"token"`     // e.g. "e80053df-fb0c-432f-a54d-ff540a6902c7"
	Type      string        `json:"type"`      // e.g. "Linux Compute"
	Updated   Timestamp     `json:"updated"`   // e.g. "2018-08-29 14:55:20.850380"
	Variables []interface{} `json:"variables"` // e.g. []
}

// GetInstanceService fetches service details of @instanceId.
func (c *Client) GetInstanceService(instanceId string) (res InstanceService, err error) {
	return res, c.Get(fmt.Sprintf("/services/instances/%s/service", instanceId), &res)
}
