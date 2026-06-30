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

func (r *Rest) CheckPhone(c *gin.Context) {
	var req model.CheckPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "request tidak valid", err)
		return
	}

	resp, err := r.service.AuthService.CheckPhoneNumber(req.PhoneNumber)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	msg := "nomor ditemukan"
	if !resp.HasPIN {
		msg = "nomor ditemukan, silakan buat PIN terlebih dahulu"
	}
	response.Success(c, http.StatusOK, msg, resp)
}

func (r *Rest) SetPIN(c *gin.Context) {
	sessionToken := c.GetHeader("X-Session-Token")
	if sessionToken == "" {
		response.Error(c, http.StatusUnauthorized, "X-Session-Token diperlukan", nil)
		return
	}

	var req model.SetPINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "request tidak valid", err)
		return
	}

	resp, err := r.service.AuthService.SetPIN(sessionToken, req.PIN)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, resp.Message, resp)
}

func (r *Rest) LoginWithPIN(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "request tidak valid", err)
		return
	}

	resp, err := r.service.AuthService.LoginWithPIN(req.PhoneNumber, req.PIN)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, resp.Message, resp)
}
