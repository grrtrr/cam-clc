package clccam

import (
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
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
		Created   Timestamp  `json:"created"`   // e.g. "2018-08-29 14:52:47.423508"
		Event     InstanceOp `json:"event"`     // e.g. "deploy"
		Workspace string     `json:"workspace"` // e.g. "gerritrenker"
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
		Box string `json:"box"`

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
	ID           string        `json:"id"`           // e.g. "eb-e775t"
	Type         string        `json:"type"`         // e.g. "Linux Compute"
	ClcAlias     string        `json:"clc_alias"`    // e.g. "AVCR"
	Organization string        `json:"organization"` // e.g. "centurylink"
	ProviderID   uuid.UUID     `json:"provider_id"`  // e.g. "8c50965d-4fd0-481a-b161-eff9fab52e51"
	Created      Timestamp     `json:"created"`      // e.g. "2018-08-29 14:52:47.357210"
	Updated      Timestamp     `json:"updated"`      // e.g. "2018-08-29 14:55:20.850380"
	Deleted      interface{}   `json:"deleted"`      // e.g. null
	Operation    InstanceOp    `json:"operation"`    // e.g. "deploy"
	State        InstanceState `json:"state"`        // e.g. "done"

	Machines []struct {
		Address        MachineAddress `json:"address"`
		AgentVersion   semver.Version `json:"agent_version"`    // e.g. "6.11"
		ExternalID     string         `json:"external_id"`      // e.g. "i-000727071fbfcec16"
		Hostname       string         `json:"hostname"`         // e.g. ""
		LastAgentClose Timestamp      `json:"last_agent_close"` // e.g. "2018-08-29 14:55:20.850348"
		LastAgentPing  Timestamp      `json:"last_agent_ping"`  // e.g. "2018-08-29 14:53:39.801080"
		Name           string         `json:"name"`             // e.g. "cms1-eb-e775t-1"
		Schema         string         `json:"schema"`           // e.g. "http://elasticbox.net/schemas/aws/service-machine"
		State          InstanceState  `json:"state"`            // e.g. "done"
		SupportID      string         `json:"support_id"`       // e.g. "AWSUSW29417"
		Token          uuid.UUID      `json:"token"`            // e.g. "b7445ed9-a4ba-4e93-9462-703ab2bd700e"
	} `json:"machines"`
	StateHistory []struct {
		State     string    `json:"state"`     // e.g. "up"
		Started   Timestamp `json:"started"`   // e.g. "2018-09-04 19:53:33.008055"
		Completed Timestamp `json:"completed"` // e.g. "2018-09-04 20:05:48.840864"
	} `json:"state_history"`

	Profile   Profile       `json:"profile"`
	Tags      []interface{} `json:"tags"`      // e.g. []
	Variables []interface{} `json:"variables"` // e.g. []

	// Token seems to be the JTI of the service token used by the service
	Token  uuid.UUID `json:"token"`  // e.g. "e80053df-fb0c-432f-a54d-ff540a6902c7"
	Schema URI       `json:"schema"` // e.g. "http://elasticbox.net/schemas/service"

	Icon string `json:"icon"` // e.g. "images/platform/linux.png"
}

// MachineAddress contains IP address information of a (Virtual) Machine.
type MachineAddress struct {
	Private string  `json:"private"` // e.g.  "172.31.1.161"
	Public  *string `json:"public"`  // e.g. null
}

func (m MachineAddress) String() string {
	if m.Public != nil && *m.Public != "" {
		return fmt.Sprintf("%s (%s)", *m.Public, m.Private)
	} else if m.Private == "" {
		return "n/a"
	}
	return m.Private
}

// GetInstanceService fetches service details of @instanceId.
func (c *Client) GetInstanceService(instanceId string) (res InstanceService, err error) {
	return res, c.Get(fmt.Sprintf("/services/instances/%s/service", instanceId), &res)
}

