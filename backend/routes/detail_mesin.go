package routes

import (
	"github.com/gofiber/fiber/v2"
	"backend/controllers"
	"backend/middlewares"
)

func DetailMesinRoutes(api fiber.Router) {
	detailMesin := api.Group("/detail-mesin", middlewares.AuthMiddleware)
	detailMesin.Get("/:id", controllers.GetDetailMesin)
	detailMesin.Put("/:id", controllers.UpdateMesin)
}