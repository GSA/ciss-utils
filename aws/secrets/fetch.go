package secrets

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	smt "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
)

func GetSecretValue(cfg aws.Config, secretID string) (string, error) {
	svc := secretsmanager.NewFromConfig(cfg)
	output, err := svc.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretID),
	})
	if err != nil {
		return "", fmt.Errorf("failed to read secret value for secret: %s -> %v", secretID, err)
	}
	return aws.ToString(output.SecretString), nil
}

func FilterSecretsByPrefix(cfg aws.Config, prefix string) ([]smt.SecretListEntry, error) {
	var secrets []smt.SecretListEntry

	svc := secretsmanager.NewFromConfig(cfg)
	input := &secretsmanager.ListSecretsInput{}
	p := secretsmanager.NewListSecretsPaginator(svc, input)

	for p.HasMorePages() {
		output, err := p.NextPage(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("failed to list secrets: %v", err)
		}
		for _, s := range output.SecretList {
			if strings.HasPrefix(aws.ToString(s.Name), prefix) {
				secrets = append(secrets, s)
			}
		}
	}

	return secrets, nil
}

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
