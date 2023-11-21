package controller

import (
	"Zenithar/core/database"
	"Zenithar/models"

	"github.com/gofiber/fiber/v2"
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
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	if existingTable.ID != 0 {
		return c.Status(422).JSON(fiber.Map{"error": "Table already exists"})
	}

	newTable := models.Table{
		TableNo:       tableData.TableNo,
		TablePassword: existingTable.TablePassword,
		Capacity:      existingTable.Capacity,
	}

	if err := database.Conn.Create(&newTable).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Table creation failed"})
	}

	return c.JSON(fiber.Map{"message": "Table created successfully", "table": newTable})

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
