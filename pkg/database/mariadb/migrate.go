package mariadb

import (
	"github.com/azmiagr/sakutera-softdev/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&entity.WorkCategory{},
		&entity.WorkPlatform{},
		&entity.User{},
		&entity.TransactionSource{},
		&entity.Transaction{},
		&entity.Attachment{},
		&entity.ForecastResult{},
		&entity.IncomePassport{},
		&entity.Organization{},
		&entity.AccessLog{},
		&entity.Consent{},
		&entity.Session{},
		&entity.OTP{},
		&entity.Notification{},
		&entity.TokenBlacklist{},
	)

	if err != nil {
		return err
	}

	return nil
}
