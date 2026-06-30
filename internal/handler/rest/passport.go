package rest

import (
	"net/http"

	"github.com/azmiagr/sakutera-softdev/model"
	"github.com/azmiagr/sakutera-softdev/pkg/response"
	"github.com/gin-gonic/gin"
)

func (r *Rest) GetPassport(c *gin.Context) {
	user, err := r.service.JwtAuth.GetLoginUser(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "user tidak ditemukan", err)
		return
	}

	resp, err := r.service.PassportService.GetPassport(user.UserID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "income passport berhasil dimuat", resp)
}

func (r *Rest) PreviewPassport(c *gin.Context) {
	user, err := r.service.JwtAuth.GetLoginUser(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "user tidak ditemukan", err)
		return
	}

	period := c.Query("period")
	if period == "" {
		response.Error(c, http.StatusBadRequest, "period diperlukan", nil)
		return
	}

	resp, err := r.service.PassportService.PreviewPassport(user.UserID, period)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "preview passport berhasil dimuat", resp)
}

func (r *Rest) IssuePassport(c *gin.Context) {
	user, err := r.service.JwtAuth.GetLoginUser(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "user tidak ditemukan", err)
		return
	}

	var req model.IssuePassportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "request tidak valid", err)
		return
	}

	resp, err := r.service.PassportService.IssuePassport(user.UserID, req.Period)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "income passport berhasil diterbitkan", resp)
}
