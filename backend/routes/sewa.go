package routes

import (
	"github.com/gofiber/fiber/v2"
	"backend/controllers"
	"backend/middlewares"
)	

func sewaRoutes(api fiber.Router) {
	sewa := api.Group("/sewa", middlewares.AuthMiddleware)

	sewa.Get("/summary", controllers.GetSewaSummary)
	sewa.Get("/list", controllers.GetSewaList)
	sewa.Get("/search", controllers.SearchSewa)
}