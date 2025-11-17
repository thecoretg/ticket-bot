package newserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thecoretg/ticketbot/internal/external/psa"
	"github.com/thecoretg/ticketbot/internal/middleware"
	"github.com/thecoretg/ticketbot/internal/service/ticket"
)

func registerTicketHookRoute(r *gin.Engine, svc *ticket.Service) {
	h := func(c *gin.Context) {
		w := &psa.WebhookPayload{}
		if err := c.ShouldBindJSON(w); err != nil {
			c.Error(fmt.Errorf("bad json payload: %w", err))
			return
		}

		switch w.Action {
		case "added", "updated":
			t, err := svc.ProcessTicket(c.Request.Context(), w.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, t)
		case "deleted":
			if err := svc.DeleteTicket(c.Request.Context(), w.ID); err != nil {
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

	r.POST("hooks/cw/tickets", h, middleware.ErrorHandler(), middleware.RequireConnectwiseSignature())
}
