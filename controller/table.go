package controller

import (
	"Zenithar/core/database"
	"Zenithar/models"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

/*
AddTable:

	AddTable, yeni bir masa ekler.
	İlgili masa daha önce varsa hata döner, aksi takdirde yeni masa oluşturur ve veritabanına ekler.
*/
func AddTable(c *fiber.Ctx) error {
	var tableData struct {
		TableNo  string `json:"tableno"`
		Capacity int    `json:"capacity"`
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
		TableNo:  tableData.TableNo,
		Capacity: tableData.Capacity,
	}

	if err := database.Conn.Create(&newTable).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Table creation failed"})
	}

	return c.JSON(fiber.Map{"message": "Table created successfully", "table": newTable})
}

/*
UpdateTable:

	UpdateTable, mevcut bir masa bilgisini günceller.
	İlgili masa bilgilerini veritabanında bulur, gelen verilerle günceller ve tekrar veritabanına kaydeder.
*/
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

/*
DeleteTable:

	DeleteTable, bir masa bilgisini veritabanından siler.
	İlgili masa bilgisini masa numarasına göre bulur ve siler.
*/
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

/*
ListTables:

	ListTables, tüm masaları listeler.
	Veritabanından tüm masaları çeker ve istemciye gönderir.
*/
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

/*
MakeBusy:

	MakeBusy, belirli bir masa numarasındaki masanın durumunu "busy" olarak günceller.
*/
func MakeBusy(TableNo string) error {
	var existingTable models.Table
	if err := database.Conn.Where("table_no = ?", TableNo).First(&existingTable).Error; err != nil {
		return fmt.Errorf("error finding table: %v", err)
	}

	currentStatus := existingTable.Status

	if currentStatus != "busy" {
		if err := database.Conn.Model(&existingTable).Update("status", "busy").Error; err != nil {
			return fmt.Errorf("failed to update table status: %v", err)
		}

		return nil
	}

	return fmt.Errorf("table is already busy")
}

/*
MakeAvailable:

	MakeAvailable, belirli bir masa numarasındaki masanın durumunu "available" olarak günceller.
*/
func MakeAvailable(TableNo string) error {
	var existingTable models.Table
	if err := database.Conn.Where("table_no = ?", TableNo).First(&existingTable).Error; err != nil {
		return fmt.Errorf("error finding table: %v", err)
	}

	currentStatus := existingTable.Status

	if currentStatus == "busy" {
		if err := database.Conn.Model(&existingTable).Update("status", "available").Error; err != nil {
			return fmt.Errorf("failed to update table status: %v", err)
		}

		return nil
	}

	return fmt.Errorf("table is already available")
}
