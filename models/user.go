package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"not null"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"unique;not null" json:"password"`
	Role     string `json:"role"`
}
