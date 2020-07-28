package splunk

import "testing"

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
