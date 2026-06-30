package middleware

import (
	"github.com/azmiagr/sakutera-softdev/internal/service"
	"github.com/azmiagr/sakutera-softdev/pkg/jwt"
)

type Interface interface {
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
