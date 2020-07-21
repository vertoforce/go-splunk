# Splunk Enterprise API Go Library

[![Go Report Card](https://goreportcard.com/badge/github.com/vertoforce/go-splunk)](https://goreportcard.com/report/github.com/vertoforce/go-splunk)
[![Documentation](https://godoc.org/github.com/vertoforce/go-splunk?status.svg)](https://godoc.org/github.com/vertoforce/go-splunk)

Go library to interact with the [splunk enterprise API](https://docs.splunk.com/Documentation/Splunk/8.0.5/RESTREF/RESTprolog)

I did not implement each API call, but I made a helper to make it easy to call any other endpoint

## Implemented

* [x] Create Search Job
* [x] Find Search Job
* [x] Wait on Search Job
* [x] Get Results from Search Job

## Custom API Call

If you want to make your own API call (that isn't implemented as a function), do the following:

```go
// Create the client
client, _ := splunk.NewClient(ctx, username, password, baseURL)

// Make your own custom request
resp, _ := client.BuildResponse(ctx, "GET", "/saved/searches/{name}/history", map[string]string{"savedsearch": "nameOfSearch"})

// Parse the result
result := YourStruct{}
json.NewDecoder(resp.Body).Decode(&result)
```
