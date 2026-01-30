package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	authRoutes(app)

	api := app.Group("/api")

	DashboardRoutes(api)
	OverdueRoutes(api)
	RekapRoutes(api)
	sewaRoutes(api)
	DetailMesinRoutes(api)
	ExcelRoutes(api)
}
