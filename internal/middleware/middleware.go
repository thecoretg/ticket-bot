package middleware

import "github.com/gin-gonic/gin"

func NoOp() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
