package main

import (
	"log"

	"task/database"
	"task/handlers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.InitPostgres()
	database.InitRedis()

	app := fiber.New()

	handlers.RegisterRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
