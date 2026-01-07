package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	req := new(LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	if req.Username == "admin" && req.Password == "admin123" {
		return c.JSON(fiber.Map{
			"message": "Login success",
			"token":   "dummy-token",
		})
	}

	return c.Status(401).JSON(fiber.Map{
		"message": "Invalid credentials",
	})
}
