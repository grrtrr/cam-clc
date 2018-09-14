package cmd

type Foo []struct {
	Services []struct {
		AutomaticReconfiguration bool   `json:"automatic_reconfiguration"`
		AutomaticUpdates         string `json:"automatic_updates"`
		Box                      struct {
			ID        string `json:"id"`
			Latest    bool   `json:"latest"`
			Variables []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"variables"`
		} `json:"box"`
		ManagedOs bool   `json:"managed_os"`
		Name      string `json:"name"`
		Policy    struct {
			Requirements []string      `json:"requirements"`
			Variables    []interface{} `json:"variables"`
		} `json:"policy"`
		Tags []string `json:"tags"`
	} `json:"services"`

	Type string `json:"type"`
}
