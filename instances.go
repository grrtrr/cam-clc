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
	AutomaticReconfiguration bool          `json:"automatic_reconfiguration"` // e.g. true
	AutomaticUpdates         string        `json:"automatic_updates"`         // e.g. "off"
	Bindings                 []interface{} `json:"bindings"`                  // e.g. []
	Box                      uuid.UUID     `json:"box"`                       // e.g. "37cb9262-d04f-4bf0-97d7-4429d2bad6c3"
	Boxes                    []struct {
		AutomaticUpdates string     `json:"automatic_updates"` // e.g. "off"
		Created          Timestamp  `json:"created"`           // e.g. "2018-05-23 22:20:07.916219"
		Deleted          *Timestamp `json:"deleted"`           // e.g. null
		DraftFrom        uuid.UUID  `json:"draft_from"`        // e.g. "e211992f-04a7-4d8d-ae23-1b38bd22be65"
		Events           struct {
			Install struct {
				ContentType     string `json:"content_type"`     // e.g. "text/x-shellscript"
				DestinationPath string `json:"destination_path"` // e.g. "scripts"
				Length          int    `json:"length"`           // e.g. 2368
				URL             string `json:"url"`              // e.g. "/services/blobs/download/5b8568651873ed2f4dd461b1/install"
			} `json:"install"`
			PreInstall struct {
				ContentType     string `json:"content_type"`     // e.g. "text/x-shellscript"
				DestinationPath string `json:"destination_path"` // e.g. "scripts"
				Length          int    `json:"length"`           // e.g. 1153
				URL             string `json:"url"`              // e.g. "/services/blobs/download/5b0c5cd11873ed47f11cc0ff/pre_install"
			} `json:"pre_install"`
		} `json:"events"`
		Icon         string `json:"icon"` // e.g. "/icons/boxes/37cb9262-d04f-4bf0-97d7-4429d2bad6c3"
		IconMetadata struct {
			Border string `json:"border"` // e.g. "#027333"
			Fill   string `json:"fill"`   // e.g. "#03924e"
			Image  string `json:"image"`  // e.g. "/services/blobs/download/5b7b4110759a1802697fa3b6/SAHA5.png"
		} `json:"icon_metadata"`
		ID      uuid.UUID `json:"id"` // e.g. "37cb9262-d04f-4bf0-97d7-4429d2bad6c3"
		Members []struct {
			Role      string `json:"role"`      // e.g. "collaborator", "read"
			Workspace string `json:"workspace"` // e.g. "safehaven"
		} `json:"members"`
		Name         string `json:"name"`         // e.g. "CMS"
		Organization string `json:"organization"` // e.g. "centurylink"
		Owner        string `json:"owner"`        // e.g. "sh"
		Readme       struct {
			ContentType string    `json:"content_type"` // e.g. "text/x-markdown"
			Length      int       `json:"length"`       // e.g. 194
			UploadDate  Timestamp `json:"upload_date"`  // e.g. "2018-05-24 16:26:53.996892"
			URL         string    `json:"url"`          // e.g. "/services/blobs/download/5b06e7cd159b89785d75b5d0/README.md"
		} `json:"readme"`
		Requirements []string  `json:"requirements"` // e.g. [ "safehaven-cms" ]
		Schema       string    `json:"schema"`       // e.g. "http://elasticbox.net/schemas/boxes/script
		Updated      Timestamp `json:"updated"`      // e.g. "2018-08-28 15:21:09.895659"
		Variables    []struct {
			Name       string `json:"name"` // e.g. "DEB_URL"
			Required   bool   `json:"required"`
			Type       string `json:"type"`       // e.g. "Text", "Password", "Port"
			Value      string `json:"value"`      // e.g. "http://10.55.220.31/downloads/safehaven-nightly.deb"
			Visibility string `json:"visibility"` // e.g. "internal", "private"
		} `json:"variables"`
		Visibility string `json:"visibility"` // e.g. "workspace"
	} `json:"boxes"`
	Created      Timestamp     `json:"created"`     // e.g. "2018-08-29 14:52:47.335395"
	Deleted      *Timestamp    `json:"deleted"`     // e.g. null
	Description  string        `json:"description"` // Descriptive text of this deployment
	ID           string        `json:"id"`          // e.g. "i-z48wub"
	IsDeployOnly bool          `json:"is_deploy_only"`
	Members      []interface{} `json:"members"`
	Name         string        `json:"name"` // Name of the deployment
	Operation    struct {
		Created   Timestamp `json:"created"`   // e.g. "2018-08-29 14:52:47.423508"
		Event     string    `json:"event"`     // e.g. "deploy"
		Workspace string    `json:"workspace"` // e.g. "gerritrenker"
	} `json:"operation"`
	Owner     string `json:"owner"` // e.g. "gerritrenker"
	PolicyBox struct {
		AutomaticUpdates string     `json:"automatic_updates"` // e.g. "off"
		Claims           []string   `json:"claims"`            // e.g. [ "safehaven-cms" ]
		Created          Timestamp  `json:"created"`           // e.g. "2018-08-21 15:08:55.409444"
		Deleted          *Timestamp `json:"deleted"`           // e.g. null
		ID               uuid.UUID  `json:"id"`                // e.g. "f54d2970-ff97-4a28-8c12-f6b6fa1a00dd"
		Lifespan         struct {
			Operation string `json:"operation"` // e.g. "none"
		} `json:"lifespan"`
		Members      []interface{} `json:"members"`
		Name         string        `json:"name"`         // e.g. "CMS"
		Organization string        `json:"organization"` // e.g. "centurylink"
		Owner        string        `json:"owner"`        // e.g. "gerritrenker"
		Profile      struct {
			Cloud          string `json:"cloud"`           // e.g. "vpc-0e77a66a",
			ElasticIP      bool   `json:"elastic_ip"`      // e.g. false
			Flavor         string `json:"flavor"`          // e.g. "c1.medium"
			Image          string `json:"image"`           // e.g. ami-96d129ee"
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
			Subnet         string        `json:"subnet"`          // e.g. "subnet-e3336a95",
			Volumes        []interface{} `json:"volumes"`         // e.g. []
		} `json:"profile"`
		ProviderID uuid.UUID     `json:"provider_id"` // e.g. "8c50965d-4fd0-481a-b161-eff9fab52e51"
		Schema     string        `json:"schema"`      // e.g. "http://elasticbox.net/schemas/boxes/policy"
		Updated    Timestamp     `json:"updated"`     // e.g. "2018-08-21 15:12:08.486998"
		Variables  []interface{} `json:"variables"`   // e.g. []
		Visibility string        `json:"visibility"`  // e.g. "workspace"
	} `json:"policy_box"`
	PricingHistory []struct {
		From        Timestamp `json:"from"` // e.g. "2018-08-29 14:52:47.329214"
		PricingInfo struct {
			EstimatedMonthly int    `json:"estimated_monthly"` // e.g. 9360000
			Factor           int    `json:"factor"`            // e.g. 100000
			HourlyPrice      int    `json:"hourly_price"`      // e.g. 13000
			ProviderType     string `json:"provider_type"`     // e.g. "Amazon Web Services"
		} `json:"pricing_info"`
	} `json:"pricing_history"`
	Schema  string `json:"schema"` // e.g. "http://elasticbox.net/schemas/instance"
	Service struct {
		ID       string `json:"id"` // e.g. "eb-e775t"
		Machines []struct {
			Name     string        `json:"name"`     // e.g. "cms1-eb-e775t-1"
			State    string        `json:"state"`    // e.g. "done"
			Workflow []interface{} `json:"workflow"` // e.g. []
		} `json:"machines"`
		Type string `json:"type"` // e.g. "Linux Compute"
	} `json:"service"`
	State     string        `json:"state"`     // e.g. "done"
	Tags      []string      `json:"tags"`      // e.g. [ "myTag" ]
	Updated   Timestamp     `json:"updated"`   // e.g. "2018-08-29 14:55:23.604873"
	URI       string        `json:"uri"`       // e.g. "/services/instances/i-z48wub"
	Variables []interface{} `json:"variables"` // e.g. []
}

