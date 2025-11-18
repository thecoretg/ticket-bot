package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thecoretg/ticketbot/internal/models"
)

type UserHandler struct {
	UserRepo models.APIUserRepository
	KeyRepo  models.APIKeyRepository
}

func NewUserHandler(u models.APIUserRepository, k models.APIKeyRepository) *UserHandler {
	return &UserHandler{
		UserRepo: u,
		KeyRepo:  k,
	}
}

func (h *UserHandler) CreateAPIKey(c *gin.Context) {
	p := &struct {
		Email string `json:"email"`
	}{}

	if err := c.ShouldBindJSON(p); err != nil {
		c.Error(fmt.Errorf("bad json payload: %w", err))
		return
	}

	u, err := h.UserRepo.GetByEmail(c.Request.Context(), p.Email)
	if err != nil {
		if errors.Is(err, models.ErrAPIUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.Error(fmt.Errorf("querying user by email: %w", err))
		return
	}

}
