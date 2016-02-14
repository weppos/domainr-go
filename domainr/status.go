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

// Status performs a /status request and returns the results.
func (c *Client) Status(domains string) (*StatusResponse, error) {
	req, err := c.NewRequest(fmt.Sprintf("/v2/status?domain=%s", domains))
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

// Status is a shortcut to checks the status of a domains and get the domains contained in the response.
func Status(client *Client, domains string) ([]Domain, error) {
	statusResponse, err := client.Status(domains)
	return statusResponse.Domains, err
}

// Status is a shortcut to checks the status of a single domains
func SingleStatus(client *Client, domain string) (*Domain, error) {
	statusResponse, err := client.Status(domain)
	if err != nil {
		return nil, err
	}

	if len(statusResponse.Domains) < 1 {
		return nil, errors.New("status returned 0 domains")
	}

	return &statusResponse.Domains[0], nil
}
