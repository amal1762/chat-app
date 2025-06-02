package history

import "github.com/gofiber/fiber/v2"

func SetupRoutes(r fiber.Router) {
	r.Get("/messages", GetMessages)
}

func GetMessages(c *fiber.Ctx) error {
	return c.SendString("Fetch messages placeholder")
}