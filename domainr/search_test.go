package domainr

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestClient_Search(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/search", func(w http.ResponseWriter, r *http.Request) {
		if want, got := "GET", r.Method; want != got {
			t.Errorf("Search() METHOD expected to be `%v`, got `%v`", want, got)
		}
		if want, got := "application/json", r.Header.Get("Accept"); want != got {
			t.Errorf("Search() Content-Type expected to be `%s`, got `%s`", want, got)
		}
		if want, got := userAgent, r.Header.Get("User-Agent"); want != got {
			t.Errorf("Search() User-Agent expected to be `%s`, got `%s`", want, got)
		}

		reqUrl := r.URL
		if want, got := "/v2/search", reqUrl.Path; want != got {
			t.Errorf("Search() /path expected to be `%s`, got `%s`", want, got)
		}
		wantQuery, _ := url.ParseQuery(fmt.Sprintf("client_id=%s&query=coffee.io", client.ClientID))
		if want, got := wantQuery, reqUrl.Query(); !reflect.DeepEqual(want, got) {
			t.Errorf("Search() ?query expected to be `%s`, got `%s`", want, got)
		}

		fmt.Fprint(w, `
			{"results":[{"domain":"coffee.io","host":"","subdomain":"coffee.","zone":"io","path":"","registerURL":"https://api.domainr.com/v2/register?client_id=3a4e71dbc14a4cc496fca2d7ee65aded\u0026domain=coffee.io\u0026registrar=\u0026source="},{"domain":"coffee.it","host":"","subdomain":"coffee.","zone":"it","path":"","registerURL":"https://api.domainr.com/v2/register?client_id=3a4e71dbc14a4cc496fca2d7ee65aded\u0026domain=coffee.it\u0026registrar=\u0026source="},{"domain":"coffee.co.it","host":"","subdomain":"coffee.","zone":"co.it","path":"","registerURL":"https://api.domainr.com/v2/register?client_id=3a4e71dbc14a4cc496fca2d7ee65aded\u0026domain=coffee.co.it\u0026registrar=\u0026source="}]}
		`)
	})

	searchResponse, err := client.Search("coffee.io", nil)
	if err != nil {
		t.Fatalf("Search() returned error: %v", err)
	}

	domains := searchResponse.Domains
	if want, got := 3, len(domains); want != got {
		t.Errorf("Search() expected to return %v domains, got %v", want, got)
	}

	var wantDomain *Domain
	wantDomain = &Domain{Name: "coffee.io", Host: "", Subdomain: "coffee.", Zone: "io", Path: "", RegisterURL: "https://api.domainr.com/v2/register?client_id=3a4e71dbc14a4cc496fca2d7ee65aded\u0026domain=coffee.io\u0026registrar=\u0026source="}
	if !reflect.DeepEqual(&domains[0], wantDomain) {
		t.Fatalf("Search() returned %+v, want %+v", domains[0], *wantDomain)
	}
}
