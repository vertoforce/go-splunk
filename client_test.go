package splunk

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"testing"
)

func TestClientSimple(t *testing.T) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	_, err := NewClient(context.Background(), os.Getenv("USERNAME"), os.Getenv("PASSWORD"), &Config{BaseURL: os.Getenv("BASEURL"), HTTPClient: client})
	if err != nil {
		t.Error(err)
		return
	}
}
