package routes

import (
	"Zenithar/controller"
	"Zenithar/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Router(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST, PATCH, DELETE",
		AllowCredentials: true,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Zenithar says Hello world")
	})

	auth := app.Group("/auth")
	auth.Post("/signup", controller.Signup)
	auth.Post("/login", controller.Login)
	auth.Post("/logout", controller.Logout)

	product := app.Group("/product", middlewares.IsAuthorized(), middlewares.Authorization())
	product.Post("/add", controller.AddProduct)
	product.Post("/update", controller.UpdateProduct)
	product.Delete("/delete/:product_name", controller.DeleteProduct)
	product.Get("/list", controller.ListProducts)

	table := app.Group("/table", middlewares.IsAuthorized(), middlewares.Authorization())
	table.Post("/add", controller.AddTable)
	table.Post("/update", controller.UpdateTable)
	table.Delete("/delete/:tableNo", controller.DeleteTable)
	table.Get("/list", controller.ListTables)

	order := app.Group("/order", middlewares.IsAuthorized(), middlewares.Authorization())
	order.Post("/take", controller.TakeOrder)
	order.Post("/make_payment", controller.MakePayment)
	order.Get("/list", controller.ListOrders)
}
