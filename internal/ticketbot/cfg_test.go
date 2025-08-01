package ticketbot

import "testing"

func TestInit(t *testing.T) {
	cfg, err := InitCfg()
	if err != nil {
		t.Fatalf("error initializing config: %v", err)
	}

	t.Logf("got config: %v", cfg)
}
