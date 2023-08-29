package main

import (
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
	"zappy.sh/database"
	"zappy.sh/middleware"
	"zappy.sh/routes"
)

var authMiddleware *middleware.AuthMiddleware

func setUpRoutes(app *fiber.App, handler *routes.Handler) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})
	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})

	app.Get("/:alias", handler.GetAlias)
	app.Post("/alias/create", handler.CreateAlias)
	app.Get("/requests/:alias", authMiddleware.Middleware, handler.AllRequests)
}

func setupMiddleware() {
	authMiddleware = middleware.NewAuthMiddleware(os.Getenv("API_KEY"))
}

func main() {
	// Connect to database
	database.ConnectDB()

	// Start fiber app
	app := fiber.New()

	// Set up middleware
	setupMiddleware()

	// Handle routes
	handler := routes.NewHandler()
	setUpRoutes(app, handler)

	// Set default 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})

	// Listen on port 8080
	log.Fatal(app.Listen("0.0.0.0:8080"))
}
