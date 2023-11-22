package controller

import (
	"Zenithar/core/database"
	"Zenithar/models"

	"github.com/gofiber/fiber/v2"
)

func TakeOrder(c *fiber.Ctx) error {
	var orderRequest struct {
		TableNo  string `json:"tableno"`
		Products []struct {
			ProductName string `json:"product_name"`
			Quantity    int    `json:"quantity"`
		} `json:"products"`
	}

	if err := c.BodyParser(&orderRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request data: " + err.Error(),
		})
	}

	var existingTable models.Table
	if err := database.Conn.Where("table_no", orderRequest.TableNo).First(&existingTable).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Table not found"})
	}

	newOrder := models.Order{
		TableNo: orderRequest.TableNo,
		TableID: existingTable.ID,
	}

	for _, productRequest := range orderRequest.Products {
		var existingProduct models.Product
		if err := database.Conn.Where("product_name = ?", productRequest.ProductName).First(&existingProduct).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Product not found"})
		}

		Price := existingProduct.ProductPrice * float64(productRequest.Quantity)

		orderDetail := models.OrderDetail{
			Product:  existingProduct,
			Quantity: productRequest.Quantity,
			Price:    Price,
			TableID:  existingTable.ID,
		}

		newOrder.OrderDetails = append(newOrder.OrderDetails, orderDetail)
	}

	if err := database.Conn.Create(&newOrder).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to save order to the database",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Order placed successfully",
		"data":    newOrder,
	})
}
