package newserver

import (
	"github.com/gin-gonic/gin"
	"github.com/thecoretg/ticketbot/internal/handler"
	"github.com/thecoretg/ticketbot/internal/middleware"
)

func (a *App) addRoutes(g *gin.Engine) {
	errh := middleware.ErrorHandler()
	auth := middleware.APIKeyAuth(a.Svc.User.Keys)
	cws := middleware.RequireConnectwiseSignature()

	u := g.Group("users", errh, auth)
	uh := handler.NewUserHandler(a.Svc.User)
	registerUserRoutes(u, uh)

	c := g.Group("config", errh, auth)
	ch := handler.NewConfigHandler(a.Svc.Config)
	registerConfigRoutes(c, ch)

	b := g.Group("boards", errh, auth)
	bh := handler.NewBoardHandler(a.Stores.CW.Board)
	registerBoardRoutes(b, bh)

	n := g.Group("notifiers", errh, auth)
	nh := handler.NewNotifierHandler(a.Stores.Notifiers, a.Stores.CW.Board, a.Stores.WebexRoom)
	registerNotifierRoutes(n, nh)

	th := handler.NewTicketHandler(a.Svc.Ticket)
	g.POST("hooks/cw/tickets", th.ProcessTicket, errh, cws)
}

func registerUserRoutes(r *gin.RouterGroup, h *handler.UserHandler) {
	r.GET("", h.ListUsers)
	r.GET(":id", h.GetUser)
	r.DELETE(":id")

	k := r.Group("keys")
	k.GET("", h.ListAPIKeys)
	k.GET(":id", h.GetAPIKey)
	k.POST("", h.AddAPIKey)
	k.DELETE(":id", h.DeleteAPIKey)
}

func registerConfigRoutes(r *gin.RouterGroup, h *handler.ConfigHandler) {
	r.GET("", h.Get)
	r.PUT("", h.Update)
}

func registerBoardRoutes(r *gin.RouterGroup, h *handler.BoardHandler) {
	r.GET("", h.ListBoards)
	r.GET(":id", h.GetBoard)
}

func registerNotifierRoutes(r *gin.RouterGroup, h *handler.NotifierHandler) {
	r.GET("", h.ListNotifiers)
	r.POST("", h.AddNotifier)
}
