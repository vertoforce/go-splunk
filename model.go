package splunk

type Generator struct {
	Build   string `json:"build"`
	Version string `json:"version"`
}

type Paging struct {
	Total   int64 `json:"total"`
	PerPage int64 `json:"perPage"`
	Offset  int64 `json:"offset"`
}

type UpdateSearchConcurrencySettingsScheduleReq struct {
	// MaxSearchesPer The maximum number of searches the scheduler can run as a percentage of the maximum number of concurrent searches. Default: 50.
	MaxSearchesPer int
	// AutoSummaryPer The maximum number of concurrent searches to be allocated for auto summarization, as a percentage of the concurrent searches that the scheduler can run. Default: 50.
	AutoSummaryPer int
}
