package connectwise

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"testing"
)

func TestClient_Get(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		t.Fatalf("creating aws config: %v", err)
	}

	c, err := NewClientFromAWS(ctx, ssm.NewFromConfig(cfg), "/connectwise/creds/ticketbot", true)
	if err != nil {
		t.Fatalf("creating connectwise client from aws: %v", err)
	}
	t.Logf("Got CW creds from AWS - client ID: %s", c.creds.ClientId)
	t.Logf("resty client: %v", c.restClient)

	b, err := c.GetBoard(1, nil)
	if err != nil {
		t.Fatalf("getting board: %v", err)
	}

	t.Logf("Board name: %s", b.Name)
}
