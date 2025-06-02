package chat

import "github.com/gofiber/fiber/v2"

func SetupRoutes(r fiber.Router) {
	r.Get("/ws", HandleWebSocket)
}

func HandleWebSocket(c *fiber.Ctx) error {
	return c.SendString("WebSocket placeholder")
}
