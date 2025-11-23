package middleware

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AuthMiddleware(sessionManager *scs.SessionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := sessionManager.GetInt64(c.Request.Context(), "userID")
		if userID == 0 { // Если userID не найден или равен 0 (значение по умолчанию для int64)
			zap.S().Warnw("Unauthorized access attempt", "path", c.Request.URL.Path)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.Set("userID", userID) // Устанавливаем userID в контекст Gin
		c.Next()
	}
}
