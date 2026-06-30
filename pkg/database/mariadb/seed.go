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

	if err := seedTransactionSources(db); err != nil {
		return err
	}

	return seedOrganizations(db)
}

func seedTransactionSources(db *gorm.DB) error {
	sources := []entity.TransactionSource{
		{TransactionSourceID: uuid.New(), Name: "Gojek", Provider: "GoPay", IsActive: true},
		{TransactionSourceID: uuid.New(), Name: "Grab", Provider: "OVO", IsActive: true},
		{TransactionSourceID: uuid.New(), Name: "ShopeeFood", Provider: "ShopeePay", IsActive: true},
		{TransactionSourceID: uuid.New(), Name: "GoSend", Provider: "GoPay", IsActive: true},
		{TransactionSourceID: uuid.New(), Name: "Tokopedia", Provider: "Manual", IsActive: true},
		{TransactionSourceID: uuid.New(), Name: "Lainnya", Provider: "Manual", IsActive: true},
	}

	for _, src := range sources {
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&src).Error; err != nil {
			return err
		}
	}

	return nil
}

func seedOrganizations(db *gorm.DB) error {
	orgs := []entity.Organization{
		{OrganizationID: uuid.New(), Name: "Koperasi Sejahtera Jawa", Type: "fintech"},
		{OrganizationID: uuid.New(), Name: "PT BPR Sentosa Digital", Type: "bank"},
		{OrganizationID: uuid.New(), Name: "Bank XYZ Leasing", Type: "bank"},
		{OrganizationID: uuid.New(), Name: "KoinWorks", Type: "fintech"},
		{OrganizationID: uuid.New(), Name: "Amartha Mikro Fintek", Type: "fintech"},
		{OrganizationID: uuid.New(), Name: "PT Astra Credit Companies", Type: "employer"},
	}

	for _, org := range orgs {
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&org).Error; err != nil {
			return err
		}
	}

	return nil
}
