package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thecoretg/ticketbot/internal/external/psa"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			if errors.Is(err, psa.ErrNotFound) {
				c.Status(http.StatusNoContent)
				return
			}
			slog.Error("error occurred in request", "error", err)
			c.Status(http.StatusInternalServerError)
			c.Abort()
			c.Writer.Flush()
		}
	}
}
