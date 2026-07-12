package rest

import (
	"net/http"
	"strconv"

	"github.com/azmiagr/sakutera-softdev/model"
	apperr "github.com/azmiagr/sakutera-softdev/pkg/errors"
	"github.com/azmiagr/sakutera-softdev/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *Rest) GetTransactionSources(c *gin.Context) {
	provider := c.Query("provider")
	resp, err := r.service.TransactionService.GetSources(provider)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "sumber penghasilan berhasil diambil", resp)
}

func (r *Rest) CreateTransaction(c *gin.Context) {
	user, err := r.service.JwtAuth.GetLoginUser(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "user tidak ditemukan", err)
		return
	}

	var req model.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "request tidak valid", err)
		return
	}

	resp, err := r.service.TransactionService.CreateTransaction(user.UserID, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, resp.Message, resp)
}

func (r *Rest) UploadTransactionAttachment(c *gin.Context) {
	_, err := r.service.JwtAuth.GetLoginUser(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "user tidak ditemukan", err)
		return
	}

	file, err := c.FormFile("photo")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "foto tidak ditemukan", err)
		return
	}

	var req model.UploadAttachmentRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "request tidak valid", err)
		return
	}

	resp, err := r.service.TransactionService.UploadAttachment(file, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "foto berhasil diupload", resp)
}

func (r *Rest) GetTransactions(c *gin.Context) {
	user, err := r.service.JwtAuth.GetLoginUser(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "user tidak ditemukan", err)
		return
	}

	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	resp, err := r.service.TransactionService.GetTransactions(user.UserID, limit)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "transaksi berhasil diambil", resp)
}

func (r *Rest) GetLedger(c *gin.Context) {
	user, err := r.service.JwtAuth.GetLoginUser(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "user tidak ditemukan", err)
		return
	}

	period := c.Query("period")
	if period == "" {
		period = "all"
	}

	var sourceID *uuid.UUID
	if raw := c.Query("source_id"); raw != "" {
		parsed, err := uuid.Parse(raw)
		if err != nil {
			response.HandleError(c, apperr.BadRequest("source_id tidak valid"))
			return
		}
		sourceID = &parsed
	}

	resp, err := r.service.TransactionService.GetLedger(user.UserID, period, sourceID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "buku kas berhasil diambil", resp)
}
