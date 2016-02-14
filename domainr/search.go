package domainr

import (
	"fmt"
	"net/http"
	"net/url"
)

// StatusResponse represents the response from a /search API call.
type SearchResponse struct {
	httpResponse *http.Response

	Query   string   `json:"query"`
	Domains []Domain `json:"results"`
}

// SearchOptions represents the parameters you can pass to a /search API requet.
// Query is not included because it is mandatory, it's not an option.
type SearchOptions struct {
	Location  string
	Registrar string
	Defaults  string
}

func (o *SearchOptions) append(query url.Values) {
	if o.Location != "" {
		query.Add("location", o.Location)
	}
	if o.Registrar != "" {
		query.Add("registrar", o.Location)
	}
	if o.Defaults != "" {
		query.Add("defaults", o.Location)
	}
}

// Status performs a /status request and returns the results.
func (c *Client) Search(query string, options *SearchOptions) (*SearchResponse, error) {
	qs := url.Values{}
	qs.Add("query", query)
	if options != nil {
		options.append(qs)
	}

	req, err := c.NewRequest(fmt.Sprintf("/v2/search?%s", qs.Encode()))
	if err != nil {
		return nil, err
	}

	searchResponse := &SearchResponse{}
	resp, err := c.Do(req, searchResponse)
	if err != nil {
		return nil, err
	}

	searchResponse.httpResponse = resp
	return searchResponse, nil
}
