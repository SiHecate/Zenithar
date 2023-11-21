package models

import (
	"gorm.io/gorm"
)

type Table struct {
	gorm.Model
	TableNo       string  `json:"tableno" gorm:"not null"`
	TablePassword string  `json:"table_password" gorm:"not null"`
	Capacity      int     `json:"capacity" gorm:"not null"`
	Status        string  `json:"status" gorm:"not null;default:'available'"`
	Orders        []Order `gorm:"foreignKey:TableID"`
}
