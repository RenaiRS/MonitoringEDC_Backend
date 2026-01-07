package routes

import (
	"backend/controllers"

	"github.com/gofiber/fiber/v2"
)

func authRoutes(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Post("/login", controllers.Login)
}
