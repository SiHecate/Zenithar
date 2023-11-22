package controller

import (
	"Zenithar/core/database"
	"Zenithar/models"

	"github.com/gofiber/fiber/v2"
)

/*
TakeOrder:

	TakeOrder, bir müşterinin bir masaya sipariş vermesini işler.
	İlgili masanın durumunu kontrol eder, sipariş verilebilir durumdaysa yeni bir sipariş oluşturur ve verilen ürünleri bu siparişe ekler.
	Eğer masanın durumu sipariş verilemezse veya masada zaten bir aktif sipariş varsa hata döner.
*/
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

	if existingTable.Status != "available" {
		return c.Status(500).JSON(fiber.Map{"error": "Table already busy."})
	}

	var existingOrder models.Order
	if err := database.Conn.Where("table_id = ? AND payment_status = false", existingTable.ID).First(&existingOrder).Error; err == nil {
		return c.Status(500).JSON(fiber.Map{"error": "There is already an active order for this table."})
	}

	newOrder := models.Order{
		TableNo:       orderRequest.TableNo,
		TableID:       existingTable.ID,
		PaymentStatus: false,
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

	MakeBusy(orderRequest.TableNo)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Order placed successfully",
		"data":    newOrder,
	})
}

/*
MakePayment:

	MakePayment, bir müşterinin bir masada yaptığı ödemeyi işler.
	İlgili masanın durumunu kontrol eder, eğer müşteri yoksa veya masada ödenmemiş bir sipariş yoksa hata döner.
	Eğer her şey yolunda ise, siparişin ödeme durumunu günceller.
*/
func MakePayment(c *fiber.Ctx) error {
	var paymentRequest struct {
		TableNo string `json:"tableno"`
	}

	if err := c.BodyParser(&paymentRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request data: " + err.Error(),
		})
	}

	var existingTable models.Table
	if err := database.Conn.Where("table_no", paymentRequest.TableNo).First(&existingTable).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Table not found"})
	}

	var existingOrder models.Order
	if err := database.Conn.Where("table_no", paymentRequest.TableNo).First(&existingOrder).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Order not found"})
	}

	if existingTable.Status != "busy" {
		return c.Status(500).JSON(fiber.Map{"error": "Table already available; there is no customer at this table."})
	}

	if existingOrder.PaymentStatus {
		return c.Status(500).JSON(fiber.Map{"error": "Table doesn't have an order"})
	}

	if err := database.Conn.Model(&existingOrder).Update("payment_status", true).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to update payment status",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Payment processed successfully",
		"data":    existingOrder,
	})
}

func ListOrders(c *fiber.Ctx) error {
	var existingOrders []models.Order

	if err := database.Conn.Find(&existingOrders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch orders", "message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    existingOrders,
	})
}
