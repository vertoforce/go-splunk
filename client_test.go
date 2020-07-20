package splunk

import (
	"context"
	"crypto/tls"
	"fmt"
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
	c, err := NewClient(context.Background(), os.Getenv("USERNAME"), os.Getenv("PASSWORD"), &Config{BaseURL: os.Getenv("BASEURL"), HTTPClient: client})
	if err != nil {
		t.Error(err)
		return
	}

	// Try creating a simple job
	searchID, err := c.CreateSearchJob(context.Background(), `* TEST`)
	if err != nil {
		t.Error(err)
		return
	}

	jobDetails, err := c.GetSearchJob(context.Background(), searchID)
	if err != nil {
		t.Error(err)
		return
	}
	if len(jobDetails.Entry) == 0 {
		t.Errorf("Not enough entries returned")
	}
	fmt.Printf("%s: %s\n", jobDetails.Entry[0].Name, jobDetails.Entry[0].Content.Search)
}
