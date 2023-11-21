package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	TableID       uint
	PaymentStatus string  `json:"payment_status" gorm:"not null;default:'unpaid'"`
	DueAmount     float64 `json:"due_amount"`
	PaymentType   string  `json:"payment_type"`
	OrderDetails  []struct {
		gorm.Model
		Product  Product `json:"product" gorm:"foreignKey:ProductID"`
		Quantity int     `json:"quantity"`
		Price    float64 `json:"price"`
	} `json:"order_details"`
}

/*
	Başla

	Bir restorana girildi.
	Kullanıcı masa seçti.
	Kullanıcı masaya oturdu Masa no = 1
	Masaya görevli geldi
	Görevli menüyü verdi
	Görevli masanın durumunu doluya çekti
	Kullanıcı menüden seçim yaptı
	Kullanıcı seçilen ürünü görevliye söyledi
	Görevli kullanıcının istediği ürün yazdı
	Görevli mutfağa gidip ürünü söyledi ve ürün hazırlanmaya başladı
	Görevli masanın borcunu yazdı
	Görevli ödenme durumunu ödenmedi olarak ayarladı (bunu otomatik yapacağız)
	Kullanıcının ürünü geldi
	Kullanıcı işini bitirip görevliye gitti
	Kullanıcı ödeme yöntemini söyledi
	Görevli ödeme yöntemini kaydetti
	Kullanıcı borcunu ödedi
	Görevli ödenme durumunu ödendi olarak değiştirdi
	Masanın görevini bitirdi
	Masadan alınan ürün, ödenme yöntemi, ödenme durumu, sipariş tarihi, ödenme tarihi gibi veriler farklı bir sayfaya kaydedildi
	Masayı tekrardan müsait durumuna getirdi

	Bitir
*/
