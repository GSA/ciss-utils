package sm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func Pivot(secretID string) func(context.Context) (aws.Credentials, error) {
	return func(context.Context) (aws.Credentials, error) {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return aws.Credentials{}, fmt.Errorf("failed to load configuration: %v", err)
		}
		svc := secretsmanager.NewFromConfig(cfg)
		output, err := svc.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secretID),
		})
		if err != nil {
			return aws.Credentials{}, fmt.Errorf("failed to read secret value for secret: %s -> %v", secretID, err)
		}

		cred := struct {
			KeyID     string `json:"aws_access_key_id"`
			SecretKey string `json:"aws_secret_access_key"`
		}{}
		err = json.Unmarshal([]byte(aws.ToString(output.SecretString)), &cred)
		if err != nil {
			return aws.Credentials{}, fmt.Errorf("failed to unmarshal secret value for secret: %s -> %v", secretID, err)
		}

		return aws.Credentials{
			AccessKeyID:     cred.KeyID,
			SecretAccessKey: cred.SecretKey,
		}, nil
	}
}
