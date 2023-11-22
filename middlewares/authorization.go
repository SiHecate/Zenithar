package middlewares

import (
	"Zenithar/core/database"
	"Zenithar/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Authorization() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		role := c.Locals("Role")
		id := c.Locals("Issuer")

		var user models.User
		if err := database.Conn.Model(&user).Where("id = ?", id).First(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
		}

		if role != "Admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
		}

		fmt.Println("Admin role: " + user.Role)
		fmt.Println("Admin username: " + user.Name)

		return c.Next()
	}
}
