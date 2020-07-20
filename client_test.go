package splunk

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
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
	fmt.Printf("Job %s created\n", searchID)

	jobDetails, err := c.GetSearchJob(context.Background(), searchID)
	if err != nil {
		t.Error(err)
		return
	}
	if len(jobDetails.Entry) == 0 {
		t.Errorf("Not enough entries returned")
	}

	// Wait on job
	c.WaitOnJob(context.Background(), searchID)
	fmt.Println("Job done")
	time.Sleep(time.Millisecond * 50)

	// Try getting results
	results, err := c.GetSearchJobResults(context.Background(), searchID)
	totalResults := 0
	for range results {
		totalResults++
	}
	fmt.Printf("Total resutls: %d\n", totalResults)
}
