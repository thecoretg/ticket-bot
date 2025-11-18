package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thecoretg/ticketbot/internal/external/psa"
	"github.com/thecoretg/ticketbot/internal/service/ticket"
)

type TicketHandler struct {
	Service *ticket.Service
}

func NewTicketHandler(svc *ticket.Service) *TicketHandler {
	return &TicketHandler{Service: svc}
}

func (h *TicketHandler) ProcessTicket(c *gin.Context) {
	w := &psa.WebhookPayload{}
	if err := c.ShouldBindJSON(w); err != nil {
		c.Error(fmt.Errorf("bad json payload: %w", err))
		return
	}

	switch w.Action {
	case "added", "updated":
		t, err := h.Service.ProcessTicket(c.Request.Context(), w.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, t)
	case "deleted":
		if err := h.Service.DeleteTicket(c.Request.Context(), w.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusOK)
		return
	default:
		c.Error(fmt.Errorf("invalid action; expected 'added', 'updated', or 'deleted'; got '%s'", w.Action))
		return
	}
}
