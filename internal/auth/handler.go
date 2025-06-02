package auth

import "github.com/gofiber/fiber/v2"

func SetupRoutes(r fiber.Router) {
	r.Post("/login", Login)
	r.Post("/register", Register)
}

func Login(c *fiber.Ctx) error {
	// Placeholder login handler
	return c.SendString("Login endpoint")
}

func Register(c *fiber.Ctx) error {
	// Placeholder register handler
	return c.SendString("Register endpoint")
}