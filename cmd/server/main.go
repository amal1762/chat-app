package main

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/amal1762/chat-app/internal/router"
)

func main() {
	app := fiber.New()
	router.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}