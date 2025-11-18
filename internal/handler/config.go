package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thecoretg/ticketbot/internal/models"
	"github.com/thecoretg/ticketbot/internal/service/config"
)

type ConfigHandler struct {
	Service *config.Service
}

func NewConfigHandler(svc *config.Service) *ConfigHandler {
	return &ConfigHandler{Service: svc}
}

func (h *ConfigHandler) Get(c *gin.Context) {
	cfg, err := h.Service.Get(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, cfg)
}

func (h *ConfigHandler) Update(c *gin.Context) {
	p := &models.Config{}
	if err := c.ShouldBindJSON(p); err != nil {
		c.Error(fmt.Errorf("bad json payload: %w", err))
		return
	}

	cfg, err := h.Service.Update(c.Request.Context(), p)
	if err != nil {
		c.Error(fmt.Errorf("updating config: %w", err))
		return
	}

	c.JSON(http.StatusOK, cfg)
}
