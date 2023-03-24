package main

import (
	"log"
	"to-do-backend/database"
	"to-do-backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Health(c *fiber.Ctx) error {
	return c.Status(200).SendString("Healthy")
}

func setupRoutes(app *fiber.App) {
	app.Get("/api/health", Health)

	// Authentication Endpoints
	app.Post("/api/auth/login", routes.Login)
	app.Post("/api/auth/signup", routes.Signup)
	app.Post("/api/auth/authenticate", routes.AuthenticateJWTToken)
	app.Post("/api/auth/logout", routes.Logout)

	// User Endpoints
	app.Post("/api/user/create", routes.CreateUser)
	app.Get("/api/user/get", routes.GetUserByJWTToken)
	app.Post("/api/user/update", routes.UpdateUser)

	// Task Endpoints
	app.Post("/api/task/create", routes.CreateTask)
	app.Get("/api/task", routes.GetAllTasks)
	app.Post("/api/task/update", routes.UpdateTask)
	app.Delete("/api/task/:id", routes.DeleteTask)
}

func main() {
	database.ConnectDb()
	app := fiber.New()

	// Default config
	app.Use(cors.New())

	// For CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	setupRoutes(app)
	log.Fatal(app.Listen("0.0.0.0:4000"))
}
