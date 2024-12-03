package helpers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	integrationtests "github.com/newrelic/aws-unified-lambda-logging/common"
)

// FetchLogsFromNewRelic function to fetch logs
func FetchLogsFromNewRelic(userKey, accountID, fileName string) ([]integrationtests.LogEvent, error) {
	// setup
	accountIDTemp, _ := strconv.ParseInt(accountID, 10, 32)
	nrqlQuery := "SELECT * FROM Log SINCE " + integrationtests.FetchLogsTimeRange + " ago WHERE " + integrationtests.LogObjectKey + " LIKE '" + fileName + "'"
	variables := map[string]interface{}{
		"id":   int(accountIDTemp),
		"nrql": nrqlQuery,
	}
	query := `query($id: Int!, $nrql: Nrql!) {
				actor {
					account(id: $id) {
						nrql(query: $nrql) {
							results
						}
					}
				}
			}`

	body := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", integrationtests.QueryEndpoint, bytes.NewBuffer(jsonBody))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("API-Key", userKey)

	// execute
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	// Read and output the response
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
		return nil, err
	}

	results, err := fetchEventsFromLogs(bodyBytes)
	if err != nil {
		log.Fatalf("Error to unmarshal data: %v", err)
		return nil, err
	}

	return results, nil
}

func fetchEventsFromLogs(jsonData []byte) ([]integrationtests.LogEvent, error) {
	var response integrationtests.APIResponse
	err := json.Unmarshal(jsonData, &response)
	if err != nil {
		return nil, err
	}
	return response.Data.Actor.Account.Nrql.Results, nil
}
