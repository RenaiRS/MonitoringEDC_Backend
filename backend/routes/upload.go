package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func ExcelRoutes(api fiber.Router) {
	excel := api.Group("/excel", middlewares.AuthMiddleware)

	excel.Post("/upload-vendor", controllers.UploadVendorExcel)
	excel.Post("/upload-bank", controllers.UploadBankExcel)
}
