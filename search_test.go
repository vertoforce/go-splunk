package splunk

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetURL(t *testing.T) {
	s := Search{
		SearchID: "TestID",
		client: &Client{
			config: &Config{
				BaseURL: "http://localhost:8090",
			},
		},
	}

	if url := s.URL(); url != "http://localhost/en-US/app/search/search?sid=TestID" {
		t.Errorf("Bad URL: %s", url)
	}
	if url := s.URL("http://localhost:81"); url != "http://localhost:81/en-US/app/search/search?sid=TestID" {
		t.Errorf("Bad URL: %s", url)
	}
}

func TestClient_DeleteSearchJob(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			switch {
			case req.Method != http.MethodDelete:
				rw.WriteHeader(http.StatusMethodNotAllowed)
				return
			case req.RequestURI != "/services/search/jobs/job_id_1?output_mode=json":
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
			rw.WriteHeader(http.StatusOK)
		}))
		defer server.Close()
		client := &Client{
			config: &Config{
				BaseURL:    server.URL,
				HTTPClient: http.DefaultClient,
			},
		}
		err := client.DeleteSearchJob(context.Background(), "job_id_1")
		require.NoError(t, err)
	})
}

func TestClient_UpdateSearchConcurrencySettingsScheduler(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			switch {
			case req.Method != http.MethodPost:
				rw.WriteHeader(http.StatusMethodNotAllowed)
				return
			case req.RequestURI != "/services/search/concurrency-settings/scheduler":
				rw.WriteHeader(http.StatusBadRequest)
				return
			}

			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}

			if string(b) != `auto_summary_perc=60&max_searches_perc=60&output_mode=json` {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}

			rw.WriteHeader(http.StatusOK)
		}))
		defer server.Close()
		client := &Client{
			config: &Config{
				BaseURL:    server.URL,
				HTTPClient: http.DefaultClient,
			},
		}
		req := &UpdateSearchConcurrencySettingsScheduleReq{
			MaxSearchesPer: 60,
			AutoSummaryPer: 60,
		}
		err := client.UpdateSearchConcurrencySettingsScheduler(context.Background(), req)
		require.NoError(t, err)
	})
}

func TestClient_RunSearchJobControlCommand(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			switch {
			case req.Method != http.MethodPost:
				rw.WriteHeader(http.StatusMethodNotAllowed)
				return
			case req.RequestURI != "/services/search/jobs/job_id_1/control":
				rw.WriteHeader(http.StatusBadRequest)
				return
			}

			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}

			if string(b) != "action=cancel&output_mode=json" {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}

			rw.WriteHeader(http.StatusOK)
		}))
		defer server.Close()
		client := &Client{
			config: &Config{
				BaseURL:    server.URL,
				HTTPClient: http.DefaultClient,
			},
		}
		err := client.RunSearchJobControlCommand(context.Background(), "job_id_1", ControlCommandCancel)
		require.NoError(t, err)
	})
}
