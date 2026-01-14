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

	AlertsResp struct {
		Items    []Alert `json:"items"`
		Metadata `json:"metadata"`
	}

	Device struct {
		Facts          map[string]Fact `json:"facts"`
		OrgID          string          `json:"orgid"`
		AgentID        string          `json:"agentid"`
		AuditData      time.Time       `json:"audit_date"`
		AgentAuditDate time.Time       `json:"agent_audit_date"`
	}

	AlertSearchParams struct {
		Query   AlertsQuery `json:"query"`
		Page    int         `json:"page"`
		PerPage int         `json:"per_page"`
	}

	AlertsQuery struct {
		AlertIDs          []string `json:"ids,omitempty"`
		AgentIDs          []string `json:"agent_ids,omitempty"`
		Category          string   `json:"category,omitempty"`
		RemediationStatus string   `json:"remediation_status,omitempty"`
		Statuses          []string `json:"statuses,omitempty"`
		SortDirection     string   `json:"sort_direction,omitempty"`
		SortField         string   `json:"sort_field,omitempty"`
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

	Alert struct {
		ID                    string    `json:"id"`
		OrgID                 string    `json:"orgid"`
		Name                  string    `json:"name"`
		Status                string    `json:"status"`
		AgentID               string    `json:"agent_id"`
		FactName              string    `json:"fact_name"`
		FactIdentifier        string    `json:"fact_identifier"`
		Value                 any       `json:"value"`
		ValueType             string    `json:"value_type"`
		Selector              string    `json:"selector"`
		Level                 string    `json:"level"`
		Category              string    `json:"category"`
		Emails                []string  `json:"emails"`
		RemediationEnabled    bool      `json:"remediation_enabled"`
		RemediationTime       int       `json:"remediation_time"`
		RemediationStatus     string    `json:"remediation_status"`
		CreatedDate           time.Time `json:"created_date"`
		ResolvedUserEmail     string    `json:"resolved_user_email"`
		ResolvedDate          time.Time `json:"resolved_date"`
		SecondsToResolved     int       `json:"seconds_to_resolved"`
		AcknowledgedUserEmail string    `json:"ack_user_email"`
		AcknowledgedDate      time.Time `json:"ack_date"`
		SecondsToAcknowledge  int       `json:"seconds_to_ack"`
		Muted                 bool      `json:"muted"`
	}
)
