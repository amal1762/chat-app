package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetTokenFromRequest(c *fiber.Ctx) string {
	authHeader := c.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	return c.Query("token")
}
