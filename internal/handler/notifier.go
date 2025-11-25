package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thecoretg/ticketbot/internal/models"
)

type NotifierHandler struct {
	BoardRepo    models.BoardRepository
	RoomRepo     models.WebexRoomRepository
	RulesRepo    models.NotifierRepository
	ForwardsRepo models.UserForwardRepository
}

func NewNotifierHandler(r models.NotifierRepository, br models.BoardRepository, wr models.WebexRoomRepository, fr models.UserForwardRepository) *NotifierHandler {
	return &NotifierHandler{
		BoardRepo:    br,
		RoomRepo:     wr,
		RulesRepo:    r,
		ForwardsRepo: fr,
	}
}

func (h *NotifierHandler) ListNotifierRules(c *gin.Context) {
	n, err := h.RulesRepo.ListAll(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, n)
}

func (h *NotifierHandler) GetNotifierRule(c *gin.Context) {
	id, err := convertID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, badIntErrorOutput(c.Param("id")))
		return
	}

	n, err := h.RulesRepo.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotifierNotFound) {
			c.JSON(http.StatusNotFound, errorOutput(err))
			return
		}
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, n)
}

func (h *NotifierHandler) AddNotifierRule(c *gin.Context) {
	ctx := c.Request.Context()
	p := &models.Notifier{}
	if err := c.ShouldBindJSON(p); err != nil {
		c.JSON(http.StatusBadRequest, errorOutput(fmt.Errorf("bad json payload: %w", err)))
		return
	}

	exists, err := h.RulesRepo.Exists(ctx, p.CwBoardID, p.WebexRoomID)
	if err != nil {
		c.Error(err)
		return
	}

	if exists {
		err = fmt.Errorf("notifier with board id %d and room id %d already exists", p.CwBoardID, p.WebexRoomID)
		c.JSON(http.StatusConflict, errorOutput(err))
	}

	if _, err = h.BoardRepo.Get(ctx, p.CwBoardID); err != nil {
		if errors.Is(err, models.ErrBoardNotFound) {
			c.JSON(http.StatusNotFound, err)
			return
		}
		c.Error(err)
		return
	}

	if _, err = h.RoomRepo.Get(ctx, p.WebexRoomID); err != nil {
		if errors.Is(err, models.ErrWebexRoomNotFound) {
			c.JSON(http.StatusNotFound, err)
			return
		}
		c.Error(err)
		return
	}

	n, err := h.RulesRepo.Insert(c.Request.Context(), p)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, n)
}

func (h *NotifierHandler) DeleteNotifierRule(c *gin.Context) {
	id, err := convertID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, badIntErrorOutput(c.Param("id")))
		return
	}

	if err := h.RulesRepo.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, models.ErrNotifierNotFound) {
			c.JSON(http.StatusNotFound, errorOutput(err))
			return
		}
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *NotifierHandler) ListForwards(c *gin.Context) {
	n, err := h.ForwardsRepo.ListAll(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, n)
}

func (h *NotifierHandler) GetForward(c *gin.Context) {
	id, err := convertID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, badIntErrorOutput(c.Param("id")))
		return
	}

	f, err := h.ForwardsRepo.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrUserForwardNotFound) {
			c.JSON(http.StatusNotFound, errorOutput(err))
			return
		}
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, f)
}

func (h *NotifierHandler) AddUserForward(c *gin.Context) {
	p := &models.UserForward{}
	if err := c.ShouldBindJSON(p); err != nil {
		c.JSON(http.StatusBadRequest, errorOutput(fmt.Errorf("bad json payload: %w", err)))
		return
	}

	f, err := h.ForwardsRepo.Insert(c.Request.Context(), *p)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, f)
}

func (h *NotifierHandler) DeleteUserForward(c *gin.Context) {
	id, err := convertID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, badIntErrorOutput(c.Param("id")))
		return
	}

	if err := h.ForwardsRepo.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, models.ErrUserForwardNotFound) {
			c.JSON(http.StatusNotFound, errorOutput(err))
			return
		}
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}
