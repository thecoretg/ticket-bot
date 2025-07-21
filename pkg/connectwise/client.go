package connectwise

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"resty.dev/v3"
	"tctg-automation/pkg/amz"
)

type Creds struct {
	PublicKey  string
	PrivateKey string
	ClientId   string
	CompanyId  string // The company name you enter when you log in to the PSA
}

type Client struct {
	restClient *resty.Client
	creds      *Creds
}

func NewClient(creds *Creds) *Client {
	c := resty.New()
	c.SetBasicAuth(fmt.Sprintf("%s+%s", creds.CompanyId, creds.PublicKey), creds.PrivateKey)
	c.SetHeader("Content-Type", "application/json")
	c.SetHeader("Accept", "application/json")
	c.SetHeader("clientId", creds.ClientId)
	c.SetRetryCount(3)

	return &Client{restClient: c, creds: creds}
}

func NewClientFromAWS(ctx context.Context, s *ssm.Client, paramName string, withDecryption bool) (*Client, error) {
	creds, err := GetCredsFromAWS(ctx, s, paramName, withDecryption)
	if err != nil {
		return nil, fmt.Errorf("getting creds from AWS: %w", err)
	}

	return NewClient(creds), nil
}

func GetCredsFromAWS(ctx context.Context, s *ssm.Client, paramName string, withDecryption bool) (*Creds, error) {
	c := &Creds{}
	if err := amz.GetAndUnmarshalParam(ctx, s, paramName, withDecryption, c); err != nil {
		return nil, fmt.Errorf("getting connectwise creds from AWS: %w", err)
	}

	return c, nil
}

func basicAuth(creds *Creds) string {
	username := fmt.Sprintf("%s+%s", creds.CompanyId, creds.PublicKey)
	auth := fmt.Sprintf("%s:%s", username, creds.PrivateKey)
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
