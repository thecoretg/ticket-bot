package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thecoretg/ticketbot/internal/models"
	"github.com/thecoretg/ticketbot/internal/service/webexsvc"
)

type WebexHandler struct {
	Service *webexsvc.Service
}

func NewWebexHandler(svc *webexsvc.Service) *WebexHandler {
	return &WebexHandler{Service: svc}
}

func (h *WebexHandler) ListRooms(c *gin.Context) {
	r, err := h.Service.ListRooms(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, r)
}

func (h *WebexHandler) GetRoom(c *gin.Context) {
	id, err := convertID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, badIntErrorOutput(c.Param("id")))
		return
	}

	r, err := h.Service.GetRoom(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrWebexRoomNotFound) {
			c.JSON(http.StatusNotFound, errorOutput(err))
			return
		}
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, r)
}

func (h *WebexHandler) SyncRooms(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"result": "webex room sync started"})

	go func() {
		if err := h.Service.SyncRooms(context.Background()); err != nil {
			slog.Error("syncing webex rooms", "error", err)
		}
	}()
}
