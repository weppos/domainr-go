package domainr

import (
	"fmt"
	"os"
	"testing"
)

var (
	domainrLiveTest      bool
	domainrClientID      string
	domainrStatusDomains string
)

func init() {
	domainrClientID = os.Getenv("DOMAINR_CLIENT_ID")
	domainrStatusDomains = os.Getenv("DOMAINR_STATUS_DOMAINS")
	if len(domainrClientID) > 0 {
		domainrLiveTest = true
	}
	if len(domainrStatusDomains) == 0 {
		domainrStatusDomains = "domainr.com"
	}
}

func TestLiveGetStatus(t *testing.T) {
	if !domainrLiveTest {
		t.Skip("skipping live test")
	}

	client := NewClient(domainrClientID)

	statusResponse, err := client.GetStatus(domainrStatusDomains)
	fmt.Println(err)
	fmt.Println(statusResponse)
}

func TestLiveGetSingleStatus(t *testing.T) {
	if !domainrLiveTest {
		t.Skip("skipping live test")
	}

	client := NewClient(domainrClientID)
	var domain *Domain

	domain, err := client.GetSingleStatus(domainrStatusDomains)
	fmt.Println(err)
	fmt.Println(domain)
}
