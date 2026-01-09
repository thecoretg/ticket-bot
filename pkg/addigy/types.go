package addigy

import (
	"time"
)

type (
	Metadata struct {
		Page        int `json:"page"`
		PageCount   int `json:"page_count"`
		PerPage     int `json:"per_page"`
		ResultCount int `json:"result_count"`
		Total       int `json:"total"`
	}

	DevicesResp struct {
		Items    []Device `json:"items"`
		Metadata `json:"metadata"`
	}

	Device struct {
		Facts          map[string]Fact `json:"facts"`
		OrgID          string          `json:"orgid"`
		AgentID        string          `json:"agentid"`
		AuditData      time.Time       `json:"audit_date"`
		AgentAuditDate time.Time       `json:"agent_audit_date"`
	}

	DeviceSearchParams struct {
		DesiredFactIdentifiers []string    `json:"desired_fact_identifiers"`
		Query                  DeviceQuery `json:"query"`
		Page                   int         `json:"page"`
		PerPage                int         `json:"per_page"`
	}

	DeviceQuery struct {
		Filters []DeviceQueryFilter `json:"filters,omitempty"`
	}

	DeviceQueryFilter struct {
		AuditField string `json:"audit_field,omitempty"`
		Operation  string `json:"operation,omitempty"`
		RangeValue string `json:"range_value,omitempty"`
		Type       string `json:"type,omitempty"`
		Value      any    `json:"value,omitempty"`
	}

	Fact struct {
		Value    any    `json:"value,omitempty"`
		Type     string `json:"type,omitempty"`
		ErrorMsg string `json:"error_msg,omitempty"`
	}
)
