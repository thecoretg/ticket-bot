package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

var ErrAddigyAlertConfigNotFound = errors.New("addigy alert config not found")

type AddigyAlertConfig struct {
	ID                   int       `json:"id"`
	CWBoardID            int       `json:"cw_board_id"`
	UnattendedStatusID   int       `json:"unattended_status_id"`
	AcknowledgedStatusID int       `json:"acknowledged_status_id"`
	Mute1DayStatusID     int       `json:"mute_1_day_status_id"`
	Mute5DayStatusID     int       `json:"mute_5_day_status_id"`
	Mute10DayStatusID    int       `json:"mute_10_day_status_id"`
	Mute30DayStatusID    int       `json:"mute_30_day_status_id"`
	UpdatedOn            time.Time `json:"updated_on"`
	AddedOn              time.Time `json:"added_on"`
}

type AddigyAlertConfigRepository interface {
	WithTx(tx pgx.Tx) AddigyAlertConfigRepository
	Get(ctx context.Context) (*AddigyAlertConfig, error)
	Upsert(ctx context.Context, c *AddigyAlertConfig) (*AddigyAlertConfig, error)
	Delete(ctx context.Context) error
}

var ErrAddigyAlertNotFound = errors.New("addigy alert not found")

type AddigyAlert struct {
	ID              string     `json:"id"`
	TicketID        *int       `json:"ticket_id"`
	Level           string     `json:"level"`
	Category        string     `json:"category"`
	Name            string     `json:"name"`
	FactName        string     `json:"fact_name"`
	FactIdentifier  string     `json:"fact_identifier"`
	FactType        string     `json:"fact_type"`
	Selector        string     `json:"selector"`
	Status          string     `json:"status"`
	Value           *string    `json:"value"`
	Muted           bool       `json:"muted"`
	Remediation     bool       `json:"remediation"`
	ResolvedByEmail *string    `json:"resolved_by_email"`
	ResolvedOn      *time.Time `json:"resolved_on"`
	AcknowledgedOn  *time.Time `json:"acknowledged_on"`
	AddedOn         time.Time  `json:"added_on"`
}

type AddigyAlertRepository interface {
	WithTx(tx pgx.Tx) AddigyAlertRepository
	List(ctx context.Context) ([]*AddigyAlert, error)
	ListByStatus(ctx context.Context, status string) ([]*AddigyAlert, error)
	ListUnresolved(ctx context.Context) ([]*AddigyAlert, error)
	ListByTicket(ctx context.Context, ticketID int) ([]*AddigyAlert, error)
	Get(ctx context.Context, id string) (*AddigyAlert, error)
	Create(ctx context.Context, a *AddigyAlert) (*AddigyAlert, error)
	Update(ctx context.Context, a *AddigyAlert) (*AddigyAlert, error)
	UpdateTicket(ctx context.Context, id string, ticketID *int) error
	UpdateStatus(ctx context.Context, id string, status string) error
	Acknowledge(ctx context.Context, id string, acknowledgedOn time.Time) error
	Resolve(ctx context.Context, id string, resolvedOn time.Time, resolvedByEmail string) error
	Delete(ctx context.Context, id string) error
}
