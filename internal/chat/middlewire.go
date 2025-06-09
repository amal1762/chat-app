package chat

import (
	"fmt"
	"strings"

	"github.com/amal1762/chat-app/internal/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func WsAuthMiddleware(c *fiber.Ctx) error {

	authHeader := c.Get("Authorization")
	var tokenString string

	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		tokenString = strings.TrimPrefix(authHeader, "Bearer ")
	} else {
		// ✅ Most WebSocket clients send token as query param
		tokenString = c.Query("token")
	}

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	fmt.Println("WebSocket Token:", tokenString)
	token, err := jwt.ParseWithClaims(tokenString, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return auth.JwtKey, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	claims, ok := token.Claims.(*auth.Claims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	// ✅ Store claims in context
	c.Locals("user", claims)

	// ✅ Allow the WebSocket handler to run
	return c.Next()
}
