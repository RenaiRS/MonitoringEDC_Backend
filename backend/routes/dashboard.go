package routes

import (
	"github.com/gofiber/fiber/v2"
	"backend/controllers"
	"backend/middlewares"
)

func DashboardRoutes(api fiber.Router) {
	dashboard := api.Group("/dashboard", middlewares.AuthMiddleware)

	// === STATISTIK UTAMA ===
	dashboard.Get("/total-mesin", controllers.GetTotalMesin)
	dashboard.Get("/terdata-bank", controllers.GetMesinTerdataBank)
	// dashboard.Get("/vendor-only", controllers.GetMesinVendorOnly)

	// === STATUS ===
	dashboard.Get("/status-mesin", controllers.GetStatusMesin)
	dashboard.Get("/status-overdue", controllers.GetOverdueStatus)

	// === NOTIFIKASI ===
	dashboard.Get("/mesin-baru-vendor", controllers.GetMesinBaruVendor)
	dashboard.Get("/monitoring-overdue", controllers.GetMonitoringOverdue)
}
