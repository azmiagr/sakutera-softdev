package rest

import (
	"net/http"

	"github.com/azmiagr/sakutera-softdev/model"
	"github.com/azmiagr/sakutera-softdev/pkg/response"
	"github.com/gin-gonic/gin"
)

func (r *Rest) Register(c *gin.Context) {
	var req model.RegisterRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "request tidak valid", err)
		return
	}

	resp, err := r.service.AuthService.Register(req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, resp.Message, resp)
}

func (r *Rest) VerifyOTP(c *gin.Context) {
	sessionToken := c.GetHeader("X-Session-Token")
	if sessionToken == "" {
		response.Error(c, http.StatusUnauthorized, "X-Session-Token diperlukan", nil)
		return
	}

	var req model.VerifyOTPRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "request tidak valid", err)
		return
	}

	resp, err := r.service.AuthService.VerifyOTP(sessionToken, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, resp.Message, resp)
}
