package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/thecoretg/ticketbot/internal/db"
	"github.com/thecoretg/ticketbot/internal/psa"
	"github.com/thecoretg/ticketbot/internal/webex"
)

type connectwiseData struct {
	ticket *psa.Ticket
	note   *psa.ServiceTicketNote
}

type storedData struct {
	ticket       db.CwTicket
	company      db.CwCompany
	contact      db.CwContact
	owner        db.CwMember
	note         db.CwTicketNote
	board        db.CwBoard
	enabledRooms []db.WebexRoom
}

type requestState struct {
	logger *slog.Logger
	cwData *connectwiseData
	dbData *storedData

	action         string
	attemptNotify  bool
	syncing        bool
	webexMock      bool
	messagesToSend []webex.Message
	notified       bool
	noNotiReason   string
	roomsNotify    []string
	peopleNotify   []string
}

func (cl *Client) handleTickets(c *gin.Context) {
	w := &psa.WebhookPayload{}
	if err := c.ShouldBindJSON(w); err != nil {
		c.Error(fmt.Errorf("unmarshaling connectwise webhook payload: %w", err))
		return
	}

	if w.ID == 0 {
		c.Error(errors.New("received no ticket ID or zero"))
		return
	}

	switch w.Action {
	case "added", "updated":
		if err := cl.processTicket(c.Request.Context(), w.ID, w.Action, false); err != nil {
			c.Error(fmt.Errorf("ticket %d: adding or updating the ticket into data storage: %w", w.ID, err))
			return
		}

		c.Status(http.StatusOK)

	case "deleted":
		if err := cl.softDeleteTicket(c.Request.Context(), w.ID); err != nil {
			c.Error(fmt.Errorf("ticket %d: deleting ticket and its notes: %w", w.ID, err))
			return
		}

		c.Status(http.StatusOK)
	}
}

func (cl *Client) getTicketLock(ticketID int) *sync.Mutex {
	lockIface, _ := cl.ticketLocks.LoadOrStore(ticketID, &sync.Mutex{})
	return lockIface.(*sync.Mutex)
}

// processTicket serves as the primary handler for updating the data store with ticket data. It also will handle
// extra functionality such as ticket notifications.
func (cl *Client) processTicket(ctx context.Context, ticketID int, action string, syncing bool) error {
	// Lock the ticket so that extra calls don't interfere. Due to the nature of Connectwise updates will often
	// result in other hooks and actions taking place, which means a ticket rarely only sends one webhook payload.
	lock := cl.getTicketLock(ticketID)
	if !lock.TryLock() {
		lock.Lock()
	}

	defer func() {
		lock.Unlock()
	}()

	rs, err := cl.getInitialRequestState(ctx, action, ticketID, syncing)
	if err != nil {
		return fmt.Errorf("getting initial request state: %w", err)
	}

	defer func() {
		logTicketResult(rs)
	}()

	// Upsert the ticket data into the database
	rs, err = cl.upsertTicket(ctx, rs)
	if err != nil {
		return fmt.Errorf("upserting ticket data: %w", err)
	}

	// If a note exists and notifications are on, run the ticket notification action,
	// which checks if it meets message criteria and then notifies if valid.
	// AttemptNotify and the bypassNotis (used for preloads) acts as a hard block from even attempting.
	rs, err = cl.runNotificationAction(ctx, rs)
	if err != nil {
		return fmt.Errorf("running notifier: %w", err)
	}

	return nil
}

func (cl *Client) runNotificationAction(ctx context.Context, rs *requestState) (*requestState, error) {
	eligible, reason := rs.meetsMessageCriteria()
	if !eligible {
		rs.noNotiReason = reason
	}

	// set notified regardless, even if it doesn't meet critera - this is so it doesnt attempt again
	if err := cl.setNotified(ctx, rs.dbData.note.ID, true); err != nil {
		rs.noNotiReason = "SET_NOTIFIED_ERROR"
		return rs, fmt.Errorf("setting notified to true: %w", err)
	}

	if rs.syncing {
		rs.noNotiReason = "TICKET_SYNC"
		return rs, nil
	}

	if !eligible {
		return rs, nil
	}

	if !rs.attemptNotify {
		rs.noNotiReason = "ATTEMPT_NOTIFY_OFF"
		return rs, nil
	}

	rs, err := cl.makeAndSendMessages(ctx, rs)
	if err != nil {
		rs.noNotiReason = "MAKE_SEND_MESSAGES"
		return rs, fmt.Errorf("processing webex messages: %w", err)
	}
	rs.notified = true

	return rs, nil
}

func (cl *Client) softDeleteTicket(ctx context.Context, ticketID int) error {
	if err := cl.Queries.SoftDeleteTicket(ctx, ticketID); err != nil {
		return fmt.Errorf("soft deleting ticket: %w", err)
	}

	return nil
}

