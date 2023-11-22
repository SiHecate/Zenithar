package controller

import (
	"Zenithar/core/database"
	"Zenithar/models"
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AddTable(c *fiber.Ctx) error {
	var tableData struct {
		TableNo       string `json:"tableno"`
		TablePassword string `json:"table_password"`
		Capacity      int    `json:"capacity"`
	}

	if err := c.BodyParser(&tableData); err != nil {
		return c.Status(422).JSON(fiber.Map{"error": "Invalid request data: " + err.Error()})
	}

	var existingTable models.Table
	if err := database.Conn.Where("table_no", tableData.TableNo).First(&existingTable).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Record not found, proceed with table creation
		} else {
			return c.Status(500).JSON(fiber.Map{"error": "Database error"})
		}
	} else {
		// Table already exists
		return c.Status(422).JSON(fiber.Map{"error": "Table already exists"})
	}

	newTable := models.Table{
		TableNo:       tableData.TableNo,
		TablePassword: tableData.TablePassword,
		Capacity:      tableData.Capacity,
	}

	if err := database.Conn.Create(&newTable).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Table creation failed"})
	}

	return c.JSON(fiber.Map{"message": "Table created successfully", "table": newTable})
}

func UpdateTable(c *fiber.Ctx) error {
	var tableData struct {
		TableNo       string `json:"tableno"`
		TablePassword string `json:"table_password"`
		Capacity      int    `json:"capacity"`
	}

	if err := c.BodyParser(&tableData); err != nil {
		return c.Status(422).JSON(fiber.Map{"error": "Invalid request data: " + err.Error()})
	}

	var existingTable models.Table
	if err := database.Conn.Where("table_no", tableData.TableNo).First(&existingTable).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	if existingTable.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Table not found"})
	}

	// Masalara yeni özellikler girildikçe eklentiler gerekebilir.
	existingTable.Capacity = tableData.Capacity

	if err := database.Conn.Save(&existingTable).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Table update failed"})
	}

	return c.JSON(fiber.Map{"message": "Table updated successfully", "table": existingTable})
}

func DeleteTable(c *fiber.Ctx) error {
	tableNo := c.Params("tableNo")

	var existingTable models.Table
	if err := database.Conn.Where("table_no", tableNo).First(&existingTable).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Table not found"})
	}

	if err := database.Conn.Delete(&existingTable).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Table deletion failed"})
	}

	return c.JSON(fiber.Map{"message": "Table deleted successfully"})
}

func ListTables(c *fiber.Ctx) error {
	var existingTables []models.Table
	if err := database.Conn.Find(&existingTables).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	if len(existingTables) == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "No products found"})
	}

	return c.JSON(fiber.Map{"message": existingTables})

}

func TakeBusy(c *fiber.Ctx) error {
	var tableData struct {
		TableNo string `json:"tableno"`
	}

	if err := c.BodyParser(&tableData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request data: " + err.Error()})
	}

	var existingTable models.Table
	if err := database.Conn.Where("table_no = ?", tableData.TableNo).First(&existingTable).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	currentStatus := existingTable.Status

	if currentStatus != "busy" {
		if err := database.Conn.Model(&existingTable).Update("status", "busy").Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update table status"})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Table status updated to Busy",
		})
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Table is already busy"})
	}
}

func TakeAvailable(c *fiber.Ctx) error {
	var tableData struct {
		TableNo string `json:"tableno"`
	}

	if err := c.BodyParser(&tableData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request data: " + err.Error()})
	}

	var existingTable models.Table
	if err := database.Conn.Where("table_no = ?", tableData.TableNo).First(&existingTable).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	currentStatus := existingTable.Status

	if currentStatus == "busy" {
		if err := database.Conn.Model(&existingTable).Update("status", "available").Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update table status"})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Table status updated to Available",
		})
	} else if currentStatus == "Available" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Table is already available"})
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Table status is invalid"})
	}
}
