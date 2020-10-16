package splunk

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	searchJobsSuffix                = "/services/search/jobs"
	searchJobSuffix                 = "/services/search/jobs/%s"
	searchControlJobSuffix          = "/services/search/jobs/%s/control"
	searchConcurrencySettingsSuffix = "/services/search/concurrency-settings/scheduler"
)

// Search is a search in splunk
type Search struct {
	SearchID string `json:"sid"`
	client   *Client
}

// CreateSearchJob Creates a search and returns the search object
//
// Params are any other parameters you want to specific from [the documentation](https://docs.splunk.com/Documentation/Splunk/8.0.5/RESTREF/RESTsearch#search.2Fjobs)
func (c *Client) CreateSearchJob(ctx context.Context, query string, params map[string]string) (*Search, error) {
	// Build params
	paramsToSend := params
	if paramsToSend == nil {
		paramsToSend = map[string]string{}
	}
	paramsToSend["search"] = fmt.Sprintf("search %s", query)

	// Make request
	resp, err := c.BuildResponse(ctx, "POST", searchJobsSuffix, paramsToSend)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 201 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("bad status code: %d, body: %s", resp.StatusCode, string(body))
	}

	search := &Search{}
	err = json.NewDecoder(resp.Body).Decode(search)
	resp.Body.Close()
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
	DiskUsage                         float64       `json:"diskUsage"`
	DispatchState                     string        `json:"dispatchState"`
	DoneProgress                      float64       `json:"doneProgress"`
	DropCount                         float64       `json:"dropCount"`
	EarliestTime                      string        `json:"earliestTime"`
	EventAvailableCount               float64       `json:"eventAvailableCount"`
	EventCount                        float64       `json:"eventCount"`
	EventFieldCount                   float64       `json:"eventFieldCount"`
	EventIsStreaming                  bool          `json:"eventIsStreaming"`
	EventIsTruncated                  bool          `json:"eventIsTruncated"`
	EventSearch                       string        `json:"eventSearch"`
	EventSorting                      string        `json:"eventSorting"`
	IndexEarliestTime                 float64       `json:"indexEarliestTime"`
	IndexLatestTime                   float64       `json:"indexLatestTime"`
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
	NumPreviews                       float64       `json:"numPreviews"`
	OptimizedSearch                   string        `json:"optimizedSearch"`
	Phase0                            string        `json:"phase0"`
	Phase1                            string        `json:"phase1"`
	PID                               string        `json:"pid"`
	Priority                          float64       `json:"priority"`
	Provenance                        string        `json:"provenance"`
	RemoteSearch                      string        `json:"remoteSearch"`
	ReportSearch                      string        `json:"reportSearch"`
	ResultCount                       float64       `json:"resultCount"`
	ResultIsStreaming                 bool          `json:"resultIsStreaming"`
	ResultPreviewCount                float64       `json:"resultPreviewCount"`
	RunDuration                       float64       `json:"runDuration"`
	SampleRatio                       string        `json:"sampleRatio"`
	SampleSeed                        string        `json:"sampleSeed"`
	ScanCount                         float64       `json:"scanCount"`
	Search                            string        `json:"search"`
	SearchCanBeEventType              bool          `json:"searchCanBeEventType"`
	SearchTotalBucketsCount           float64       `json:"searchTotalBucketsCount"`
	SearchTotalEliminatedBucketsCount float64       `json:"searchTotalEliminatedBucketsCount"`
	Sid                               string        `json:"sid"`
	StatusBuckets                     float64       `json:"statusBuckets"`
	TTL                               float64       `json:"ttl"`
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
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %s", err)
	}

	return &result, nil
}

// DeleteSearchJob Delete current search job
func (c *Client) DeleteSearchJob(ctx context.Context, searchID string) error {
	resp, err := c.BuildResponse(ctx, http.MethodDelete, fmt.Sprintf(searchJobSuffix, searchID), nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}
	return nil
}

// RunSearchJobControlCommand Run a job control command for the {search_id} search.
func (c *Client) RunSearchJobControlCommand(ctx context.Context, searchID string, action ControlCommand) error {
	params := map[string]string{"action": string(action)}
	resp, err := c.BuildResponse(ctx, http.MethodPost, fmt.Sprintf(searchControlJobSuffix, searchID), params)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}
	return nil
}

