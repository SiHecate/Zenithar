package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	ProductName     string        `json:"product_name"`
	ProductCategory string        `json:"product_category"`
	ProductPrice    float64       `json:"product_price"`
	OrderDetails    []OrderDetail `gorm:"foreignKey:ProductID"` // Bu alan, OrderDetail modelindeki ProductID'ye referans olmalıdır.
}
