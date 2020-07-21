package splunk

import "time"

const (
	// SplunkTimeFormat is the time format splunk expects from time parameters
	SplunkTimeFormat = "2006-01-02T15:04:05.000-07:00"
)

// FormatTime Will format a time correct to be sent as a parameter
func FormatTime(time time.Time) string {
	return time.UTC().Format(SplunkTimeFormat)
}

// ParseTime Parses a splunk returned time
func ParseTime(timeString string) (time.Time, error) {
	return time.Parse(SplunkTimeFormat, timeString)
}
