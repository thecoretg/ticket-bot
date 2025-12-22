package handlers

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/thecoretg/ticketbot/internal/models"
	"github.com/thecoretg/ticketbot/internal/service/messenger"
	"github.com/thecoretg/ticketbot/internal/service/webexsvc"
	"github.com/thecoretg/ticketbot/pkg/webex"
)

type WebexHandler struct {
	WebexSvc     *webexsvc.Service
	MessengerSvc *messenger.Service
}

func NewWebexHandler(wx *webexsvc.Service, ms *messenger.Service) *WebexHandler {
	return &WebexHandler{
		WebexSvc:     wx,
		MessengerSvc: ms,
	}
}

func (h *WebexHandler) HandleMessageToBot(c *gin.Context) {
	w := &webex.MessageHookPayload{}
	if err := c.ShouldBindJSON(w); err != nil {
		badPayloadError(c, err)
		return
	}
	slog.Debug("bot messages: got message from webex", "id", w.Data.ID)

	ctx := context.WithoutCancel(c.Request.Context())
	go func() {
		msg, err := h.WebexSvc.GetMessage(ctx, w)
		if err != nil {
			if errors.Is(err, webexsvc.ErrMessageFromBot) {
				// messages the bot sends, sends a hook payload. No need to do anything
				// with these.
				return
			}
			slog.Error("bot messages: error fetching webex message details", "error", err)
			return
		}

		if err := h.MessengerSvc.ParseAndRespond(ctx, msg); err != nil {
			slog.Error("bot messages: error parsing/responding", "error", err)
			return
		}
	}()

	resultJSON(c, "received webex message")
}

func (h *WebexHandler) HandleAttachmentActions(c *gin.Context) {
	w := &webex.MessageHookPayload{}
	if err := c.ShouldBindJSON(w); err != nil {
		badPayloadError(c, err)
		return
	}
	slog.Debug("bot messages: got attachment action hook from webex", "id", w.Data.ID)
	ctx := context.WithoutCancel(c.Request.Context())
	go func() {
		ach, err := h.WebexSvc.GetAttachmentAction(ctx, w)
		if err != nil {
			if errors.Is(err, webexsvc.ErrMessageFromBot) {
				return
			}
			slog.Error("bot messages: error fetching attachment action", "error", err)
			return
		}

		if err := h.MessengerSvc.ParseAndRespondAttachment(ctx, ach); err != nil {
			slog.Error("bot messages: error parsing/responding", "error", err)
			return
		}
	}()

	resultJSON(c, "received webex attachment action")
}

func (h *WebexHandler) ListRecipients(c *gin.Context) {
	r, err := h.WebexSvc.ListRecipients(c.Request.Context())
	if err != nil {
		internalServerError(c, err)
		return
	}

	outputJSON(c, r)
}

func (h *WebexHandler) GetRoom(c *gin.Context) {
	id, err := convertID(c)
	if err != nil {
		badIntError(c)
		return
	}

	r, err := h.WebexSvc.GetRecipient(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrWebexRecipientNotFound) {
			notFoundError(c, err)
			return
		}
		internalServerError(c, err)
		return
	}

	outputJSON(c, r)
}
