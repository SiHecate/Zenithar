package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) error {
	adminPassword := "Admin"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	SeedData := []User{
		{
			Name:     "Admin",
			Email:    "admin@admin.com",
			Password: string(hashedPassword),
			Role:     "Admin",
		},
	}

	for _, user := range SeedData {
		var existingUser User
		if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if existingUser.ID == 0 {
			if err := db.Create(&user).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
