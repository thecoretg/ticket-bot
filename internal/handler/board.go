package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thecoretg/ticketbot/internal/models"
)

type BoardHandler struct {
	Repo models.BoardRepository
}

func NewBoardHandler(r models.BoardRepository) *BoardHandler {
	return &BoardHandler{Repo: r}
}

func (h *BoardHandler) ListBoards(c *gin.Context) {
	b, err := h.Repo.List(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, b)
}

func (h *BoardHandler) GetBoard(c *gin.Context) {
	id, err := convertID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, badIntErrorOutput(c.Param("id")))
		return
	}

	b, err := h.Repo.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrBoardNotFound) {
			c.JSON(http.StatusNotFound, errorOutput(err))
			return
		}
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, b)
}
