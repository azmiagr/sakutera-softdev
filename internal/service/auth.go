package service

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/azmiagr/sakutera-softdev/entity"
	"github.com/azmiagr/sakutera-softdev/internal/repository"
	"github.com/azmiagr/sakutera-softdev/model"
	"github.com/azmiagr/sakutera-softdev/pkg/database/mariadb"
	apperrors "github.com/azmiagr/sakutera-softdev/pkg/errors"
	"github.com/azmiagr/sakutera-softdev/pkg/jwt"
	"github.com/azmiagr/sakutera-softdev/pkg/whatsapp"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type IAuthService interface {
	Register(req model.RegisterRequest) (*model.RegisterResponse, error)
	VerifyOTP(sessionToken string, req model.VerifyOTPRequest) (*model.VerifyOTPResponse, error)
	GetUserByID(userID uuid.UUID) (*entity.User, error)
}

type AuthService struct {
	db          *gorm.DB
	userRepo    repository.IUserRepository
	sessionRepo repository.ISessionRepository
	otpRepo     repository.IOTPRepository
	jwtAuth     jwt.Interface
	whatsapp    whatsapp.Interface
}

func NewAuthService(userRepo repository.IUserRepository, sessionRepo repository.ISessionRepository, otpRepo repository.IOTPRepository, jwtAuth jwt.Interface) IAuthService {
	return &AuthService{
		db:          mariadb.Connection,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		otpRepo:     otpRepo,
		jwtAuth:     jwtAuth,
	}
}

func (s *AuthService) Register(req model.RegisterRequest) (*model.RegisterResponse, error) {
	existingUser, err := s.userRepo.GetUser(s.db, model.UserParam{
		PhoneNumber: req.PhoneNumber,
	})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, apperrors.InternalServer("gagal memeriksa nomor telepon")
	}

	var user *entity.User
	if existingUser != nil && existingUser.Status == "active" {
		return nil, apperrors.Conflict("nomor HP sudah terdaftar")
	}
	if existingUser != nil && existingUser.Status == "inactive" {
		user = existingUser
	} else {
		newUser := &entity.User{
			UserID:      uuid.New(),
			FullName:    req.FullName,
			PhoneNumber: req.PhoneNumber,
			Status:      "inactive",
		}
		if err := s.userRepo.CreateUser(s.db, newUser); err != nil {
			return nil, apperrors.InternalServer("gagal membuat akun")
		}
		user = newUser
	}

	_ = s.otpRepo.DeleteOTPByUserID(s.db, user.UserID.String())
	_ = s.sessionRepo.DeleteSessionByUserID(s.db, user.UserID.String())

	otpCode, err := generateOTPCode()
	if err != nil {
		return nil, apperrors.InternalServer("gagal membuat kode OTP")
	}

	otpEntity := &entity.OTP{
		OTPID:     uuid.New(),
		UserID:    user.UserID,
		Code:      otpCode,
		Type:      "registration",
		ExpiredAt: time.Now().Add(5 * time.Minute),
	}
	if err := s.otpRepo.CreateOTP(s.db, otpEntity); err != nil {
		return nil, apperrors.InternalServer("gagal menyimpan OTP")
	}

	message := fmt.Sprintf("Kode OTP Sakutera kamu adalah: *%s*\nBerlaku 5 menit. Jangan bagikan kode ini kepada siapapun.", otpCode)
	if err := s.whatsapp.SendMessage(req.PhoneNumber, message); err != nil {
		return nil, apperrors.Wrap(err, http.StatusServiceUnavailable, "gagal mengirim OTP via WhatsApp")
	}

	sessionToken, err := s.jwtAuth.CreateJWTToken(user.UserID, "session")
	if err != nil {
		return nil, apperrors.InternalServer("gagal membuat session token")
	}

	sessionEntity := &entity.Session{
		SessionID: uuid.New(),
		UserID:    user.UserID,
		Token:     sessionToken,
		ExpiredAt: time.Now().Add(5 * time.Minute),
	}
	if err := s.sessionRepo.CreateSession(s.db, sessionEntity); err != nil {
		return nil, apperrors.InternalServer("gagal menyimpan session")
	}
	return &model.RegisterResponse{
		SessionToken: sessionToken,
		PhoneMasked:  maskPhoneNumber(req.PhoneNumber),
		Message:      "Kode OTP telah dikirim via WhatsApp",
	}, nil
}

func (s *AuthService) VerifyOTP(sessionToken string, req model.VerifyOTPRequest) (*model.VerifyOTPResponse, error) {
	session, err := s.sessionRepo.GetSessionByToken(s.db, sessionToken)
	if err != nil {
		return nil, apperrors.Unauthorized("session tidak valid atau sudah kedaluwarsa")
	}

	if time.Now().After(session.ExpiredAt) {
		return nil, apperrors.Unauthorized("session sudah kedaluwarsa, silakan daftar ulang")
	}

	otp, err := s.otpRepo.GetOTP(s.db, model.OTPParam{
		UserID: session.UserID,
		Code:   req.Code,
		Type:   "registration",
	})
	if err != nil {
		return nil, apperrors.BadRequest("kode OTP tidak valid")
	}

	if time.Now().After(otp.ExpiredAt) {
		return nil, apperrors.BadRequest("kode OTP sudah kedaluwarsa, silakan kirim ulang")
	}

	user, err := s.userRepo.GetUser(s.db, model.UserParam{UserID: session.UserID})
	if err != nil {
		return nil, apperrors.InternalServer("gagal mengambil data pengguna")
	}
	user.Status = "active"
	if err := s.userRepo.UpdateUser(s.db, user); err != nil {
		return nil, apperrors.InternalServer("gagal mengaktifkan akun")
	}

	_ = s.otpRepo.DeleteOTPByUserID(s.db, user.UserID.String())
	_ = s.sessionRepo.DeleteSessionByUserID(s.db, user.UserID.String())

	token, err := s.jwtAuth.CreateJWTToken(user.UserID, "user")
	if err != nil {
		return nil, apperrors.InternalServer("gagal membuat token autentikasi")
	}
	return &model.VerifyOTPResponse{
		Token:   token,
		Message: "Akun berhasil diverifikasi",
	}, nil
}

func (s *AuthService) GetUserByID(userID uuid.UUID) (*entity.User, error) {
	user, err := s.userRepo.GetUser(s.db, model.UserParam{UserID: userID})
	if err != nil {
		return nil, apperrors.NotFound("user tidak ditemukan")
	}
	return user, nil
}

func generateOTPCode() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func maskPhoneNumber(phone string) string {
	if len(phone) < 8 {
		return phone
	}
	visible := 4
	masked := strings.Repeat("*", len(phone)-visible*2)
	return phone[:visible] + masked + phone[len(phone)-visible:]
}
