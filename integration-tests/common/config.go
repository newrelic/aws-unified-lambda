package common

import "time"

const (
	// LocalEndpoint endpoint where localstack is running
	LocalEndpoint = "http://localhost:4566"
	// LambdaHandler name of the lambda handler
	LambdaHandler = "main"
	// LambdaRuntime type of lambda runtime
	LambdaRuntime = "go1.x"
	// QueryEndpoint query endpoint to fetch logs from new relic
	QueryEndpoint = "https://api.newrelic.com/graphql"
	// ZipFileName executable zip file name
	ZipFileName = "main.zip"
	// ExecutableFileName executable file name
	ExecutableFileName = "main"
	// NewRelicRegion region where lambda is deployed
	NewRelicRegion = "US"
	// LogObjectKey parameter of entity synthesis
	LogObjectKey = "logObjectKey"
	// FetchLogsTimeRange time range used in query to fetch logs
	FetchLogsTimeRange = "30 seconds"
	// NoOfRetriesForEventCreation number of retries while creating resources
	NoOfRetriesForEventCreation = 3
	// NoOfRetriesForLogsFetch number of retries while fetching logs
	NoOfRetriesForLogsFetch = 3
	// WaitTimeForResourceCreation wait time before resource creation
	WaitTimeForResourceCreation = 5 * time.Second
	// WaitTimeForLogsFetch wait time before fetching logs again
	WaitTimeForLogsFetch = 2 * time.Second
	// RetryError error code for retry
	RetryError = "400"
)

// SecretData structure to store license key json
type SecretData struct {
	LicenseKey string `json:"LicenseKey"`
}

// LogEvent entities/attributes in logs fetched
type LogEvent struct {
	AwsAccountID            string `json:"aws.accountId"`
	AwsRealm                string `json:"aws.realm"`
	AwsRegion               string `json:"aws.region"`
	InstrumentationName     string `json:"instrumentation.name"`
	InstrumentationProvider string `json:"instrumentation.provider"`
	LogBucketName           string `json:"logBucketName"`
	LogObjectKey            string `json:"logObjectKey"`
	Message                 string `json:"message"`
	NewRelicLogPattern      string `json:"newrelic.logPattern"`
	NewRelicSource          string `json:"newrelic.source"`
	Timestamp               int64  `json:"timestamp"`
}

type nrql struct {
	Results []LogEvent `json:"results"`
}

type account struct {
	Nrql nrql `json:"nrql"`
}

type actor struct {
	Account account `json:"account"`
}

type data struct {
	Actor actor `json:"actor"`
}

// APIResponse api response from new relic logs fetch
type APIResponse struct {
	Data data `json:"data"`
}
