package util

import (
	"context"
	"github.com/newrelic/aws-unified-lambda-logging/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
	logging "github.com/newrelic/newrelic-client-go/v2/pkg/logs"
	"github.com/newrelic/newrelic-client-go/v2/pkg/region"
	"os"
	"sync"
)

// NewRelicClientAPI is an interface that defines the methods for interacting with the New Relic Logs API.
type NewRelicClientAPI interface {
	CreateLogEntry(logEntry interface{}) error
}

// ConsumeLogBatches consumes log batches from a channel and creates log entries using the provided NewRelicClientAPI.
// The function returns when the channel is closed or the context is cancelled.
func ConsumeLogBatches(ctx context.Context, channel <-chan common.DetailedLogsBatch, wg *sync.WaitGroup, nrClientAPI NewRelicClientAPI) {
	// Defer the Done() method of the WaitGroup to indicate that the goroutine has finished processing
	defer wg.Done()

	for {
		select {
		case batch, ok := <-channel:
			if !ok {
				return
			}
			if err := nrClientAPI.CreateLogEntry(batch); err != nil {
				log.Fatalf("error posting Log entry: %v", err)
			}
		case <-ctx.Done():
			// Context has been cancelled, exit the goroutine
			return
		}
	}
}

// NewNRClient Initializes a new NRClient with debug level and region
// It returns a NewRelicClientAPI interface and an error if there is a problem setting the region.
func NewNRClient() (NewRelicClientAPI, error) {
	nrRegion, _ := region.Get(region.Name(os.Getenv(common.NewRelicRegion)))
	var nrClient logging.Logs
	cfg := config.Config{
		Compression: config.Compression.Gzip,
	}

	if os.Getenv(common.DebugEnabled) == "true" {
		cfg.LogLevel = "debug"
	} else {
		cfg.LogLevel = "info"
	}

	if err := cfg.SetRegion(nrRegion); err != nil {
		return &nrClient, err
	}

	licenseKey, err := GetLicenseKey()
	cfg.LicenseKey = licenseKey
	nrClient = logging.New(cfg)
	return &nrClient, err
}
