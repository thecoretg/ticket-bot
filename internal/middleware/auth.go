package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thecoretg/ticketbot/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func APIKeyAuth(r models.APIKeyRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid authorization header"})
			return
		}

		key := strings.TrimPrefix(auth, "Bearer ")

		keys, err := r.List(c.Request.Context())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}

		var userID int32
		found := false
		for _, k := range keys {
			if bcrypt.CompareHashAndPassword(k.KeyHash, []byte(key)) == nil {
				userID = int32(k.UserID)
				found = true
				break
			}
		}

		if !found {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid api key"})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
