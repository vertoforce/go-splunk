package splunk

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

const (
	searchJobsSuffix = "/services/search/jobs"
	searchJobSuffix  = "/services/search/jobs/%s"
)

// Search is a search in splunk
type Search struct {
	SearchID string `json:"sid"`
	client   *Client
}

// CreateSearchJob Creates a search and returns the search object
func (c *Client) CreateSearchJob(ctx context.Context, query string) (*Search, error) {
	resp, err := c.BuildResponse(ctx, "POST", searchJobsSuffix, map[string]string{
		"search": fmt.Sprintf("search %s", query),
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	search := &Search{}
	err = json.NewDecoder(resp.Body).Decode(search)
	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to unmarshal: %s, body: %s", err, string(body))
	}
	search.client = c

	return search, nil
}

// JobSearchResult is what splunk returns when searching for a search job
type JobSearchResult struct {
	Generator Generator `json:"generator"`
	Entry     []Entry   `json:"entry"`
	Paging    Paging    `json:"paging"`
}

// Entry is a search job stored on the server is some state
type Entry struct {
	Name          string        `json:"name"`
	ID            string        `json:"id"`
	Updated       string        `json:"updated"`
	Links         interface{}   `json:"links"`
	Published     string        `json:"published"`
	Author        string        `json:"author"`
	SearchContent SearchContent `json:"content"`
}

// SearchContent is the details about a searchJob
type SearchContent struct {
	BundleVersion                     string        `json:"bundleVersion"`
	CanSummarize                      bool          `json:"canSummarize"`
	CursorTime                        string        `json:"cursorTime"`
	DefaultSaveTTL                    string        `json:"defaultSaveTTL"`
	DefaultTTL                        string        `json:"defaultTTL"`
	Delegate                          string        `json:"delegate"`
	DiskUsage                         int64         `json:"diskUsage"`
	DispatchState                     string        `json:"dispatchState"`
	DoneProgress                      int64         `json:"doneProgress"`
	DropCount                         int64         `json:"dropCount"`
	EarliestTime                      string        `json:"earliestTime"`
	EventAvailableCount               int64         `json:"eventAvailableCount"`
	EventCount                        int64         `json:"eventCount"`
	EventFieldCount                   int64         `json:"eventFieldCount"`
	EventIsStreaming                  bool          `json:"eventIsStreaming"`
	EventIsTruncated                  bool          `json:"eventIsTruncated"`
	EventSearch                       string        `json:"eventSearch"`
	EventSorting                      string        `json:"eventSorting"`
	IndexEarliestTime                 int64         `json:"indexEarliestTime"`
	IndexLatestTime                   int64         `json:"indexLatestTime"`
	IsBatchModeSearch                 bool          `json:"isBatchModeSearch"`
	IsDone                            bool          `json:"isDone"`
	IsEventsPreviewEnabled            bool          `json:"isEventsPreviewEnabled"`
	IsFailed                          bool          `json:"isFailed"`
	IsFinalized                       bool          `json:"isFinalized"`
	IsPaused                          bool          `json:"isPaused"`
	IsPreviewEnabled                  bool          `json:"isPreviewEnabled"`
	IsRealTimeSearch                  bool          `json:"isRealTimeSearch"`
	IsRemoteTimeline                  bool          `json:"isRemoteTimeline"`
	IsSaved                           bool          `json:"isSaved"`
	IsSavedSearch                     bool          `json:"isSavedSearch"`
	IsTimeCursored                    bool          `json:"isTimeCursored"`
	IsZombie                          bool          `json:"isZombie"`
	Keywords                          string        `json:"keywords"`
	Label                             string        `json:"label"`
	NormalizedSearch                  string        `json:"normalizedSearch"`
	NumPreviews                       int64         `json:"numPreviews"`
	OptimizedSearch                   string        `json:"optimizedSearch"`
	Phase0                            string        `json:"phase0"`
	Phase1                            string        `json:"phase1"`
	PID                               string        `json:"pid"`
	Priority                          int64         `json:"priority"`
	Provenance                        string        `json:"provenance"`
	RemoteSearch                      string        `json:"remoteSearch"`
	ReportSearch                      string        `json:"reportSearch"`
	ResultCount                       int64         `json:"resultCount"`
	ResultIsStreaming                 bool          `json:"resultIsStreaming"`
	ResultPreviewCount                int64         `json:"resultPreviewCount"`
	RunDuration                       float64       `json:"runDuration"`
	SampleRatio                       string        `json:"sampleRatio"`
	SampleSeed                        string        `json:"sampleSeed"`
	ScanCount                         int64         `json:"scanCount"`
	Search                            string        `json:"search"`
	SearchCanBeEventType              bool          `json:"searchCanBeEventType"`
	SearchTotalBucketsCount           int64         `json:"searchTotalBucketsCount"`
	SearchTotalEliminatedBucketsCount int64         `json:"searchTotalEliminatedBucketsCount"`
	Sid                               string        `json:"sid"`
	StatusBuckets                     int64         `json:"statusBuckets"`
	TTL                               int64         `json:"ttl"`
	Messages                          []interface{} `json:"messages"`
	SearchProviders                   []string      `json:"searchProviders"`
	RemoteSearchLogs                  []string      `json:"remoteSearchLogs"`
}

// GetSearchJob Gets details about a current search job
func (c *Client) GetSearchJob(ctx context.Context, searchID string) (*JobSearchResult, error) {
	resp, err := c.BuildResponse(ctx, "GET", fmt.Sprintf(searchJobSuffix, searchID), nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	result := JobSearchResult{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %s", err)
	}

	return &result, nil
}

// Wait for a search job to be done and the results available.
// It waits for the dispatchState to be "DONE".
//
// If there is an error it returns.  If no jobs is found, it returns.
//
func (s *Search) Wait(ctx context.Context) {
	for {
		job, err := s.client.GetSearchJob(ctx, s.SearchID)
		if err != nil {
			return
		}
		if len(job.Entry) == 0 {
			return
		}
		if job.Entry[0].SearchContent.DispatchState == "DONE" {
			return
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second * 3):
		}
	}
}

// SearchResults is the response when fetching a single page of results
type SearchResults struct {
	Preview    bool                     `json:"preview"`
	InitOffset int64                    `json:"init_offset"`
	Fields     []map[string]string      `json:"fields"`
	Results    []map[string]interface{} `json:"results"`
}

// GetResults Gets a channel of results from the search job
func (s *Search) GetResults(ctx context.Context) (chan map[string]interface{}, error) {
	count := 100

	// Make results channel with 4 page buffer
	results := make(chan map[string]interface{}, count*4)

	go func() {
		defer close(results)

		page := 0
		for {
			params := map[string]string{
				"count":  fmt.Sprintf("%d", count),
				"offset": fmt.Sprintf("%d", page*count),
			}

			resp, err := s.client.BuildResponse(ctx, "GET", fmt.Sprintf(searchJobSuffix, s.SearchID)+"/results", params)
			if err != nil {
				return
			}
			if resp.StatusCode != 200 {
				return
			}

			result := SearchResults{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			if err != nil {
				return
			}

			if len(result.Results) == 0 {
				// Done
				return
			}

			for _, result := range result.Results {
				select {
				case results <- result:
				case <-ctx.Done():
					return
				}
			}

			page++
		}
	}()

	return results, nil
}
