package sm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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
		var s Secret
		err = json.Unmarshal([]byte(aws.ToString(output.SecretString)), &s)
		if err != nil {
			return aws.Credentials{}, fmt.Errorf("failed to unmarshal secret value for secret: %s -> %v", secretID, err)
		}
		creds := strings.Split(s.Value, ":")
		return aws.Credentials{
			AccessKeyID:     creds[0],
			SecretAccessKey: creds[1],
		}, nil
	}
}
