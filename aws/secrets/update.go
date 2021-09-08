package secrets

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type secret struct {
	Username   string `json:"-"`
	KeyID      string `json:"aws_access_key_id"`
	Secret     string `json:"aws_secret_access_key"`
	SecretName string `json:"-"`
	SecretID   string `json:"-"`
}

func (s *secret) ToJSON() (string, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return "", fmt.Errorf("failed to marshal secret to JSON: %q -> %v", s.SecretName, err)
	}
	return string(b), nil
}

func StoreKeys(sec []*secret, kms string) error {
	for _, f := range sec {
		err := StoreKey(f, kms)
		if err != nil {
			return fmt.Errorf("failed to store key: %v", err)
		}
	}
	return nil
}

func StoreKey(s *secret, kms string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to load configuration: %v", err)
	}

	svc := secretsmanager.NewFromConfig(cfg)

	secret, err := s.ToJSON()
	if err != nil {
		return err
	}

	if len(s.SecretID) == 0 {
		_, err := svc.CreateSecret(context.TODO(), &secretsmanager.CreateSecretInput{
			Name:         aws.String(s.SecretName),
			KmsKeyId:     aws.String(kms),
			SecretString: aws.String(secret),
		})
		if err != nil {
			return fmt.Errorf("failed to create secret with name: %q -> %v", s.SecretName, err)
		}
		return nil
	}

	_, err = svc.UpdateSecret(context.TODO(), &secretsmanager.UpdateSecretInput{
		SecretId:     aws.String(s.SecretID),
		KmsKeyId:     aws.String(kms),
		SecretString: aws.String(secret),
	})
	if err != nil {
		return fmt.Errorf("failed to update secret with name: %q -> %v", s.SecretName, err)
	}
	return nil
}
