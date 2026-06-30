package rest

import (
	"net/http"

	"github.com/azmiagr/sakutera-softdev/model"
	"github.com/azmiagr/sakutera-softdev/pkg/response"
	"github.com/gin-gonic/gin"
)

func (r *Rest) GetWorkCategories(c *gin.Context) {
	resp, err := r.service.OnboardingService.GetWorkCategories()
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "data pekerjaan berhasil diambil", resp)
}

func (r *Rest) SelectWorkPlatform(c *gin.Context) {
	user, err := r.service.JwtAuth.GetLoginUser(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "user tidak ditemukan", err)
		return
	}

	var req model.SelectWorkPlatformRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "request tidak valid", err)
		return
	}

	resp, err := r.service.OnboardingService.SelectWorkPlatform(user.UserID, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, resp.Message, resp)
}

func (r *Rest) GetIncomeSources(c *gin.Context) {
	resp, err := r.service.OnboardingService.GetIncomeSources()
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "sumber penghasilan berhasil diambil", resp)
}

func (r *Rest) SelectIncomeSource(c *gin.Context) {
	user, err := r.service.JwtAuth.GetLoginUser(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "user tidak ditemukan", err)
		return
	}

	var req model.SelectIncomeSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "request tidak valid", err)
		return
	}

	resp, err := r.service.OnboardingService.SelectIncomeSource(user.UserID, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, resp.Message, resp)
}
