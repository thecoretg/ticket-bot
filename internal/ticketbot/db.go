package ticketbot

import (
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"time"
)

type DBHandler struct {
	db *sqlx.DB
}

var tablesStmt = `
CREATE TABLE IF NOT EXISTS ticket (
  id INTEGER PRIMARY KEY,
  board INTEGER NOT NULL,
  summary VARCHAR(100) NOT NULL,
  company INTEGER NOT NULL,
  contact INTEGER,
  owner INTEGER,
  resources TEXT,
  latest_note INTEGER,
  closed BOOLEAN DEFAULT FALSE,
  created TIMESTAMP NOT NULL,
  updated TIMESTAMP NOT NULL,
  closed_on TIMESTAMP
);
`

type Ticket struct {
	ID         int        `db:"id"`
	Board      int        `db:"board"`
	Company    int        `db:"company"`
	Contact    *int       `db:"contact"`
	LatestNote *int       `db:"latest_note"`
	Owner      *int       `db:"owner"`
	Summary    string     `db:"summary"`
	Resources  *string    `db:"resources"`
	Created    time.Time  `db:"created"`
	Updated    time.Time  `db:"updated"`
	ClosedOn   *time.Time `db:"closed_on"`
	Closed     bool       `db:"closed"`
}

func newTicket(ticketID, boardID, companyID, contactID, latestNoteID, ownerID int, summary, resources string, createdOn, updatedOn, closedOn time.Time, closed bool) *Ticket {
	return &Ticket{
		ID:         ticketID,
		Board:      boardID,
		Company:    companyID,
		Contact:    intToPtr(contactID),
		LatestNote: intToPtr(latestNoteID),
		Owner:      intToPtr(ownerID),
		Summary:    summary,
		Resources:  strToPtr(resources),
		Created:    createdOn,
		Updated:    updatedOn,
		ClosedOn:   timeToPtr(closedOn),
		Closed:     closed,
	}
}

func initDB(connString string) (*DBHandler, error) {
	db, err := sqlx.Connect("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("connecting to db: %w", err)
	}

	db.MustExec(tablesStmt)

	return &DBHandler{
		db: db,
	}, nil
}

func upsertTicketSQL() string {
	return `INSERT INTO ticket (id, board, summary, company, contact, owner, resources, latest_note, closed, created, updated)
		VALUES (:id, :board, :summary, :company, :contact, :owner, :resources, :latest_note, :closed, :created, :updated)
		ON CONFLICT (id) DO UPDATE SET
			board = EXCLUDED.board,
			summary = EXCLUDED.summary,
			company = EXCLUDED.company,
			contact = EXCLUDED.contact,
			owner = EXCLUDED.owner,
			resources = EXCLUDED.resources,
			latest_note = EXCLUDED.latest_note,
			closed = EXCLUDED.closed,
			created = EXCLUDED.created,
			updated = EXCLUDED.updated`
}

func (h *DBHandler) getTicket(ticketID int) (*Ticket, error) {
	t := &Ticket{}
	err := h.db.Get(t, "SELECT * FROM ticket WHERE id = $1", ticketID)
	if err != nil {
		return nil, fmt.Errorf("getting ticket by id: %w", err)
	}

	if t.ID == 0 {
		return nil, nil
	}

	return t, nil
}

func (h *DBHandler) listTickets() ([]Ticket, error) {
	var tickets []Ticket
	if err := h.db.Select(&tickets, "SELECT * FROM ticket"); err != nil {
		return nil, fmt.Errorf("listing tickets: %w", err)
	}

	return tickets, nil
}

func (h *DBHandler) upsertTicket(t *Ticket) error {
	_, err := h.db.NamedExec(upsertTicketSQL(), t)
	if err != nil {
		return fmt.Errorf("inserting ticket: %w", err)
	}
	return nil
}

func (h *DBHandler) deleteTicket(ticketID int) error {
	_, err := h.db.Exec("DELETE FROM ticket WHERE id = $1", ticketID)
	if err != nil {
		return err
	}

	return nil
}
