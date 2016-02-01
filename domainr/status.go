package domainr

import (
	"errors"
	"fmt"
	"net/http"
)

// Domain represents a domain name and the corresponding
// status information as defined by Domainr.
type Domain struct {
	Name    string `json:"domain"`
	Zone    string `json:"zone"`
	Status  string `json:"status"`
	Summary string `json:"summary"`
}

type StatusResponse struct {
	Domains []Domain `json:"status"`
	// Omit "errors" for now

	httpResponse *http.Response
}

func (c *Client) getStatus(domain string) (*StatusResponse, error) {
	req, err := c.NewRequest(fmt.Sprintf("/v2/status?domain=%s", domain))
	if err != nil {
		return nil, err
	}

	statusResponse := &StatusResponse{}
	resp, err := c.Do(req, statusResponse)
	if err != nil {
		return nil, err
	}

	statusResponse.httpResponse = resp
	return statusResponse, nil
}

// GetStatus checks the status of a single domain.
func (c *Client) GetStatus(domain string) (*Domain, error) {
	statusResponse, err := c.getStatus(domain)
	if err != nil {
		return nil, err
	}

	if len(statusResponse.Domains) < 1 {
		return nil, errors.New("status returned 0 domains")
	}

	return &statusResponse.Domains[0], nil
}
