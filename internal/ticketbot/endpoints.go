package ticketbot

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"tctg-automation/pkg/connectwise"
	"tctg-automation/pkg/util"
	"tctg-automation/pkg/webex"
)

const (
	maxNoteLength = 300
)

func (s *Server) listBoardsEndpoint(c *gin.Context) {
	b := s.Boards
	if len(b) == 0 {
		b = []boardSetting{}
	}

	c.JSON(http.StatusOK, gin.H{"boards": b})
}

func (s *Server) addOrUpdateBoardEndpoint(c *gin.Context) {
	b := &boardSetting{}
	if err := c.ShouldBindJSON(b); err != nil {
		util.ErrorJSON(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := s.addBoardSetting(b); err != nil {
		util.ErrorJSON(c, http.StatusInternalServerError, "failed to add or update board setting")
		return
	}

	if err := s.refreshBoards(); err != nil {
		util.ErrorJSON(c, http.StatusInternalServerError, "failed to refresh boards")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Board setting added or updated successfully"})
}

func (s *Server) deleteBoardEndpoint(c *gin.Context) {
	boardIDStr := c.Param("board_id")
	boardID, err := strconv.Atoi(boardIDStr)
	if err != nil {
		util.ErrorJSON(c, http.StatusBadRequest, "board_id must be a valid integer")
		return
	}

	if err := s.deleteBoardSetting(boardID); err != nil {
		util.ErrorJSON(c, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := s.refreshBoards(); err != nil {
		util.ErrorJSON(c, http.StatusInternalServerError, "failed to refresh boards")
		return
	}

	c.Status(http.StatusNoContent)
}

func (s *Server) handleTicketEndpoint(c *gin.Context) {
	// Parse webhook payload
	w := &connectwise.WebhookPayload{}
	if err := c.ShouldBindJSON(w); err != nil {
		util.ErrorJSON(c, http.StatusInternalServerError, "invalid request body")
		return
	}

	// Validate action and ID
	if !validAction(w.Action) {
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("action '%s' is not of 'added' or 'updated'", w.Action)})
		return
	}
	if w.ID == 0 {
		util.ErrorJSON(c, http.StatusBadRequest, "ticket ID is required")
		return
	}

	// Fetch ticket and validate board
	ticket, err := s.cwClient.GetTicket(c.Request.Context(), w.ID, nil)
	if err != nil {
		util.ErrorJSON(c, http.StatusInternalServerError, "failed to get ticket")
		return
	}

	bs := s.ticketInEnabledBoard(ticket)
	if bs == nil {
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("ticket %d received, but board %d is not an enabled board for notifications", ticket.ID, ticket.Board.ID)})
		return
	}

	switch w.Action {
	case "added":
		s.handleNewTicket(c, ticket, bs)
	}
}

func (s *Server) handleNewTicket(c *gin.Context, ticket *connectwise.Ticket, bs *boardSetting) {
	if bs.WebexRoomID == "" {
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("board %d is enabled, but has no specified webex room ID for new tickets", ticket.Board.ID)})
		return
	}

	p := &connectwise.QueryParams{OrderBy: "_info/dateEntered desc"}
	notes, err := s.cwClient.ListServiceTicketNotes(c.Request.Context(), ticket.ID, p)
	if err != nil {
		util.ErrorJSON(c, http.StatusInternalServerError, "error receiving ticket notes")
		return
	}
	log.Printf("webex room id: %s\n", bs.WebexRoomID)

	m := buildNewTicketMessage(ticket, notes)
	w := webex.NewMessageToRoom(bs.WebexRoomID, m)
	if err := s.webexClient.SendMessage(c.Request.Context(), w); err != nil {
		slog.Error("sending new ticket message", "boardName", bs.BoardName, "webexRoomId", bs.WebexRoomID, "ticketId", ticket.ID, "ticketSummary", ticket.Summary, "error", err)
		util.ErrorJSON(c, http.StatusInternalServerError, fmt.Sprintf("sending message to webex room: %v", err))
		return
	}
	slog.Info("successfully sent new ticket message", "boardName", bs.BoardName, "webexRoomId", bs.WebexRoomID, "ticketId", ticket.ID, "ticketSummary", ticket.Summary)
	c.Status(http.StatusNoContent)
}

func (s *Server) ticketInEnabledBoard(ticket *connectwise.Ticket) *boardSetting {
	for _, b := range s.Boards {
		if b.BoardID == ticket.Board.ID {
			return &b
		}
	}

	return nil
}

func validAction(action string) bool {
	return action == "added" || action == "updated"
}

func buildNewTicketMessage(t *connectwise.Ticket, n []connectwise.ServiceTicketNoteAll) string {
	m := fmt.Sprintf("**New:** %s %s", ticketLink(t), t.Summary)

	// add requester contact and company name, or just company if no contact (rare)
	r := fmt.Sprintf("**Requester:** %s (No Contact)", t.Company.Name)
	if t.Contact.Name != "" {
		r = fmt.Sprintf("**Requester:** %s (%s)", t.Contact.Name, t.Company.Name)
	}
	m += fmt.Sprintf("\n%s", r)

	// add most recent note if present
	mr := mostRecentNote(n)
	if mr != nil {
		// trim note text and add ... if it exceeds the maximum
		noteTxt := mr.Text
		if len(noteTxt) > maxNoteLength {
			noteTxt = noteTxt[:maxNoteLength] + "..."
		}

		m += fmt.Sprintf("\n**Latest Note:** %s\n"+
			"%s",
			noteSenderName(mr),
			addBlockQuotes(noteTxt),
		)
	}

	return m
}

func noteSenderName(n *connectwise.ServiceTicketNoteAll) string {
	if n.Member.Name != "" {
		return n.Member.Name
	} else if n.Contact.Name != "" {
		return n.Contact.Name
	} else {
		return "N/A"
	}
}

func mostRecentNote(n []connectwise.ServiceTicketNoteAll) *connectwise.ServiceTicketNoteAll {
	for _, note := range n {
		if note.Text != "" {
			return &note
		}
	}

	return nil
}

func ticketLink(t *connectwise.Ticket) string {
	return fmt.Sprintf("[%d](https://na.myconnectwise.net/v4_6_release/services/system_io/Service/fv_sr100_request.rails?service_recid=%d&companyName=securenetit)", t.ID, t.ID)
}

func addBlockQuotes(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if line == "" {
			lines[i] = ">"
		} else {
			lines[i] = "> " + line
		}
	}

	return strings.Join(lines, "\n")
}