func (cl *Client) getInitialRequestState(ctx context.Context, action string, ticketID int, syncing bool) (*requestState, error) {
	cd, err := cl.getCwData(ticketID)
	if err != nil {
		return nil, fmt.Errorf("getting data from connectwise: %w", err)
	}

	sd, err := cl.ensureStoredData(ctx, cd)
	if err != nil {
		return nil, fmt.Errorf("ensuring stored data: %w", err)
	}

	return &requestState{
		logger:         slog.Default(),
		cwData:         cd,
		dbData:         sd,
		action:         action,
		webexMock:      cl.testing.mockWebex,
		messagesToSend: []webex.Message{},
		attemptNotify:  cl.Config.AttemptNotify,
		syncing:        syncing,
		roomsNotify:    []string{},
		peopleNotify:   []string{},
	}, nil
}

func (cl *Client) getCwData(ticketID int) (*connectwiseData, error) {
	ticket, err := cl.CWClient.GetTicket(ticketID, nil)
	if err != nil {
		return nil, fmt.Errorf("getting ticket: %w", err)
	}

	note, err := cl.CWClient.GetMostRecentTicketNote(ticketID)
	if err != nil {
		return nil, fmt.Errorf("getting most recent note: %w", err)
	}

	if note == nil {
		note = &psa.ServiceTicketNote{}
	}

	return &connectwiseData{
		ticket: ticket,
		note:   note,
	}, nil
}

func (cl *Client) ensureStoredData(ctx context.Context, cd *connectwiseData) (*storedData, error) {
	// first check for or create board since it needs to exist before the ticket
	board, err := cl.ensureBoardInStore(ctx, cd.ticket.Board.ID)
	if err != nil {
		return nil, fmt.Errorf("ensuring board in store: %w", err)
	}

	company, err := cl.ensureCompanyInStore(ctx, cd.ticket.Company.ID)
	if err != nil {
		return nil, fmt.Errorf("ensuring company in store: %w", err)
	}

	contact := db.CwContact{}
	if cd.ticket.Contact.ID != 0 {
		contact, err = cl.ensureContactInStore(ctx, cd.ticket.Contact.ID)
		if err != nil {
			return nil, fmt.Errorf("ensuring contact in store: %w", err)
		}
	}

	owner := db.CwMember{}
	if cd.ticket.Owner.ID != 0 {
		owner, err = cl.ensureMemberInStore(ctx, cd.ticket.Owner.ID)
		if err != nil {
			return nil, fmt.Errorf("ensuring owner in store: %w", err)
		}
	}

	// check for, or create ticket
	ticket, err := cl.ensureTicketInStore(ctx, cd)
	if err != nil {
		return nil, fmt.Errorf("ensuring ticket in store: %w", err)
	}

	// start with empty note, use existing or created note if there is a note in the ticket
	note := db.CwTicketNote{}
	if cd.note.ID != 0 {
		note, err = cl.ensureNoteInStore(ctx, cd)
		if err != nil {
			return nil, fmt.Errorf("ensuring note in store: %w", err)
		}
	}

	cons, err := cl.Queries.ListNotifierConnectionsByBoard(ctx, board.ID)
	if err != nil {
		return nil, fmt.Errorf("getting rooms to notify: %w", err)
	}

	return &storedData{
		ticket:       ticket,
		company:      company,
		contact:      contact,
		owner:        owner,
		note:         note,
		board:        board,
		enabledRooms: roomsFromNotifiers(cons),
	}, nil
}

func (cl *Client) ensureTicketInStore(ctx context.Context, cd *connectwiseData) (db.CwTicket, error) {
	ticket, err := cl.Queries.GetTicket(ctx, cd.ticket.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p := db.UpsertTicketParams{
				ID:        cd.ticket.ID,
				Summary:   cd.ticket.Summary,
				BoardID:   cd.ticket.Board.ID,
				OwnerID:   intToPtr(cd.ticket.Owner.ID),
				CompanyID: cd.ticket.Company.ID,
				ContactID: intToPtr(cd.ticket.Contact.ID),
				Resources: &cd.ticket.Resources,
				UpdatedBy: &cd.ticket.Info.UpdatedBy,
			}

			ticket, err = cl.Queries.UpsertTicket(ctx, p)
			if err != nil {
				return db.CwTicket{}, fmt.Errorf("inserting ticket into db: %w", err)
			}

			return ticket, nil
		} else {
			return db.CwTicket{}, fmt.Errorf("getting ticket from storage: %w", err)
		}
	}

	return ticket, nil
}

func (cl *Client) upsertTicket(ctx context.Context, rs *requestState) (*requestState, error) {
	var err error
	p := cwDataToUpdateTicketParams(rs.cwData)
	rs.dbData.ticket, err = cl.Queries.UpsertTicket(ctx, p)
	if err != nil {
		return rs, fmt.Errorf("updating ticket in store: %w", err)
	}

	return rs, nil
}

func logTicketResult(rs *requestState) {
	tg := slog.Group("ticket",
		slog.Int("id", rs.dbData.ticket.ID),
		slog.String("action", rs.action),
		slog.Bool("notified", rs.notified),
		notifyLogGroup(rs),
		boardLogGroup(rs.dbData.board),
		companyLogGroup(rs.dbData.company),
		contactLogGroup(rs.dbData.contact),
		ownerLogGroup(rs.dbData.owner),
		noteLogGroup(rs.dbData.note),
	)

	msg := "ticket processed"
	if rs.webexMock {
		msg = "ticket processed with webex mocking"
	}

	rs.logger.With(tg).Info(msg)
}

