package middlewares

import "github.com/gofiber/fiber/v2"

func AuthMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")

	if token == "" {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Validasi token (dummy)
	if token != "Bearer dummy-token" {
		return c.Status(401).JSON(fiber.Map{
			"message": "Invalid token",
		})
	}

	return c.Next()
}
