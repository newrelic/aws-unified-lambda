# Integration Tests with LocalStack

## Overview

Integration tests for the Lambda function verify its behavior during invocation, including data fetching, processing, and pushing to New Relic. Configuration of triggers and permissions related to triggers are beyond the scope of these tests.

## Prerequisites

1. install `localstack` using 
   ```
   brew install localstack
   ```

2. Docker must be running. 
3. Run `localstack start`
4. In the `/integration-tests` directory, create a `.env` file with following content:

    ```shell
    NEW_RELIC_INGEST_KEY=<YOUR_KEY>
    NEW_RELIC_ACCOUNT_ID=<ACCOUNT_ID>
    NEW_RELIC_LICENSE_KEY_SECRET_NAME=nr-license-key
    NEW_RELIC_USER_KEY=<YOUR_KEY>
    ```

5. In the `/integration-tests` directory, run `go test -v`
6. If you encounter a compiler error, `compiler "x86_64-linux-musl-gcc" not found` then the cross-compilation toolchain using:

    ```shell
    brew install filosottile/musl-cross/musl-cross
    export PATH="/usr/local/bin:$PATH"
    ```

## Test Structure
```
aws-unified-logging
├── src/
│   ├── main.go
│   ├── go.mod
│   └── go.sum
└── integration-tests/
    ├── common (contants and other common files)
    ├── helpers (different resource creation files)
    ├── test-files (all the log files used for testing)
    ├── .env (temporary env variables files)
    ├── helpers.go (calls all helper functions)
    ├── main_test.go
    ├── go.mod
    └── go.sum