// InstanceActivity represents an activity log of an instance.
type InstanceActivity struct {
	Box       string    `json:"box"`        // e.g. "",
	Created   Timestamp `json:"created"`    // e.g. "2018-09-17 19:04:41.534077"
	Finished  Timestamp `json:"finished"`   // e.g. "2018-09-17 19:04:43.557408"
	Event     BoxEvent  `json:"event"`      // e.g. "install"
	ExitCode  int64     `json:"exit_code"`  // e.g. 0
	Level     string    `json:"level"`      // Seems to be an enum, e.g. waiting, start, info, install
	Machine   string    `json:"machine"`    // e.g. "cms-eb-k599a-1"
	RequestID uuid.UUID `json:"request_id"` // e.g. "0b0f6c28-c776-4f5e-9142-50711af5fda6"
	Text      string    `json:"text"`       // e.g. "Create in Progress for cmsebk599a1Machine (AWS::EC2::Instance)."
}

// GetInstanceActivity retrieves activity logs of @instanceId.
// @op: optional operation to filter by; either an empty string or a valid InstanceOp
func (c *Client) GetInstanceActivity(instanceId, op string) (res []InstanceActivity, err error) {
	var filter string

	if op != "" {
		if _, err := InstanceOpFromString(op); err != nil {
			//			return nil, err
		}
		filter = fmt.Sprintf("?operation=%s", op)
	}
	return res, c.Get(fmt.Sprintf("/services/instances/%s/activity%s", instanceId, filter), &res)
}

// GetInstanceMachineLogs retrieves the logs of machine @machineId on instance @instanceId.
func (c *Client) GetInstanceMachineLogs(instanceId, machineId string) (res string, err error) {
	return res, c.Get(fmt.Sprintf("/services/instances/%s/machine_logs?machine_name=%s", instanceId, machineId), &res)
}

// InstanceBinding is returned by the instance-binding API call.
type InstanceBinding struct {
	Box struct {
		Name    string    `json:"name"`    // e.g. "Wordpress"
		Version uuid.UUID `json:"version"` // e.g. "1186ce20-96f2-458d-ae62-19496036b275"
	} `json:"box"`
	Created      Timestamp `json:"created"`       // e.g. "2014-03-21 18:20:04.921745"
	DefaultStamp float64   `json:"default_stamp"` // e.g. 1395426004.921258
	ID           string    `json:"id"`            // e.g. "--profile id--", FIXME: may be uuid.UUID
	Instances    []struct {
		Bindings []interface{} `json:"bindings"` // e.g. []
		Box      struct {
			Name    string    `json:"name"`    // e.g. "Wordpress"
			Version uuid.UUID `json:"version"` // e.g. "1186ce20-96f2-458d-ae62-19496036b275"
		} `json:"box"`
		Path    string `json:"path"` // e.g. "/"
		Profile struct {
			Autoscalable  bool   `json:"autoscalable"`   // e.g. false
			Cloud         string `json:"cloud"`          // e.g. "EC2"
			Flavor        string `json:"flavor"`         // e.g. "t1.micro"
			Image         string `json:"image"`          // e.g. "Linux Compute"
			Instances     int64  `json:"instances"`      // e.g. 1
			Keypair       string `json:"keypair"`        // e.g. "None"
			Location      string `json:"location"`       // e.g. "us-east-1"
			Schema        URI    `json:"schema"`         // e.g. "http://elasticbox.net/schemas/aws/ec2/profile"
			SecurityGroup string `json:"security_group"` // e.g. "Automatic"
			Subnet        string `json:"subnet"`         // e.g. "us-east-1b"
		} `json:"profile"`
		Provider  string          `json:"provider"`  // e.g. "Amazon"
		Variables []BasicVariable `json:"variables"` // e.g. [ { "type":"Port", "name":"http", "value":"80" } ]
	} `json:"instances"`
	Members []string  `json:"members"` // e.g. [ "member1","member2" ]
	Name    string    `json:"name"`    // e.g. "profile"
	Owner   string    `json:"owner"`   // e.g. "workspace1"
	Schema  URI       `json:"schema"`  // e.g. "http://elasticbox.net/schemas/deployment-profile"
	Updated Timestamp `json:"updated"` // e.g. "2014-03-21 18:20:04.921745"
	URI     URI       `json:"uri"`     // e.g. "--profile uri--"
}

// GetInstanceBindings retrieves bindings of @instanceId.
func (c *Client) GetInstanceBindings(instanceId string) (res []InstanceBinding, err error) {
	return res, c.Get(fmt.Sprintf("/services/instances/%s/bindings", instanceId), &res)
}

