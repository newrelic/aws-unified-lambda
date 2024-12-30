#!/bin/bash

TEMPLATES=(
  "logging-lambda-metric-polling.yaml"
  "logging-lambda-metric-stream.yaml"
  "logging-firehose-metric-polling.yaml"
  "logging-firehose-metric-stream.yaml"
  "logging-lambda-firehose-metric-polling.yaml"
  "logging-lambda-firehose-metric-stream.yaml"
  "lambda-template.yaml"
  "logging-lambda-firehose-template.yaml"
)

source config-file.cfg

for TEMPLATE_FILE in "${TEMPLATES[@]}"; do

  BASE_NAME=$(basename "$TEMPLATE_FILE" .yaml)
  BUILD_DIR="$BUILD_DIR_BASE/$BASE_NAME"

  # sam build and package
  sam build -u --template-file "../$TEMPLATE_FILE" --build-dir "$BUILD_DIR"
  sam package --s3-bucket "$S3_BUCKET" --template-file "$BUILD_DIR/template.yaml" --output-template-file "$BUILD_DIR/$TEMPLATE_FILE"

done

