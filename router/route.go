package routes

import (
	"Zenithar/controller"

	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Zenithar say Hello world")
	})

	product := app.Group("/product")
	product.Post("/add_product", controller.AddProduct)

	table := app.Group("/table")
	table.Post("/add_table", controller.AddTable)
}
