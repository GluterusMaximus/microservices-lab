package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseKafkaMessage(t *testing.T) {
	message := []byte(`{"ID": 123, "RequestorEmail": "test@example.com", "Status": "pending"}`)
	order, err := ParseKafkaMessage(message)

	if err != nil {
		t.Errorf("ParseKafkaMessage returned an error: %v", err)
	}

	if order.ID != 123 || order.RequestorEmail != "test@example.com" || order.Status != "pending" {
		t.Errorf("ParseKafkaMessage did not return the expected order")
	}
}

func TestSendEmailViaBrevo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "https://api.brevo.com/v3/smtp/email", r.URL.String())
		// More assertions on request headers, etc.
	}))
	defer server.Close()

	apiKey := "test-api-key"
	order := Order{ID: 123, RequestorEmail: "test@example.com", Status: "pending"}

	err := SendEmailViaBrevo(order, apiKey)
	assert.NoError(t, err)
}
