package splunk

import (
	"fmt"
	"testing"
	"time"
)

func TestMarshal(t *testing.T) {
	type TestStruct struct {
		String string    `splunk:"string"`
		Int    int       `splunk:"int"`
		Time   time.Time `splunk:"time"`
	}

	searchResult := SearchResult{
		"string": "stringVal",
		"int":    float64(3),
		"time":   "2006-01-02T15:04:00.000+00:00",
	}

	testStruct := TestStruct{}
	err := searchResult.UnMarshal(&testStruct)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(testStruct)
}
