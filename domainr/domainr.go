// Package domainr implements a client for the Domainr API.
//
// In order to use this package you will need a Domainr client_id.
package domainr

import (
	"net/http"
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