func boardLogGroup(board db.CwBoard) slog.Attr {
	return slog.Group("board",
		slog.Int("id", board.ID),
		slog.String("name", board.Name),
	)
}

func companyLogGroup(company db.CwCompany) slog.Attr {
	return slog.Group("company",
		slog.Int("id", company.ID),
		slog.String("name", company.Name),
	)
}

func contactLogGroup(contact db.CwContact) slog.Attr {
	if contact.ID == 0 {
		return slog.Bool("contact", false)
	}

	// TODO: get company name in here
	lastName := ""
	companyID := 0

	if contact.LastName != nil {
		lastName = *contact.LastName
	}

	if contact.CompanyID != nil {
		companyID = *contact.CompanyID
	}

	return slog.Group("contact",
		slog.Int("id", contact.ID),
		slog.Int("company_id", companyID),
		slog.String("first_name", contact.FirstName),
		slog.String("last_name", lastName),
	)
}

func ownerLogGroup(owner db.CwMember) slog.Attr {
	if owner.ID == 0 {
		return slog.Bool("owner", false)
	}

	return slog.Group("owner",
		slog.Int("id", owner.ID),
		slog.String("identifier", owner.Identifier),
		slog.String("first_name", owner.FirstName),
		slog.String("last_name", owner.LastName),
		slog.String("primary_email", owner.PrimaryEmail),
	)
}

func noteLogGroup(note db.CwTicketNote) slog.Attr {
	// TODO: member, contact name, already notified
	if note.ID == 0 {
		return slog.Bool("latest_note", false)
	}

	memberID := 0
	contactID := 0
	if note.MemberID != nil {
		memberID = *note.MemberID
	}

	if note.ContactID != nil {
		contactID = *note.ContactID
	}

	return slog.Group("latest_note",
		slog.Int("id", note.ID),
		slog.Int("member_id", memberID),
		slog.Int("contact_id", contactID),
	)
}

func notifyLogGroup(rs *requestState) slog.Attr {
	if !rs.attemptNotify {
		return slog.String("notifications", "disabled")
	}

	var er []string
	for _, r := range rs.dbData.enabledRooms {
		er = append(er, r.Name)
	}

	attrs := []slog.Attr{
		slog.Bool("sent", rs.notified),
		slog.String("enabled_rooms", strings.Join(er, ",")),
		slog.String("rooms", strings.Join(rs.roomsNotify, ",")),
		slog.String("people", strings.Join(rs.peopleNotify, ",")),
	}

	if rs.noNotiReason != "" {
		a := slog.String("no_noti_reason", rs.noNotiReason)
		attrs = append(attrs, a)
	}

	var anyAttrs []any
	for _, a := range attrs {
		anyAttrs = append(anyAttrs, a)
	}

	return slog.Group("webex_notifications", anyAttrs...)
}

// meetsMessageCriteria checks if a message would be allowed to send a notification,
// depending on if it was added or updated, if the note changed, and the board's notification settings.
func (rs *requestState) meetsMessageCriteria() (bool, string) {
	meetsCrit := false
	switch rs.action {
	case "added":
		meetsCrit = roomsToNotifyExist(rs.dbData)
		if !meetsCrit {
			return false, "NO_ROOMS_TO_NOTIFY"
		}

		return true, ""

	case "updated":
		roomsExist := roomsToNotifyExist(rs.dbData)
		meetsCrit = !rs.dbData.note.Notified && roomsToNotifyExist(rs.dbData)

		if !meetsCrit {
			if rs.dbData.note.Notified {
				return false, "NOTE_ALREADY_NOTIFIED"
			}

			if !roomsExist {
				return false, "NO_ROOMS_TO_NOTIFY"
			}
		}
		return true, ""

	default:
		return false, "INELIGIBLE_ACTION_TYPE"
	}
}

func roomsToNotifyExist(sd *storedData) bool {
	return sd.enabledRooms != nil && len(sd.enabledRooms) > 0
}

func cwDataToUpdateTicketParams(cd *connectwiseData) db.UpsertTicketParams {
	return db.UpsertTicketParams{
		ID:        cd.ticket.ID,
		Summary:   cd.ticket.Summary,
		BoardID:   cd.ticket.Board.ID,
		OwnerID:   intToPtr(cd.ticket.Owner.ID),
		CompanyID: cd.ticket.Company.ID,
		ContactID: intToPtr(cd.ticket.Contact.ID),
		Resources: strToPtr(cd.ticket.Resources),
		UpdatedBy: strToPtr(cd.ticket.Info.UpdatedBy),
	}
}

func roomsFromNotifiers(notifiers []db.ListNotifierConnectionsByBoardRow) []db.WebexRoom {
	var rooms []db.WebexRoom
	for _, n := range notifiers {
		rooms = append(rooms, n.WebexRoom)
	}

	return rooms
}
