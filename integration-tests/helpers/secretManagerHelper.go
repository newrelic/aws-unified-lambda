package helpers

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	integrationtests "github.com/newrelic/aws-unified-lambda-logging/common"
	"log"
)

// CreateSecret function to create secret
func CreateSecret(smClient *secretsmanager.SecretsManager, secretName, licenseKey string) error {
	secretValue, err := getSecretKeyJSON(licenseKey)

	input := &secretsmanager.CreateSecretInput{
		Name:         aws.String(secretName),
		SecretString: aws.String(secretValue),
	}

	_, err = smClient.CreateSecret(input)
	if err != nil {
		return err
	}

	return nil
}

func getSecretKeyJSON(licenseKey string) (string, error) {
	secretData := integrationtests.SecretData{
		LicenseKey: licenseKey,
	}
	secretValueBytes, err := json.Marshal(secretData)
	if err != nil {
		log.Fatalf("failed to marshal secret data: %v", err)
		return "", err
	}
	secretValue := string(secretValueBytes)

	return secretValue, nil
}
