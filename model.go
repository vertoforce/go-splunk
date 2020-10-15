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
	MaxSearchesPer *int
	AutoSummaryPer *int
}

func (u *UpdateSearchConcurrencySettingsScheduleReq) SetMaxSearchesPer(value int) *UpdateSearchConcurrencySettingsScheduleReq {
	u.MaxSearchesPer = &value
	return u
}

func (u *UpdateSearchConcurrencySettingsScheduleReq) SetAutoSummaryPer(value int) *UpdateSearchConcurrencySettingsScheduleReq {
	u.AutoSummaryPer = &value
	return u
}