// InstanceOperation represents operations recorded for an instance.
type InstanceOperation struct {
	Activity      []InstanceActivity `json:"activity"`       // e.g. []
	Created       Timestamp          `json:"created"`        // e.g. "2018-09-17 19:03:24.562476"
	Deleted       interface{}        `json:"deleted"`        // e.g. null
	ID            string             `json:"id"`             // e.g. "0e8d7950-1dc2-4c7d-9213-f4a9623f93dd"
	Instance      string             `json:"instance"`       // e.g. "i-dq17mn"
	Operation     InstanceOp         `json:"operation"`      // e.g. "deploy"
	RequestID     uuid.UUID          `json:"request_id"`     // e.g. "678c71bf-ce2b-4c04-a9fd-6e378a62973b"
	Schema        URI                `json:"schema"`         // e.g. "http://elasticbox.net/schemas/instance-operation"
	State         InstanceState      `json:"state"`          // e.g. "processing"
	InstanceState InstanceState      `json:"instance_state"` // e.g. "processing"
	Updated       Timestamp          `json:"updated"`        // e.g. "2018-09-17 19:04:45.333961"
	Username      string             `json:"username"`       // e.g. "gerrit.renker@centurylink.com"
	Workspace     string             `json:"workspace"`      // e.g. "gerritrenker"
}

// GetInstanceOperations retrieves operations of the machine@ @instanceId.
func (c *Client) GetInstanceOperations(instanceId string) (res []InstanceOperation, err error) {
	return res, c.Get(fmt.Sprintf("/services/instances/%s/operations", instanceId), &res)
}

/*
 * Instance Activities:
 */

// DeployInstance re-deploys an existing instance @instanceId.
func (c *Client) DeployInstance(instanceId string) error {
	return c.getResponse(fmt.Sprintf("/services/instances/%s/deploy", instanceId), "PUT", nil, nil)
}

// PowerOnInstance powers @instanceId on.
func (c *Client) PowerOnInstance(instanceId string) error {
	return c.getResponse(fmt.Sprintf("/services/instances/%s/poweron", instanceId), "PUT", nil, nil)
}

// ShutdownInstance shuts down @instanceId.
func (c *Client) ShutdownInstance(instanceId string) error {
	return c.getResponse(fmt.Sprintf("/services/instances/%s/shutdown", instanceId), "PUT", nil, nil)
}

// ReinstallInstance re-installs @instanceId.
func (c *Client) ReinstallInstance(instanceId string) error {
	return c.getResponse(fmt.Sprintf("/services/instances/%s/reinstall", instanceId), "PUT", nil, nil)
}

// ReconfigureInstance re-configures @instanceId.
func (c *Client) ReconfigureInstance(instanceId string) error {
	return c.getResponse(fmt.Sprintf("/services/instances/%s/reconfigure", instanceId), "PUT", struct {
		// FIXME: not sure the body is needed, since the information is all in the URL already.
		Id     string `json:"id"`
		Method string `json:"method"`
	}{Id: instanceId, Method: "reconfigure"}, nil)
}

// ImportInstance attempts to (re-)import an unregistered instance @instanceId.
func (c *Client) ImportInstance(instanceId string) error {
	return c.getResponse(fmt.Sprintf("/services/instances/%s/import", instanceId), "PUT", nil, nil)
}

// CancelImportInstance cancels a failed import of an unregistered instance @instanceId.
func (c *Client) CancelImportInstance(instanceId string) error {
	return c.getResponse(fmt.Sprintf("/services/instances/%s/cancel_import", instanceId), "PUT", nil, nil)
}

// MakeManagedInstance delegates management of an existing instance @instanceId to CenturyLink.
func (c *Client) MakeManagedInstance(instanceId string) error {
	return c.getResponse(fmt.Sprintf("/services/instances/%s/make_managed_os?accept_terms=true", instanceId), "PUT", nil, nil)
}

// DeleteInstance attempts to terminate / force-terminate, or delete @instanceId.
func (c *Client) DeleteInstance(instanceId, op string) error {
	switch op {
	case "terminate", "force_terminate", "delete":
		return c.getResponse(fmt.Sprintf("/services/instances/%s?operation=%s", instanceId, op), "DELETE", nil, nil)
	}
	return errors.Errorf("invalid operation %q", op)
}
