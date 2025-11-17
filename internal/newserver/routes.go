package newserver

import (
	"github.com/gin-gonic/gin"
	"github.com/thecoretg/ticketbot/internal/middleware"
)

func (a *App) addRoutes(r *gin.Engine) {
	auth := middleware.NoOp()
	if !a.TestFlags.skipAuth {
		// TODO: DONT SKIP: add the real auth handler
	}

	registerTicketHookRoute(r, a.Svc.Ticket)
}
