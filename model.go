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
