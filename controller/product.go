package controller

import (
	"Zenithar/core/database"
	"Zenithar/models"
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

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
