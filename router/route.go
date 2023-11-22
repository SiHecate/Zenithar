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
	product.Post("/update_product", controller.UpdateProduct)
	product.Post("/delete_product", controller.DeleteProduct)
	product.Get("/list_product", controller.ListProducts)

	table := app.Group("/table")
	table.Post("/add_table", controller.AddTable)
	table.Post("/update_table", controller.UpdateTable)
	table.Post("/delete_table", controller.DeleteTable)
	table.Get("/list_table", controller.ListTables)

	order := app.Group("/order")
	order.Post("/take_order", controller.TakeOrder)
}
