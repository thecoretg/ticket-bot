package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (cl *Client) handleBaseAdminPage(c *gin.Context) {
	boards, err := cl.Queries.ListBoards(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.HTML(http.StatusOK, "base.gohtml", gin.H{"Title": "TicketBot Admin Dashboard", "Boards": boards})
}

func (cl *Client) handleAdminUsersPage(c *gin.Context) {
	users, err := cl.Queries.ListUsers(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.HTML(http.StatusOK, "users.gohtml", gin.H{"Title": "Users", "Users": users})
}
