package domainr

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"testing"
)

var (
	domainrStatusDomains string
)

func init() {
	domainrStatusDomains = os.Getenv("DOMAINR_STATUS_DOMAINS")
	if len(domainrStatusDomains) == 0 {
		domainrStatusDomains = "domainr.com"
	}
}

func TestClient_Status(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/status", func(w http.ResponseWriter, r *http.Request) {
		if want, got := "GET", r.Method; want != got {
			t.Errorf("Status() METHOD expected to be `%v`, got `%v`", want, got)
		}
		if want, got := "application/json", r.Header.Get("Accept"); want != got {
			t.Errorf("Status() Content-Type expected to be `%s`, got `%s`", want, got)
		}
		if want, got := userAgent, r.Header.Get("User-Agent"); want != got {
			t.Errorf("Status() User-Agent expected to be `%s`, got `%s`", want, got)
		}

		reqUrl := r.URL
		if want, got := "/v2/status", reqUrl.Path; want != got {
			t.Errorf("Status() /path expected to be `%s`, got `%s`", want, got)
		}
		wantQuery, _ := url.ParseQuery(fmt.Sprintf("client_id=%s&domain=example.com,example.org", "client-id"))
		if want, got := wantQuery, reqUrl.Query(); !reflect.DeepEqual(want, got) {
			t.Errorf("Status() ?query expected to be `%s`, got `%s`", want, got)
		}

		fmt.Fprint(w, `
			{"status":[{"domain":"example.com","zone":"com","status":"active registrar","summary":"active"},{"domain":"example.org","zone":"org","status":"active","summary":"active"}]}
		`)
	})

	statusResponse, err := client.Status([]string{"example.com", "example.org"})
	if err != nil {
		t.Fatalf("Status() returned error: %v", err)
	}

	domains := statusResponse.Domains
	if want, got := 2, len(domains); want != got {
		t.Errorf("Status() expected to return %v domains, got %v", want, got)
	}

	var wantDomain *Domain
	wantDomain = &Domain{Name: "example.com", Zone: "com", Status: "active registrar", Summary: "active"}
	if !reflect.DeepEqual(&domains[0], wantDomain) {
		t.Fatalf("Status() returned %+v, want %+v", domains[0], *wantDomain)
	}
}

func Test_Status(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/status", func(w http.ResponseWriter, r *http.Request) {
		if want, got := "GET", r.Method; want != got {
			t.Errorf("Status() METHOD expected to be `%v`, got `%v`", want, got)
		}

		fmt.Fprint(w, `
			{"status":[{"domain":"example.com","zone":"com","status":"active registrar","summary":"active"},{"domain":"example.org","zone":"org","status":"active","summary":"active"}]}
		`)
	})

	domains, err := Status(client, []string{"example.com", "example.org"})
	if err != nil {
		t.Fatalf("Status() returned error: %v", err)
	}
	if want, got := 2, len(domains); want != got {
		t.Errorf("Status() expected to return %v domains, got %v", want, got)
	}

	wantDomain := &Domain{Name: "example.com", Zone: "com", Status: "active registrar", Summary: "active"}
	if !reflect.DeepEqual(&domains[0], wantDomain) {
		t.Fatalf("GetStatus() returned %+v, want %+v", domains[0], *wantDomain)
	}
}

func TestLive_Client_Status(t *testing.T) {
	if !domainrLiveTest {
		t.Skip("skipping live test")
	}

	client := NewClient(NewDomainrAuthentication(domainrClientID))

	statusResponse, err := client.Status([]string{domainrStatusDomains})
	fmt.Println(err)
	fmt.Println(statusResponse)
}

func TestLive_Status(t *testing.T) {
	if !domainrLiveTest {
		t.Skip("skipping live test")
	}

	client := NewClient(NewDomainrAuthentication(domainrClientID))
	var domains []Domain

	domains, err := Status(client, []string{domainrStatusDomains})
	fmt.Println(err)
	fmt.Println(domains)
}

func TestLive_SingleStatus(t *testing.T) {
	if !domainrLiveTest {
		t.Skip("skipping live test")
	}

	client := NewClient(NewDomainrAuthentication(domainrClientID))
	var domain *Domain

	domain, err := SingleStatus(client, domainrStatusDomains)
	fmt.Println(err)
	fmt.Println(domain)
}
