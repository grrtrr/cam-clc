package clccam

// GetInstances lists all instances
func (c *Client) GetInstances() ([]Box, error) {
	var res []Box

	return res, c.Get("/services/instances", &res)
}
