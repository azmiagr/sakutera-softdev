package rest

import (
	"net/http"
	"strconv"

	"github.com/azmiagr/sakutera-softdev/entity"
	"github.com/azmiagr/sakutera-softdev/model"
	apperr "github.com/azmiagr/sakutera-softdev/pkg/errors"
	"github.com/azmiagr/sakutera-softdev/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *Rest) CreatePairingCode(c *gin.Context) {
	user, err := r.service.JwtAuth.GetLoginUser(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "user tidak ditemukan", err)
		return
	}

	var req model.CreatePairingCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "request tidak valid", err)
		return
	}

	resp, err := r.service.CollectorService.CreatePairingCode(user.UserID, req)
	if err != nil {
		handleCollectorError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "kode pairing berhasil dibuat", resp)
}

func (r *Rest) PairDevice(c *gin.Context) {
	var req model.PairDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "request tidak valid", err)
		return
	}

	resp, err := r.service.CollectorService.PairDevice(req, c.ClientIP())
	if err != nil {
		handleCollectorError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "perangkat berhasil dihubungkan", resp)
}

// handleCollectorError adds a Retry-After header on 429 responses (per
// docs/notification.md §6.3) before delegating to the standard error handler.
func handleCollectorError(c *gin.Context, err error) {
	var appErr *apperr.AppError
	if e, ok := err.(*apperr.AppError); ok {
		appErr = e
	}
	if appErr != nil && appErr.Code == http.StatusTooManyRequests {
		c.Header("Retry-After", "600")
	}
	response.HandleError(c, err)
}

func (r *Rest) ListDevices(c *gin.Context) {
	user, err := r.service.JwtAuth.GetLoginUser(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "user tidak ditemukan", err)
		return
	}

	resp, err := r.service.CollectorService.ListDevices(user.UserID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "daftar perangkat berhasil dimuat", resp)
}

func (r *Rest) RevokeDevice(c *gin.Context) {
	user, err := r.service.JwtAuth.GetLoginUser(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "user tidak ditemukan", err)
		return
	}

	deviceID, err := uuid.Parse(c.Param("device_id"))
	if err != nil {
		response.HandleError(c, apperr.BadRequest("device_id tidak valid"))
		return
	}

	resp, err := r.service.CollectorService.RevokeDevice(user.UserID, deviceID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "akses perangkat berhasil dicabut", resp)
}

func (r *Rest) CollectorHealth(c *gin.Context) {
	device := c.MustGet("device").(*entity.Device)

	resp := r.service.CollectorService.HealthCheck(device)

	response.Success(c, http.StatusOK, "collector terhubung", resp)
}

func (r *Rest) GetConfig(c *gin.Context) {
	resp, err := r.service.CollectorService.GetConfig()
	if err != nil {
		response.HandleError(c, err)
		return
	}

	etag := "\"" + strconv.Itoa(resp.ConfigVersion) + "\""
	if match := c.GetHeader("If-None-Match"); match != "" && match == etag {
		c.Status(http.StatusNotModified)
		return
	}

	c.Header("ETag", etag)
	response.Success(c, http.StatusOK, "konfigurasi collector berhasil dimuat", resp)
}

func (r *Rest) UploadNotificationBatch(c *gin.Context) {
	device := c.MustGet("device").(*entity.Device)

	var req model.NotificationBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "request tidak valid", err)
		return
	}

	resp, err := r.service.CollectorService.IngestBatch(device, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "batch notifikasi selesai diproses", resp)
}
