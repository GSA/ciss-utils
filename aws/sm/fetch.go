package sm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	smt "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
)

func GetSecret(cfg aws.Config, secretID string) (Secret, error) {
	var sec Secret
	s, err := GetSecretValue(cfg, secretID)
	if err != nil {
		return sec, err
	}
	err = json.Unmarshal([]byte(s), &sec)
	if err != nil {
		return sec, err
	}
	return sec, nil
}

func SecretExists(cfg aws.Config, secretID string) bool {
	svc := secretsmanager.NewFromConfig(cfg)
	_, err := svc.DescribeSecret(context.TODO(), &secretsmanager.DescribeSecretInput{
		SecretId: aws.String(secretID),
	})
	// Check for other exceptions
	return err == nil
}

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
