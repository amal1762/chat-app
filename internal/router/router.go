package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/amal1762/chat-app/internal/auth"
	"github.com/amal1762/chat-app/internal/chat"
	"github.com/amal1762/chat-app/internal/history"
)

func SetupRoutes(app *fiber.App) {
	authGroup := app.Group("/auth")
	chatGroup := app.Group("/chat")
	historyGroup := app.Group("/history")

	auth.SetupRoutes(authGroup)
	chat.SetupRoutes(chatGroup)
	history.SetupRoutes(historyGroup)
}
