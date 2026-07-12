package middleware

import (
	"github.com/azmiagr/sakutera-softdev/internal/service"
	"github.com/azmiagr/sakutera-softdev/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type Interface interface {
	Cors() gin.HandlerFunc
	AuthenticateUser() gin.HandlerFunc
	AuthenticateDevice() gin.HandlerFunc
}

type middleware struct {
	service *service.Service
	jwtAuth jwt.Interface
}

func Init(service *service.Service, jwtAuth jwt.Interface) Interface {
	return &middleware{
		service: service,
		jwtAuth: jwtAuth,
	}
}