// UpdateSearchConcurrencySettingsScheduler Edit settings that determine concurrent scheduled search limits.
func (c *Client) UpdateSearchConcurrencySettingsScheduler(ctx context.Context, req *UpdateSearchConcurrencySettingsScheduleReq) error {
	params := make(map[string]string)
	if req.MaxSearchesPer != 0 {
		params["max_searches_perc"] = strconv.Itoa(req.MaxSearchesPer)
	}
	if req.AutoSummaryPer != 0 {
		params["auto_summary_perc"] = strconv.Itoa(req.AutoSummaryPer)
	}
	resp, err := c.BuildResponse(ctx, http.MethodPost, searchConcurrencySettingsSuffix, params)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}
	return nil
}

// Wait for a search job to be done.
// It waits for the dispatchState to be "DONE".
//
// If there is an error it returns.  If no jobs is found, it returns.
//
func (s *Search) Wait(ctx context.Context) error {
	for {
		job, err := s.client.GetSearchJob(ctx, s.SearchID)
		if err != nil {
			return err
		}
		if len(job.Entry) == 0 {
			return fmt.Errorf("no search found")
		}
		if job.Entry[0].SearchContent.DispatchState == "DONE" {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second * 3):
		}
	}
}

// SearchResults is the response when fetching a single page of results
type SearchResults struct {
	Preview    bool                `json:"preview"`
	InitOffset int64               `json:"init_offset"`
	Fields     []map[string]string `json:"fields"`
	Results    []SearchResult      `json:"results"`
}

// SearchResult is a single search result from a splunk search
type SearchResult map[string]interface{}

// GetFieldString retuns the string value of the field, or "" if it does not exist
func (s SearchResult) GetFieldString(fieldName string) string {
	if val, ok := s[fieldName]; ok {
		return fmt.Sprintf("%v", val)
	}
	return ""
}

// GetResults Gets a channel of results from the search job.
//
// If the search is still running, it will get the available results, and wait for
// results to continue populating.  It will not close the channel until the search is finished
// AND it sends all results
//
// *NOTE*: If you are performing a search with changing results (like a stats command)
// you must wait for the search to complete before getting results.  Otherwise you will get available
// results that will later be changed.
func (s *Search) GetResults(ctx context.Context) (chan SearchResult, error) {
	// Number of results per page
	count := 100

	// Make results channel with 4 page buffer
	results := make(chan SearchResult, count*4)

	go func() {
		defer close(results)

		offset := 0
		for {
			params := map[string]string{
				"count":  fmt.Sprintf("%d", count),
				"offset": fmt.Sprintf("%d", offset),
			}

			resp, err := s.client.BuildResponse(ctx, "GET", fmt.Sprintf(searchJobSuffix, s.SearchID)+"/results_preview", params)
			if err != nil {
				return
			}
			if resp.StatusCode != 200 {
				// No more content
				return
			}

			result := SearchResults{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			resp.Body.Close()
			if err != nil {
				return
			}

			if len(result.Results) == 0 && !result.Preview {
				// No more results and these results aren't a preview, we are done
				return
			}

			for _, result := range result.Results {
				select {
				case results <- result:
				case <-ctx.Done():
					return
				}
			}

			if result.Preview {
				// The search is still running and we've reached the end of the available results
				// Wait a bit before making the next request so we aren't spamming when
				// there are no results
				time.Sleep(time.Second)
			}

			offset += len(result.Results)
		}
	}()

	return results, nil
}

// Delete the job in splunk and remove it.  If you stop an already stopped job, it will do nothing
func (s *Search) Delete(ctx context.Context) error {
	resp, err := s.client.BuildResponse(ctx, "DELETE", fmt.Sprintf(searchJobSuffix, s.SearchID), nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	return nil
}

// StopAndFinalize the job in splunk.
func (s *Search) StopAndFinalize(ctx context.Context) error {
	resp, err := s.client.BuildResponse(ctx, "POST", fmt.Sprintf(searchControlJobSuffix, s.SearchID), map[string]string{
		"action": "finalize",
	})
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	return nil
}

// URL Returns the human vistable URL to see the results of the search
//
// By default it will use the config base URL with the port set to 80, but you can pass in a custom base URL
func (s *Search) URL(customBaseURL ...string) string {
	baseURL := ""
	if len(customBaseURL) > 0 {
		baseURL = customBaseURL[0]
	} else {
		baseURLP, err := url.Parse(s.client.config.BaseURL)
		if err != nil {
			return ""
		}
		baseURLP.Host = baseURLP.Hostname()
		baseURL = baseURLP.String()
	}
	return fmt.Sprintf("%s/en-US/app/search/search?sid=%s", baseURL, s.SearchID)
}
