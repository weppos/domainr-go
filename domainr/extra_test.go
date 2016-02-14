package domainr

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestClient_Zones(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/zones", func(w http.ResponseWriter, r *http.Request) {
		if want, got := "GET", r.Method; want != got {
			t.Errorf("Zones() METHOD expected to be `%v`, got `%v`", want, got)
		}
		if want, got := "application/json", r.Header.Get("Accept"); want != got {
			t.Errorf("Zones() Content-Type expected to be `%s`, got `%s`", want, got)
		}
		if want, got := userAgent, r.Header.Get("User-Agent"); want != got {
			t.Errorf("Zones() User-Agent expected to be `%s`, got `%s`", want, got)
		}

		reqUrl := r.URL
		if want, got := "/v2/zones", reqUrl.Path; want != got {
			t.Errorf("Zones() /path expected to be `%s`, got `%s`", want, got)
		}
		wantQuery, _ := url.ParseQuery(fmt.Sprintf("client_id=%s", client.ClientID))
		if want, got := wantQuery, reqUrl.Query(); !reflect.DeepEqual(want, got) {
			t.Errorf("Zones() ?query expected to be `%s`, got `%s`", want, got)
		}

		fmt.Fprint(w, `
			{"zones":["africa","africa.com","city","한국"]}
		`)
	})

	zonesResponse, err := client.Zones()
	if err != nil {
		t.Fatalf("Zones() returned error: %v", err)
	}

	zones := zonesResponse.Zones
	if want, got := []string{"africa", "africa.com", "city", "한국"}, zones; !reflect.DeepEqual(want, got) {
		t.Errorf("Zones() expected to return %v zones, got %v", want, got)
	}
}
