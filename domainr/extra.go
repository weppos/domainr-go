package domainr

import (
	"net/http"
)

// ZonesResponse represents the response from a /zone API call.
type ZonesResponse struct {
	httpResponse *http.Response

	Zones []string `json:"zones"`
}

// Zone performs a /zone request and returns the results.
func (c *Client) Zones() (*ZonesResponse, error) {
	req, err := c.NewRequest("/v2/zones")
	if err != nil {
		return nil, err
	}

	zonesResponse := &ZonesResponse{}
	resp, err := c.Do(req, zonesResponse)
	if err != nil {
		return nil, err
	}

	zonesResponse.httpResponse = resp
	return zonesResponse, nil
}