// GetInstance retrieves details of @instanceId
func (c *Client) GetInstance(instanceId string) (res Instance, err error) {
	return res, c.Get("/services/instances/"+instanceId, &res)
}

// Service represents the service associated with an instance.
type InstanceService struct {
	ClcAlias string     `json:"clc_alias"` // e.g. "AVCR"
	Created  Timestamp  `json:"created"`   // e.g. "2018-08-29 14:52:47.357210"
	Deleted  *Timestamp `json:"deleted"`   // e.g. null
	Icon     string     `json:"icon"`      // e.g. "images/platform/linux.png"
	ID       string     `json:"id"`        // e.g. "eb-e775t"
	Machines []struct {
		Address struct {
			Private string  `json:"private"` // e.g.  "172.31.1.161"
			Public  *string `json:"public"`  // e.g. null
		} `json:"address"`
		AgentVersion   string    `json:"agent_version"`    // e.g. "6.11"
		ExternalID     string    `json:"external_id"`      // e.g. "i-000727071fbfcec16"
		Hostname       string    `json:"hostname"`         // e.g. ""
		LastAgentClose Timestamp `json:"last_agent_close"` // e.g. "2018-08-29 14:55:20.850348"
		LastAgentPing  Timestamp `json:"last_agent_ping"`  // e.g. "2018-08-29 14:53:39.801080"
		Name           string    `json:"name"`             // e.g. "cms1-eb-e775t-1"
		Schema         string    `json:"schema"`           // e.g. "http://elasticbox.net/schemas/aws/service-machine"
		State          string    `json:"state"`            // e.g. "done"
		SupportID      string    `json:"support_id"`       // e.g. "AWSUSW29417"
		Token          uuid.UUID `json:"token"`            // e.g. "b7445ed9-a4ba-4e93-9462-703ab2bd700e"
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
	ProviderID   uuid.UUID `json:"provider_id"` // e.g. "8c50965d-4fd0-481a-b161-eff9fab52e51"
	Schema       string    `json:"schema"`      // e.g. "http://elasticbox.net/schemas/service"
	State        string    `json:"state"`       // e.g. "done"
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
