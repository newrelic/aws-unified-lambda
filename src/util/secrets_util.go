package util

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/newrelic/aws-unified-lambda-logging/common"
	"os"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/newrelic/aws-unified-lambda-logging/logger"
)

var log = logger.NewLogrusLogger(logger.WithDebugLevel())

// SecretManagerAPI is an interface for interacting with AWS Secrets Manager.
type SecretManagerAPI interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

// GetSecretFromSecretManager retrieves a secret from AWS Secrets Manager.
// It returns a map of secret data and an error if any.
func GetSecretFromSecretManager(ctx context.Context, secretManagerClient SecretManagerAPI, secretName string) (map[string]string, error) {
	// Check if the passed secret name is empty
	if secretName == "" {
		return nil, errors.New("secret name is empty")
	}

	// Fetch the response from the provided secrets manager client
	resp, err := secretManagerClient.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	})
	if err != nil {
		log.WithField("error", err).Error("secrets manager GetSecretValue error")
		return nil, err
	}
	log.Debug("successfully fetched secret from secret manager")

	// Decode the JSON response string from AWS Secrets Manager
	var secretData map[string]string
	if resp.SecretString != nil {
		if err := json.Unmarshal([]byte(*resp.SecretString), &secretData); err != nil {
			return nil, err
		}
		return secretData, nil
	}
	return nil, errors.New("secret is in binary format (likely encrypted)")
}

// NewSecretsManagerClient creates a new AWS Secrets Manager client.
// It returns a SecretManagerAPI client and an error if any.
func NewSecretsManagerClient() (SecretManagerAPI, error) {
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.WithField("error", err).Error("aws configuration couldn't be found")
		return nil, err
	}
	return secretsmanager.NewFromConfig(cfg), nil
}

// GetLicenseKey accepts returns the license key from the environment variable or the AWS Secrets manager.
// It returns the New Relic Ingest License key and an error if any.
func GetLicenseKey() (key string, err error) {
	if os.Getenv(common.EnvLicenseKey) != "" {
		log.Debugf("fetching license key from environment variable")
		return os.Getenv(common.EnvLicenseKey), nil
	}
	log.Debugf("fetching license key from secret manager")
	secretsManagerClient, err := NewSecretsManagerClient()

	if err != nil {
		return "", err
	}
	secretMap, err := GetSecretFromSecretManager(context.TODO(), secretsManagerClient, os.Getenv(common.NewRelicLicenseKeySecretName))
	if err != nil {
		return "", err
	}
	if secretMap[common.LicenseKey] != "" {
		return secretMap[common.LicenseKey], nil
	}
	return "", errors.New("either LicenseKey is empty or not present in the secrets manager")
}
