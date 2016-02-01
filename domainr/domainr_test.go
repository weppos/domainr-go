package domainr

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	clientID := "client-id"
	client := NewClient("client-id")

	if client.ClientID != "client-id" {
		t.Errorf("NewClient ClientID = %v, want %v", client.ClientID, clientID)
	}
}
