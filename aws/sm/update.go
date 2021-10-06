package sm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func (s *Secret) ToJSON() (string, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return "", fmt.Errorf("failed to marshal secret: %q -> %v", s.Name, err)
	}
	return string(b), nil
}

func StoreSecrets(sess aws.Config, sec []*Secret, kms string) error {
	for _, f := range sec {
		err := StoreSecret(sess, f, kms)
		if err != nil {
			return fmt.Errorf("failed to store key: %v", err)
		}
	}
	return nil
}

func StoreSecret(cfg aws.Config, s *Secret, kmsKeyID string) error {
	secret, err := s.ToJSON()
	if err != nil {
		return err
	}
	if len(s.ID) == 0 {
		return createSecret(cfg, s, kmsKeyID, secret)
	}
	return updateSecret(cfg, s, kmsKeyID, secret)
}

func StoreValue(cfg aws.Config, s *Secret, kmsKeyID string) error {
	if len(s.ID) == 0 {
		return createSecret(cfg, s, kmsKeyID, s.Value)
	}
	return updateSecret(cfg, s, kmsKeyID, s.Value)
}

func createSecret(cfg aws.Config, s *Secret, kmsKeyID string, value string) error {
	svc := secretsmanager.NewFromConfig(cfg)

	_, err := svc.CreateSecret(context.TODO(), &secretsmanager.CreateSecretInput{
		Name:         aws.String(s.Name),
		KmsKeyId:     aws.String(kmsKeyID),
		SecretString: aws.String(value),
	})
	if err != nil {
		return fmt.Errorf("failed to create secret with name: %q -> %v", s.Name, err)
	}
	return nil
}

func updateSecret(cfg aws.Config, s *Secret, kmsKeyID string, value string) error {
	svc := secretsmanager.NewFromConfig(cfg)

	_, err := svc.UpdateSecret(context.TODO(), &secretsmanager.UpdateSecretInput{
		SecretId:     aws.String(s.ID),
		KmsKeyId:     aws.String(kmsKeyID),
		SecretString: aws.String(value),
	})
	if err != nil {
		return fmt.Errorf("failed to update secret with name: %q -> %v", s.Name, err)
	}
	return nil
}
