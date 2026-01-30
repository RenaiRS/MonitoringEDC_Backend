package utils

import "github.com/gofiber/fiber/v2"

func Success(c *fiber.Ctx, data interface{}) error {
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   data,
	})
}

func Error(c *fiber.Ctx, message string) error {
	return c.Status(400).JSON(fiber.Map{
		"status":  "error",
		"message": message,
	})
}
