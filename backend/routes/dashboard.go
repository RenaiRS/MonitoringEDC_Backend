package routes

import (
	"github.com/gofiber/fiber/v2"
	"backend/controllers"
	"backend/middlewares"
)

func DashboardRoutes(api fiber.Router) {
	dashboard := api.Group("/dashboard", middlewares.AuthMiddleware)

	dashboard.Get("/", controllers.GetDashboard)
}

