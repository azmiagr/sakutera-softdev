package rest

import (
	"net/http"

	"github.com/azmiagr/sakutera-softdev/model"
	apperr "github.com/azmiagr/sakutera-softdev/pkg/errors"
	"github.com/azmiagr/sakutera-softdev/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *Rest) GetReviews(c *gin.Context) {
	user, err := r.service.JwtAuth.GetLoginUser(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "user tidak ditemukan", err)
		return
	}

	status := c.Query("status")

	resp, err := r.service.TransactionReviewService.GetReviews(user.UserID, status, 20)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "kandidat transaksi berhasil dimuat", resp)
}

func (r *Rest) ConfirmReview(c *gin.Context) {
	user, err := r.service.JwtAuth.GetLoginUser(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "user tidak ditemukan", err)
		return
	}

	reviewID, err := uuid.Parse(c.Param("review_id"))
	if err != nil {
		response.HandleError(c, apperr.BadRequest("review_id tidak valid"))
		return
	}

	var req model.ConfirmReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "request tidak valid", err)
		return
	}

	resp, err := r.service.TransactionReviewService.ConfirmReview(user.UserID, reviewID, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "transaksi berhasil dikonfirmasi", resp)
}

func (r *Rest) RejectReview(c *gin.Context) {
	user, err := r.service.JwtAuth.GetLoginUser(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "user tidak ditemukan", err)
		return
	}

	reviewID, err := uuid.Parse(c.Param("review_id"))
	if err != nil {
		response.HandleError(c, apperr.BadRequest("review_id tidak valid"))
		return
	}

	var req model.RejectReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "request tidak valid", err)
		return
	}

	resp, err := r.service.TransactionReviewService.RejectReview(user.UserID, reviewID, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "kandidat transaksi berhasil ditolak", resp)
}
