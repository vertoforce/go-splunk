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
	searchJob, err := c.CreateSearchJob(context.Background(), `TEST`, nil)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("Job %s created\n", searchJob.SearchID)

	jobDetails, err := c.GetSearchJob(context.Background(), searchJob.SearchID)
	if err != nil {
		t.Error(err)
		return
	}
	if len(jobDetails.Entry) == 0 {
		t.Errorf("Not enough entries returned")
	}

	// Wait on job
	// searchJob.Wait(context.Background())
	// fmt.Println("Job done")
	// time.Sleep(time.Millisecond * 50)

	// Try getting results
	searchContext, cancel := context.WithCancel(context.Background())
	results, err := searchJob.GetResults(searchContext)
	totalResults := 0
	for result := range results {
		totalResults++
		fmt.Printf("%d %s\n", totalResults, result["_time"])
		if totalResults >= 50 {
			cancel()
			break
		}
	}
	cancel()
	fmt.Printf("Total results: %d\n", totalResults)

	// Cancel and remove job
	err = searchJob.Delete(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
}
