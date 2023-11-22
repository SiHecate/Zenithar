package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	TableID      uint          `gorm:"not null" json:"-"`
	TableNo      string        `json:"tableno" gorm:"not null"`
	OrderDetails []OrderDetail `json:"order_details" gorm:"foreignKey:OrderID"`
}

type OrderDetail struct {
	gorm.Model
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	TableID   uint    `json:"table_id"`
	OrderID   uint    `gorm:"foreignKey:OrderID"` // Bu alan, Order modelindeki primary key'e referans olmalıdır.
	ProductID uint    // Bu alan, Product modelindeki primary key'e referans olmalıdır.
}

type OrderScreen struct {
	gorm.Model
	TableNo     string  `json:"tableno" gorm:"not null"`
	DueAmount   float64 `json:"due_amount"`
	PaymentType string  `json:"payment_type"`
}
