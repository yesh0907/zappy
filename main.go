package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"zappy.sh/database"
	"zappy.sh/routes"
)

func setUpRoutes(app *fiber.App, handler *routes.Handler) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})
	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})

	app.Get("/:alias", handler.GetAlias)
	app.Post("/alias/create", handler.CreateAlias)
	app.Get("/requests/:alias", handler.AllRequests)
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to database
	database.ConnectDB()

	// Start fiber app
	app := fiber.New()

	// Handle routes
	handler := routes.NewHandler()
	setUpRoutes(app, handler)

	// Set default 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})

	// Listen on port 8080
	log.Fatal(app.Listen(":8080"))
}
