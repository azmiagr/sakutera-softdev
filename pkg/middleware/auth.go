package middleware

import (
	"net/http"
	"strings"

	"github.com/azmiagr/sakutera-softdev/pkg/response"
	"github.com/gin-gonic/gin"
)

func (m *middleware) AuthenticateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(c, http.StatusUnauthorized, "token diperlukan", nil)
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := m.jwtAuth.ValidateToken(token)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "token tidak valid", err)
			c.Abort()
			return
		}

		isRevoked, err := m.service.AuthService.IsTokenRevoked(token)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "gagal memverifikasi token", err)
			c.Abort()
			return
		}
		if isRevoked {
			response.Error(c, http.StatusUnauthorized, "token sudah tidak berlaku, silakan login kembali", nil)
			c.Abort()
			return
		}

		user, err := m.service.AuthService.GetUserByID(userID)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "user tidak ditemukan", err)
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
