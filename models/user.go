package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"not null"`
	Email    string `gorm:"unique;size:10;not null" json:"email"`
	Password string `gorm:"unique;size:6;not null" json:"password"`
	Role     string `json:"role"`
}
