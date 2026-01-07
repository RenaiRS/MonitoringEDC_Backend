package routes

import (
	"github.com/gofiber/fiber/v2"
	"backend/controllers"
	"backend/middlewares"
)

func OverdueRoutes(api fiber.Router) {
	overdue := api.Group("/overdue", middlewares.AuthMiddleware)
	overdue.Get("/summary", controllers.GetOverdueSummary)
	overdue.Get("/list", controllers.GetOverdueList)
	overdue.Get("/search", controllers.SearchOverdue)
}