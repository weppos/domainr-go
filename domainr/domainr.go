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
	libraryVersion = "0.0.1"

	// baseURL is the Domainr API URL.
	baseURL = "https://api.domainr.com/"

	// userAgent represents the user agent used
	// when communicating with the Domainr API.
	userAgent = "domainr-go/" + libraryVersion
)

// Client represents a client to the Domainr API.
type Client struct {
	// HttpClient is the underlying HTTP client
	// used to communicate with the API.
	HttpClient *http.Client

	// ClientID token used for authentication.
	ClientID string
}

// NewClient returns a new Domainr API client.
func NewClient(clientID string) *Client {
	client := &Client{ClientID: clientID, HttpClient: &http.Client{}}
	return client
}

// NewRequest creates an API request.
// The path is expected to be a relative path and will be resolved according to the BaseURL of the Client.
// Paths should always be specified without a preceding slash.
//
// Domainr only requires GET requests, therefore this method assumes you want GET
// and also doesn't provide a way to specify a payload.
func (c *Client) NewRequest(path string) (*http.Request, error) {
	u, err := url.Parse(baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("Error parsing request URL: %s", err)
	}

	// Append the client_id to the query
	q := u.Query()
	q.Add("client_id", c.ClientID)
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
	Message string `json:"message"`
}

func (c *Client) checkResponse(resp *http.Response) error {
	switch resp.StatusCode {
	case 200:
		return nil
	// {"message":"unauthorized: invalid API client ID"}
	default:
		message := errorMessage{}
		body, err := ioutil.ReadAll(resp.Body)
		if err = json.Unmarshal(body, &message); err != nil {
			return err
		}

		return errors.New(message.Message)
	}
}
