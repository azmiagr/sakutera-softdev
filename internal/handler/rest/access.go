package rest

import (
	"net/http"

	"github.com/azmiagr/sakutera-softdev/model"
	"github.com/azmiagr/sakutera-softdev/pkg/response"
	"github.com/gin-gonic/gin"
)

func (r *Rest) GetConsents(c *gin.Context) {
	user, err := r.service.JwtAuth.GetLoginUser(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "user tidak ditemukan", err)
		return
	}

	resp, err := r.service.AccessService.GetConsents(user.UserID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "daftar akses berhasil dimuat", resp)
}

func (r *Rest) GrantAccess(c *gin.Context) {
	user, err := r.service.JwtAuth.GetLoginUser(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "user tidak ditemukan", err)
		return
	}

	var req model.GrantAccessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "request tidak valid", err)
		return
	}

	if err := r.service.AccessService.GrantAccess(user.UserID, req); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "akses berhasil diberikan", nil)
}

func (r *Rest) RevokeAccess(c *gin.Context) {
	user, err := r.service.JwtAuth.GetLoginUser(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "user tidak ditemukan", err)
		return
	}

	consentID := c.Param("consent_id")
	if consentID == "" {
		response.Error(c, http.StatusBadRequest, "consent_id diperlukan", nil)
		return
	}

	if err := r.service.AccessService.RevokeAccess(user.UserID, consentID); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "akses berhasil dicabut", nil)
}

func (r *Rest) GetAccessLogs(c *gin.Context) {
	user, err := r.service.JwtAuth.GetLoginUser(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "user tidak ditemukan", err)
		return
	}

	filter := c.Query("filter")

	resp, err := r.service.AccessService.GetAccessLogs(user.UserID, filter)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "riwayat akses berhasil dimuat", resp)
}

func (r *Rest) GetOrganizations(c *gin.Context) {
	orgs, err := r.service.AccessService.GetOrganizations()
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "daftar organisasi berhasil dimuat", orgs)
}
