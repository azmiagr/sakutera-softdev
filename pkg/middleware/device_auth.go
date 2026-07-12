package middleware

import (
	"net/http"
	"strings"

	"github.com/azmiagr/sakutera-softdev/pkg/response"
	"github.com/gin-gonic/gin"
)

func (m *middleware) AuthenticateDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(c, http.StatusUnauthorized, "device token diperlukan", nil)
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		device, err := m.service.CollectorService.AuthenticateDevice(token)
		if err != nil {
			response.HandleError(c, err)
			c.Abort()
			return
		}

		c.Set("device", device)
		c.Next()
	}
}
