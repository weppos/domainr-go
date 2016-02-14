// Package domainr implements a client for the Domainr API.
//
// In order to use this package you will need a Domainr client_id.
package domainr

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	// libraryVersion identifies the current library version.
	// This is a pro-forma convention given that Go dependencies
	// tends to be fetched directly from the repo.
	// It is also used in the user-agent identify the client.
	libraryVersion = "0.2.0"

	// userAgent represents the user agent used
	// when communicating with the Domainr API.
	userAgent = "domainr-go/" + libraryVersion

	domainrBaseURL = "https://api.domainr.com/"
	domainrParam   = "client_id"
	mashapeBaseURL = "https://domainr.p.mashape.com/"
	mashapeParam   = "mashape-url"
)

// Client represents a client to the Domainr API.
type Client struct {
	// HttpClient is the underlying HTTP client
	// used to communicate with the API.
	HttpClient *http.Client

	// Authenticator used for authentication.
	auth Authenticator

	// BaseURL for API requests.
	BaseURL string
}

// Domain represents a domain name in Domainr
// and is the result of either a search or status query.
type Domain struct {
	// Shared fields
	Name string `json:"domain"`
	Zone string `json:"zone"`

	// Status fields
	Status  string `json:"status"`
	Summary string `json:"summary"`

	// Search fields
	Host        string `json:"host"`
	Subdomain   string `json:"subdomain"`
	Path        string `json:"path"`
	RegisterURL string `json:"registerURL"`
}

// The Domainr API requires authentication for all requests.
// The authenticator interface exposes the requirement for an authentication implementation.
type Authenticator interface {
	Param() (string, string)
	Endpoint() (string)
}

type authentication struct {
	key      string
	value    string
	endpoint string
}

// Settings implements Authenticator.
func (a *authentication) Param() (string, string) {
	return a.key, a.value
}

// Endpoint implements Authenticator.
func (a *authentication) Endpoint() string {
	return a.endpoint
}

// Mashape users will use your Mashape API key
// and the mashape-key= query parameter to authenticate.
func NewMashapeAuthentication(clientID string) Authenticator {
	return &authentication{key: mashapeParam, value: clientID, endpoint: mashapeBaseURL}
}

// High-volume commercial users can use a client_id parameter to authenticate.
// You need to contact Domainr to obtain an API key.
func NewDomainrAuthentication(clientID string) Authenticator {
	return &authentication{key: domainrParam, value: clientID, endpoint: domainrBaseURL}
}

// NewClient returns a new Domainr API client.
func NewClient(auth Authenticator) *Client {
	client := &Client{auth: auth, BaseURL: auth.Endpoint(), HttpClient: &http.Client{}}
	return client
}

// NewRequest creates an API request.
// The path is expected to be a relative path and will be resolved according to the BaseURL of the Client.
// Paths should always be specified without a preceding slash.
//
// Domainr only requires GET requests, therefore this method assumes you want GET
// and also doesn't provide a way to specify a payload.
func (c *Client) NewRequest(path string) (*http.Request, error) {
	u, err := url.Parse(c.BaseURL + path)
	if err != nil {
		return nil, fmt.Errorf("Error parsing request URL: %s", err)
	}

	// Append the client_id to the query
	q := u.Query()
	authKey, authValue := c.auth.Param()
	q.Add(authKey, authValue)
	u.RawQuery = q.Encode()

	// Build the request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", userAgent)

	return req, nil
}

// Do sends an API request and returns the API response.
func (c *Client) Do(req *http.Request, obj interface{}) (*http.Response, error) {
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = c.checkResponse(resp)
	if err != nil {
		return resp, err
	}

	// If obj implements the io.Writer,
	// the response body is decoded into v.
	if obj != nil {
		if w, ok := obj.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(obj)
		}
	}

	return resp, err
}

type errorMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (c *Client) checkResponse(resp *http.Response) error {
	switch resp.StatusCode {
	case 200:
		return nil
	// {"message":"unauthorized: invalid API client ID"}
	// {"status":"404","message":"No results found."}
	default:
		message := errorMessage{}
		body, err := ioutil.ReadAll(resp.Body)
		if err = json.Unmarshal(body, &message); err != nil {
			return err
		}

		return errors.New(message.Message)
	}
}
