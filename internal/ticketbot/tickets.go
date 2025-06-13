package ticketbot

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tctg-automation/pkg/connectwise"
	"tctg-automation/pkg/util"
)

func addTicketRoutes(r *gin.Engine, s *server) {
	g := r.Group("/tickets")

	g.POST("", s.processTicketPayload)
}

func (s *server) processTicketPayload(c *gin.Context) {
	w := &connectwise.WebhookPayload{}
	if err := c.ShouldBindJSON(w); err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorJSON("invalid request body"))
		return
	}

	status, msg := s.processTicket(c.Request.Context(), w)
	if msg == nil {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(status, msg)
}
