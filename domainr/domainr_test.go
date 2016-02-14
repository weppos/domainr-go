package domainr

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	mux    *http.ServeMux
	client *Client
	server *httptest.Server

	domainrLiveTest bool
	domainrClientID string
)

func init() {
	domainrClientID = os.Getenv("DOMAINR_CLIENT_ID")
	if len(domainrClientID) > 0 {
		domainrLiveTest = true
	}
}

func setupMockServer() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewClient("client-id")
	client.BaseURL = server.URL
}

func teardownMockServer() {
	server.Close()
}

func TestNewClient(t *testing.T) {
	clientID := "client-id"
	client := NewClient("client-id")

	if client.ClientID != "client-id" {
		t.Errorf("NewClient ClientID = %v, want %v", client.ClientID, clientID)
	}
}

func Test_NewRequest(t *testing.T) {
	client := NewClient("client-id")

	inURL, outURL := "foo", "https://api.domainr.com/foo?client_id=client-id"
	req, _ := client.NewRequest(inURL)

	// only GET is allowed
	if method := req.Method; method != "GET" {
		t.Fatalf("NewRequest method = %v, want GET", method)
	}

	// test that relative URL was properly built
	if req.URL.String() != outURL {
		t.Errorf("NewRequest(%v) URL = %v, want %v", inURL, req.URL, outURL)
	}

	// test that the user-agent is attached to the request
	ua := req.Header.Get("User-Agent")
	if ua != userAgent {
		t.Errorf("NewRequest User-Agent = %v, want %v", ua, userAgent)
	}
}

// can't really reproduce nd generate a bad URL.
//func Test_NewRequest_BadURL(t *testing.T) {
//	client := NewClient("client-id")
//
//	req, err := client.NewRequest("// bad url")
//	if err == nil {
//		t.Fatalf("NewRequest expected to return error, but none returned: %v", req.URL)
//	}
//}
