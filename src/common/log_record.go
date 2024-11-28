package common

// DetailedLog represents a detailed log record.
//
// Reference: https://docs.newrelic.com/docs/logs/log-api/introduction-log-api/#detailed-json
type DetailedLog struct {
	CommonData Common  `json:"common"`
	Entries    LogData `json:"logs"`
}

// Common represents the common data shared by all log records.
type Common struct {
	Attributes LogAttributes `json:"attributes"` // Optional
	Timestamp  string        `json:"timestamp"`  // Optional
}

// LogData represents a collection of log records.
type LogData []Log

// Log represents a single log record.
type Log struct {
	Timestamp  string        `json:"timestamp"`  // Optional : If not timestamp is present per log and also in the common attribute then the log ingestion time is used
	Attributes LogAttributes `json:"attributes"` // Optional
	Log        string        `json:"log"`        // Required
}

// LogAttributes represents the attributes of a log record.
type LogAttributes map[string]interface{}

// DetailedLogsBatch represents a batch of detailed log records. This is the expected payload format in the API call to New Relic.
type DetailedLogsBatch []DetailedLog
