package mariadb

import (
	"github.com/azmiagr/sakutera-softdev/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Seed(db *gorm.DB) error {
	categories := []struct {
		name      string
		platforms []string
	}{
		{"Pengemudi Ojol", []string{"Gojek", "Grab", "Maxim", "InDriver"}},
		{"Kurir Online", []string{"ShopeeFood", "GoSend", "J&T", "JNE"}},
		{"Freelancer", []string{"Fiverr", "Upwork", "Freelancer.com"}},
		{"Pedagang UMKM", []string{"Tokopedia", "Shopee", "Warung Offline"}},
		{"Pekerjaan Lainnya", []string{"Lainnya"}},
	}

	for _, cat := range categories {
		category := entity.WorkCategory{
			WorkCategoryID: uuid.New(),
			Name:           cat.name,
		}
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&category).Error; err != nil {
			return err
		}

		var existingCat entity.WorkCategory
		if err := db.Where("name = ?", cat.name).First(&existingCat).Error; err != nil {
			return err
		}

		for _, pName := range cat.platforms {
			platform := entity.WorkPlatform{
				WorkPlatformID: uuid.New(),
				WorkCategoryID: existingCat.WorkCategoryID,
				Name:           pName,
			}
			if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&platform).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
