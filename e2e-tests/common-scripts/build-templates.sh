#!/bin/bash

# add all templates for subsequent test cases
# make packaging and publishing of templates more efficient so every time new al2 env is not spun up
TEMPLATES=(
  "logging-lambda-firehose-template.yaml"
  "logging-lambda-metric-polling.yaml"
  "logging-lambda-metric-stream.yaml"
  "lambda-template.yaml"
)

source config-file.cfg

for TEMPLATE_FILE in "${TEMPLATES[@]}"; do

  BASE_NAME=$(basename "$TEMPLATE_FILE" .yaml)
  BUILD_DIR="$BUILD_DIR_BASE/$BASE_NAME"

  sam build -u --template-file "../../$TEMPLATE_FILE" --build-dir "$BUILD_DIR"
  sam package --s3-bucket "$S3_BUCKET" --template-file "$BUILD_DIR/template.yaml" --output-template-file "$BUILD_DIR/$TEMPLATE_FILE"

done