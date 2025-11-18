package newserver

import (
	"github.com/gin-gonic/gin"
	"github.com/thecoretg/ticketbot/internal/handler"
	"github.com/thecoretg/ticketbot/internal/middleware"
)

func (a *App) addRoutes(r *gin.Engine) {
	errh := middleware.ErrorHandler()
	auth := middleware.APIKeyAuth(a.Stores.APIKey)

	th := handler.NewTicketHandler(a.Svc.Ticket)
	r.POST("hooks/cw/tickets", th.ProcessTicket, errh, middleware.RequireConnectwiseSignature())

	ch := handler.NewConfigHandler(a.Svc.Config)
	cfg := r.Group("config", errh, auth)
	cfg.GET("", ch.Get)
	cfg.PUT("", ch.Update)
}
