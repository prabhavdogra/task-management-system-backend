package main

import (
	"log"
	"to-do-backend/database"
	"to-do-backend/routes"

	"github.com/gofiber/fiber/v2"
)

func welcome(c *fiber.Ctx) error {
	return c.SendString("Welcome to my awesome API")
}

func setupRoutes(app *fiber.App) {
	app.Get("/api", welcome)
	// Authentication Endpoints
	app.Post("/api/auth/login", routes.Login)
	app.Post("/api/auth/signup", routes.Signup)
	// User Endpoints
	app.Post("/api/users", routes.CreateUser)
	// app.Get("/api/users", routes.GetUsers)
	// app.Get("/api/users/:id", routes.GetUser)
	// app.Post("/api/users/:id", routes.UpdateUser)
	// app.Delete("/api/users/:id", routes.DeleteUser)

	// // Task Endpoints
	// app.Post("/api/products", routes.CreateProduct)
	// app.Get("/api/products", routes.GetProducts)
	// app.Get("/api/products/:id", routes.GetProduct)
	// app.Put("/api/products/:id", routes.UpdateProduct)
}

func main() {
	database.ConnectDb()
	app := fiber.New()
	setupRoutes(app)
	log.Fatal(app.Listen(":3002"))
}
