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
		Query         AlertsQuery `json:"query"`
		Page          int         `json:"page,omitempty"`
		PerPage       int         `json:"per_page,omitempty"`
		SortDirection string      `json:"sort_direction,omitempty"`
		SortField     string      `json:"sort_field,omitempty"`
	}

	AlertsQuery struct {
		AlertIDs          []string `json:"ids,omitempty"`
		AgentIDs          []string `json:"agent_ids,omitempty"`
		Category          string   `json:"category,omitempty"`
		RemediationStatus string   `json:"remediation_status,omitempty"`
		Statuses          []string `json:"statuses,omitempty"`
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
		// TODO: Find a solution to empty time strings; Addigy's API returns "" for
		// time such as ResolvedDate instead of null, and therefore it can't unmarshal
		// into time.Time. Probably need to do something like RespAlertToAlert.

		ID                    string   `json:"id,omitempty"`
		OrgID                 string   `json:"orgid,omitempty"`
		Name                  string   `json:"name,omitempty"`
		Status                string   `json:"status,omitempty"`
		AgentID               string   `json:"agent_id,omitempty"`
		FactName              string   `json:"fact_name,omitempty"`
		FactIdentifier        string   `json:"fact_identifier,omitempty"`
		Value                 any      `json:"value,omitempty"`
		ValueType             string   `json:"value_type,omitempty"`
		Selector              string   `json:"selector,omitempty"`
		Level                 string   `json:"level,omitempty"`
		Category              string   `json:"category,omitempty"`
		Emails                []string `json:"emails,omitempty"`
		RemediationEnabled    bool     `json:"remediation_enabled,omitempty"`
		RemediationTime       int      `json:"remediation_time,omitempty"`
		RemediationStatus     string   `json:"remediation_status,omitempty"`
		CreatedDate           string   `json:"created_date,omitempty"`
		ResolvedUserEmail     string   `json:"resolved_user_email,omitempty"`
		ResolvedDate          string   `json:"resolved_date,omitempty"`
		SecondsToResolved     int      `json:"seconds_to_resolved,omitempty"`
		AcknowledgedUserEmail string   `json:"ack_user_email,omitempty"`
		AcknowledgedDate      string   `json:"ack_date,omitempty"`
		SecondsToAcknowledge  int      `json:"seconds_to_ack,omitempty"`
		Muted                 bool     `json:"muted,omitempty"`
		TicketID              *int     `json:"ticket_id,omitempty"`
	}
)
