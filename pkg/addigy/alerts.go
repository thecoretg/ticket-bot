package addigy

import (
	"errors"
	"fmt"
)

type (
	AlertStatus       string
	RemediationStatus string
)

const (
	AlertStatusUnattended   AlertStatus = "Unattended"
	AlertStatusMuted        AlertStatus = "Muted" // is this right?
	AlertStatusAcknowledged AlertStatus = "Acknowledged"
	AlertStatusResolved     AlertStatus = "Resolved"
)

var ErrAlertNotFound = errors.New("alert not found")

func (c *Client) GetAlert(id string) (*Alert, error) {
	p := &AlertSearchParams{
		Query: AlertsQuery{
			AlertIDs: []string{id},
		},
	}
	a, err := c.searchAlerts(p, 1)
	if err != nil {
		return nil, fmt.Errorf("searching alerts: %w", err)
	}

	if len(a) == 0 {
		return nil, ErrAlertNotFound
	}

	return &a[0], nil
}

func (c *Client) searchAlerts(p *AlertSearchParams, pageLimit int) ([]Alert, error) {
	var a []Alert
	if p.Page == 0 {
		p.Page = c.defaultPage
	}

	if p.PerPage == 0 {
		p.PerPage = c.defaultPerPage
	}

	if p.SortField == "" {
		p.SortField = "created_date"
	}

	if p.SortDirection == "" {
		p.SortDirection = "desc"
	}

	for {
		resp, err := Post[AlertsResp](c, "oa/monitoring/alerts/query", p)
		if err != nil {
			return nil, fmt.Errorf("getting alerts page %d: %w", p.Page, err)
		}

		a = append(a, resp.Items...)
		if p.Page >= resp.PageCount {
			break
		}

		p.Page++

		if pageLimit != 0 && p.Page > pageLimit {
			break
		}
	}

	return a, nil
}
