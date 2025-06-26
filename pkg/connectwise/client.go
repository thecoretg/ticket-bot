package connectwise

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"net/http"
	"tctg-automation/pkg/amz"
	"time"
)

type Client struct {
	httpClient   *http.Client
	encodedCreds string
	clientId     string
	retryConfig  *RetryConfig
}

type Creds struct {
	PublicKey  string
	PrivateKey string
	ClientId   string
	CompanyId  string // The company name you enter when you log in to the PSA
}

type RetryConfig struct {
	MaxRetries        int
	InitialDelay      time.Duration
	MaxDelay          time.Duration
	BackOffMultiplier float64
}

func NewClient(creds Creds, httpClient *http.Client, retryConfig *RetryConfig) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if retryConfig == nil {
		retryConfig = DefaultRetryConfig()
	}

	username := fmt.Sprintf("%s+%s", creds.CompanyId, creds.PublicKey)
	return &Client{
		httpClient:   httpClient,
		encodedCreds: basicAuth(username, creds.PrivateKey),
		clientId:     creds.ClientId,
		retryConfig:  retryConfig,
	}
}

func NewClientFromAWS(ctx context.Context, httpClient *http.Client, retryConfig *RetryConfig, s *ssm.Client, paramName string, withDecryption bool) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if retryConfig == nil {
		retryConfig = DefaultRetryConfig()
	}

	creds, err := GetCredsFromAWS(ctx, s, paramName, withDecryption)
	if err != nil {
		return nil, fmt.Errorf("getting creds from AWS: %w", err)
	}

	username := fmt.Sprintf("%s+%s", creds.CompanyId, creds.PublicKey)
	return &Client{
		httpClient:   httpClient,
		encodedCreds: basicAuth(username, creds.PrivateKey),
		clientId:     creds.ClientId,
		retryConfig:  retryConfig,
	}, nil
}

func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:        5,
		InitialDelay:      500 * time.Millisecond,
		MaxDelay:          15 * time.Second,
		BackOffMultiplier: 1.5,
	}
}

func NewRetryConfig(maxRetries int, initialDelay, maxDelay time.Duration, backOffMultiplier float64) *RetryConfig {
	return &RetryConfig{
		MaxRetries:        maxRetries,
		InitialDelay:      initialDelay,
		MaxDelay:          maxDelay,
		BackOffMultiplier: backOffMultiplier,
	}
}

func GetCredsFromAWS(ctx context.Context, s *ssm.Client, paramName string, withDecryption bool) (*Creds, error) {
	c := &Creds{}
	if err := amz.GetAndUnmarshalParam(ctx, s, paramName, withDecryption, c); err != nil {
		return nil, fmt.Errorf("getting connectwise creds from AWS: %w", err)
	}

	return c, nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
