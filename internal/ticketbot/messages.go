package ticketbot

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"tctg-automation/pkg/webex"
)

func (s *server) processMessageSent(c *gin.Context) {
	w := &webex.MessageWebhookBody{}
	if err := c.ShouldBindJSON(w); err != nil {
		c.Error(fmt.Errorf("invalid request body: %w", err))
		return
	}

	msgID := w.Data.Id
	m, err := s.webexClient.GetMessage(c.Request.Context(), msgID)
	if err != nil {
		c.Error(fmt.Errorf("getting message details: %w", err))
		return
	}

	email := m.PersonEmail
	if email == "" {
		c.Error(errors.New("sender email is blank"))
		return
	}

	if email == s.webexBotEmail {
		// ignore messages sent by the bot itself, otherwise you get an infinite loop
		c.Status(http.StatusNoContent)
		return
	}

	// sending message back for testing
	p := webex.MessagePostBody{
		Person: email,
		Text:   fmt.Sprintf("your message was: %s", m.Text),
	}
	if err := s.webexClient.SendMessage(c.Request.Context(), p); err != nil {
		c.Error(fmt.Errorf("sending message: %w", err))
		return
	}

	c.Status(http.StatusNoContent)
	return
}
