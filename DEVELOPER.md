## Implementation

The Lambda leverages the [New Relic Go client](https://github.com/newrelic/newrelic-client-go) to process the logs in batches. This means it converts the AWS source logs into [detailed JSON format](https://docs.newrelic.com/docs/logs/log-api/introduction-log-api/#detailed-json).

## Requirements

- AWS CLI must be installed and configured with Administrator permission
- Docker must be installed. Refer [Docker documentation.](https://www.docker.com/community-edition)
- Golang must be installed. Refer [Golang documentation](https://golang.org)
- Install the AWS SAM CLI. Refer [SAM CLI Documentation](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)


## Unit Testing and Coverage

Run the following commands to:

- Run the unit tests:

    ```shell
    make test
    ```

- Check the coverage:

    ```shell
    make coverage
    ```
