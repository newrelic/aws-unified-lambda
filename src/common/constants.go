// Package common provides common constants structs and variables.
package common

// InstrumentationProvider is a parameter necessary for Entity Synthesis at New Relic.
const InstrumentationProvider = "aws"

// InstrumentationName is a parameter necessary for Entity Synthesis at New Relic.
const InstrumentationName = "lambda"

// CustomMetaData is the name of the environment variable for custom meta data.
const CustomMetaData = "CUSTOM_META_DATA"

// NewRelicLicenseKeySecretName is the name of the environment variable for the New Relic license key secret.
const NewRelicLicenseKeySecretName = "NEW_RELIC_LICENSE_KEY_SECRET_NAME"

// EnvLicenseKey is the name of the environment variable for the license key.
const EnvLicenseKey = "LICENSE_KEY"

// LicenseKey is the name of the license key.
const LicenseKey = "LicenseKey"

// NewRelicRegion is the name of the environment variable for the New Relic region.
const NewRelicRegion = "NEW_RELIC_REGION"

// DebugEnabled is the name of the environment variable for enabling debug mode.
const DebugEnabled = "DEBUG_ENABLED"

// MaxBufferSize is the maximum buffer size used to read buffer readers. This is the maximum size of a log, any message larger than this will cause an error.
const MaxBufferSize = 8 * 1024 * 1024 // 8 mb

// MaxMessageSize is the maximum size of a message. Any message larger than this will be split into multiple records.
const MaxMessageSize = 1 * 1024 * 1024 // 1 mb

// MaxPayloadSize is the maximum size of a payload.
// Reference: https://docs.newrelic.com/docs/logs/log-api/introduction-log-api/#limits
const MaxPayloadSize = 1 * 1024 * 1024 // 1 mb

// MaxPayloadMessages is the maximum number of messages in a payload.
const MaxPayloadMessages = 900

// CloudTrailDigestRegex is the regex pattern for CloudTrail digest files.
const CloudTrailDigestRegex = ".*_CloudTrail-Digest_.*\\.json\\.gz$"

// CloudTrailRegex is the regex pattern for CloudTrail files.
const CloudTrailRegex = ".*_CloudTrail_.*\\.json\\.gz$"

// RequestIDRegex is the regex pattern for RequestId.
const RequestIDRegex = "RequestId:\\s([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})"

// LambdaLogGroup is prefix for identifing log group belonging to lambda
const LambdaLogGroup = "/aws/lambda"
