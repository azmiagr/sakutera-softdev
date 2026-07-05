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
	"github.com/azmiagr/sakutera-softdev/pkg/bcrypt"
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
	CheckPhoneNumber(phone string) (*model.CheckPhoneResponse, error)
	SetPIN(sessionToken string, pin string) (*model.SetPINResponse, error)
	LoginWithPIN(phone string, pin string) (*model.LoginResponse, error)
	Logout(token string) (*model.LogoutResponse, error)
	IsTokenRevoked(token string) (bool, error)
}

type AuthService struct {
	db                 *gorm.DB
	userRepo           repository.IUserRepository
	sessionRepo        repository.ISessionRepository
	otpRepo            repository.IOTPRepository
	tokenBlacklistRepo repository.ITokenBlacklistRepository
	jwtAuth            jwt.Interface
	bcrypt             bcrypt.Interface
	whatsapp           whatsapp.Interface
}

func NewAuthService(userRepo repository.IUserRepository, sessionRepo repository.ISessionRepository, otpRepo repository.IOTPRepository, tokenBlacklistRepo repository.ITokenBlacklistRepository, jwtAuth jwt.Interface, bcrypt bcrypt.Interface, wa whatsapp.Interface) IAuthService {
	return &AuthService{
		db:                 mariadb.Connection,
		userRepo:           userRepo,
		sessionRepo:        sessionRepo,
		otpRepo:            otpRepo,
		tokenBlacklistRepo: tokenBlacklistRepo,
		jwtAuth:            jwtAuth,
		bcrypt:             bcrypt,
		whatsapp:           wa,
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
	err = s.userRepo.UpdateUser(s.db, user)
	if err != nil {
		return nil, apperrors.InternalServer("gagal mengaktifkan akun")
	}

	err = s.otpRepo.DeleteOTPByUserID(s.db, user.UserID.String())
	if err != nil {
		return nil, apperrors.InternalServer("gagal menghapus OTP")
	}
	err = s.sessionRepo.DeleteSessionByUserID(s.db, user.UserID.String())
	if err != nil {
		return nil, apperrors.InternalServer("gagal menghapus session")
	}

	token, err := s.jwtAuth.CreateJWTToken(user.UserID, "user")
	if err != nil {
		return nil, apperrors.InternalServer("gagal membuat token autentikasi")
	}
	return &model.VerifyOTPResponse{
		Token:   token,
		Message: "Akun berhasil diverifikasi",
	}, nil
}

func (s *AuthService) CheckPhoneNumber(phone string) (*model.CheckPhoneResponse, error) {
	user, err := s.userRepo.GetUser(s.db, model.UserParam{PhoneNumber: phone})
	if err != nil {
		return nil, apperrors.NotFound("nomor HP tidak terdaftar")
	}
	if user.Status != "active" {
		return nil, apperrors.BadRequest("akun belum diverifikasi, silakan daftar terlebih dahulu")
	}
	if user.PINNumber != "" {
		return &model.CheckPhoneResponse{HasPIN: true}, nil
	}

	_ = s.sessionRepo.DeleteSessionByUserID(s.db, user.UserID.String())
	sessionToken, err := s.jwtAuth.CreateJWTToken(user.UserID, "session")
	if err != nil {
		return nil, apperrors.InternalServer("gagal membuat session token")
	}

	sessionEntity := &entity.Session{
		SessionID: uuid.New(),
		UserID:    user.UserID,
		Token:     sessionToken,
		ExpiredAt: time.Now().Add(10 * time.Minute),
	}
	err = s.sessionRepo.CreateSession(s.db, sessionEntity)
	if err != nil {
		return nil, apperrors.InternalServer("gagal menyimpan session")
	}

	return &model.CheckPhoneResponse{HasPIN: false, SessionToken: sessionToken}, nil
}

func (s *AuthService) SetPIN(sessionToken string, pin string) (*model.SetPINResponse, error) {
	session, err := s.sessionRepo.GetSessionByToken(s.db, sessionToken)
	if err != nil || time.Now().After(session.ExpiredAt) {
		return nil, apperrors.Unauthorized("session tidak valid atau sudah kedaluwarsa")
	}
	if !isValidPIN(pin) {
		return nil, apperrors.BadRequest("PIN harus 6 digit angka")
	}

	user, err := s.userRepo.GetUser(s.db, model.UserParam{UserID: session.UserID})
	if err != nil {
		return nil, apperrors.InternalServer("gagal mengambil data pengguna")
	}

	hashed, err := s.bcrypt.GenerateFromPassword(pin)
	if err != nil {
		return nil, apperrors.InternalServer("gagal mengenkripsi PIN")
	}
	user.PINNumber = hashed
	if err := s.userRepo.UpdateUser(s.db, user); err != nil {
		return nil, apperrors.InternalServer("gagal menyimpan PIN")
	}

	_ = s.sessionRepo.DeleteSessionByUserID(s.db, user.UserID.String())

	token, err := s.jwtAuth.CreateJWTToken(user.UserID, "user")
	if err != nil {
		return nil, apperrors.InternalServer("gagal membuat token autentikasi")
	}

	return &model.SetPINResponse{Token: token, Message: "PIN berhasil dibuat"}, nil
}

func (s *AuthService) LoginWithPIN(phone string, pin string) (*model.LoginResponse, error) {
	user, err := s.userRepo.GetUser(s.db, model.UserParam{PhoneNumber: phone})
	if err != nil {
		return nil, apperrors.NotFound("nomor HP tidak terdaftar")
	}
	if user.Status != "active" {
		return nil, apperrors.BadRequest("akun belum diverifikasi, silakan daftar terlebih dahulu")
	}
	if user.PINNumber == "" {
		return nil, apperrors.BadRequest("PIN belum dibuat, silakan buat PIN terlebih dahulu")
	}

	err = s.bcrypt.CompareAndHashPassword(user.PINNumber, pin)
	if err != nil {
		return nil, apperrors.Unauthorized("PIN tidak valid")
	}

	token, err := s.jwtAuth.CreateJWTToken(user.UserID, "user")
	if err != nil {
		return nil, apperrors.InternalServer("gagal membuat token autentikasi")
	}
	return &model.LoginResponse{Token: token, Message: "Login berhasil"}, nil
}

func (s *AuthService) Logout(token string) (*model.LogoutResponse, error) {
	expiredAt, err := s.jwtAuth.GetTokenExpiry(token)
	if err != nil {
		return nil, apperrors.Unauthorized("token tidak valid")
	}

	blacklist := &entity.TokenBlacklist{
		BlacklistID: uuid.New(),
		Token:       token,
		ExpiredAt:   expiredAt,
	}
	if err := s.tokenBlacklistRepo.CreateBlacklistToken(s.db, blacklist); err != nil {
		return nil, apperrors.InternalServer("gagal logout")
	}

	return &model.LogoutResponse{Message: "Logout berhasil"}, nil
}

func (s *AuthService) IsTokenRevoked(token string) (bool, error) {
	return s.tokenBlacklistRepo.IsTokenBlacklisted(s.db, token)
}

func isValidPIN(pin string) bool {
	if len(pin) != 6 {
		return false
	}
	for _, c := range pin {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
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
