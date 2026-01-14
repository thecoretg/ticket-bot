package addigy

import (
	"errors"
	"os"
	"testing"
)

func TestDeviceSearch(t *testing.T) {
	cl, err := testClient(t)
	if err != nil {
		t.Fatalf("creating test client: %v", err)
	}

	p := &DeviceSearchParams{
		DesiredFactIdentifiers: []string{"serial_number"},
		PerPage:                10,
		Query: DeviceQuery{
			Filters: []DeviceQueryFilter{
				{
					AuditField: "processor_type",
					Operation:  "contains",
					Value:      []string{"M3"},
				},
			},
		},
	}

	devices, err := cl.SearchDevices(p, 2)
	if err != nil {
		t.Fatalf("getting devices: %v", err)
	}

	for _, d := range devices {
		f, ok := d.Facts["serial_number"]
		if !ok {
			t.Fatal("Got no serial number")
		}
		t.Logf("Got device: %s", f.Value)
	}
}

func TestGetAlert(t *testing.T) {
	cl, err := testClient(t)
	if err != nil {
		t.Fatalf("creating test client: %v", err)
	}

	id := "69680c02c697b82ded996037"
	a, err := cl.GetAlert(id)
	if err != nil {
		t.Fatalf("getting alert: %v", err)
	}

	t.Logf("Got alert; name:%s, fact:%s", a.Name, a.FactName)
}

func testClient(t *testing.T) (*Client, error) {
	t.Helper()
	token := os.Getenv("ADDIGY_TOKEN")
	if token == "" {
		return nil, errors.New("no addigy token provided")
	}

	p := ClientParams{
		Token:          token,
		DefaultPage:    1,
		DefaultPerPage: 50,
	}

	return NewClient(p), nil
}
