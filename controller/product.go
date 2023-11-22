package controller

import (
	"Zenithar/core/database"
	"Zenithar/models"
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

/*
AddProduct:

	AddProduct, yeni bir ürün ekler.
	İlgili ürün daha önce varsa hata döner, aksi takdirde yeni ürünü oluşturur ve veritabanına ekler.
*/
func AddProduct(c *fiber.Ctx) error {
	var productData struct {
		ProductName     string  `json:"product_name"`
		ProductCategory string  `json:"product_category"`
		ProductPrice    float64 `json:"product_price"`
	}

	if err := c.BodyParser(&productData); err != nil {
		return c.Status(422).JSON(fiber.Map{"error": "Invalid request data: " + err.Error()})
	}

	var existingProduct models.Product
	if err := database.Conn.Where("product_name = ?", productData.ProductName).First(&existingProduct).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	if existingProduct.ID != 0 {
		return c.Status(422).JSON(fiber.Map{"error": "Product already exists"})
	}

	newProduct := models.Product{
		ProductName:     productData.ProductName,
		ProductCategory: productData.ProductCategory,
		ProductPrice:    productData.ProductPrice,
	}

	if err := database.Conn.Create(&newProduct).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Product creation failed"})
	}

	return c.JSON(fiber.Map{"message": "Product created successfully", "product": newProduct})
}

/*
UpdateProduct:

	UpdateProduct, mevcut bir ürünün bilgilerini günceller.
	İlgili ürünü veritabanında bulur, gelen verilerle günceller ve tekrar veritabanına kaydeder.
*/
func UpdateProduct(c *fiber.Ctx) error {
	var productData struct {
		ProductName     string  `json:"product_name"`
		ProductCategory string  `json:"product_category"`
		ProductPrice    float64 `json:"product_price"`
	}

	if err := c.BodyParser(&productData); err != nil {
		return c.Status(422).JSON(fiber.Map{"error": "Invalid request data: " + err.Error()})
	}

	var existingProduct models.Product
	if err := database.Conn.Where("product_name = ?", productData.ProductName).First(&existingProduct).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	// Ürünlere yeni özellikler girildikçe eklentiler gerekebilir.
	existingProduct.ProductCategory = productData.ProductCategory
	existingProduct.ProductPrice = productData.ProductPrice

	if err := database.Conn.Save(&existingProduct).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Product update failed"})
	}

	return c.JSON(fiber.Map{"message": "Product updated successfully", "product": existingProduct})
}

/*
DeleteProduct:

	DeleteProduct, bir ürünü veritabanından siler.
	İlgili ürünü ürün adına göre bulur ve siler.
*/
func DeleteProduct(c *fiber.Ctx) error {
	ProductName := c.Params("product_name")

	var existingProduct models.Product
	if err := database.Conn.Where("product_name", ProductName).First(&existingProduct).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
	}

	if err := database.Conn.Delete(&existingProduct).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Product deletion failed"})
	}

	return c.JSON(fiber.Map{"message": "Product deleted successfully"})
}

/*
ListProducts:

	ListProducts, tüm ürünleri listeler.
	Veritabanından tüm ürünleri çeker ve istemciye gönderir.
*/
func ListProducts(c *fiber.Ctx) error {
	var existingProducts []models.Product
	if err := database.Conn.Find(&existingProducts).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	if len(existingProducts) == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "No products found"})
	}

	return c.JSON(fiber.Map{"message": existingProducts})
}
