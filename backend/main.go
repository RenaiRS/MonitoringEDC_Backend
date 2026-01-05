package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"backend/config"
	"backend/database"
	"backend/routes"
)

func main() {
	config.LoadEnv()
	database.DBLoad()
	database.DBMigrate()

	app := fiber.New()

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
