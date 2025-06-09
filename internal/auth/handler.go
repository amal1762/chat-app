package auth

import (
	"os"

	"time"

	"github.com/amal1762/chat-app/internal/db"
	"github.com/amal1762/chat-app/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Name     string `json : "name"`
	Email    string `json : "email"`
	Password string `json : "password"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SetupRoutes(r fiber.Router) {
	r.Post("/login", Login)
	r.Post("/register", Register)
	r.Post("/refresh", RefreshToken)
	r.Post("/logout", Logout)
	r.Get("/me", AuthMiddleware, Me)

}

var JwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func Login(c *fiber.Ctx) error {
	var input LoginRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if input.Email == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "All fields are required"})
	}

	var user models.User

	db.DB.First(&user, "email = ?", input.Email)

	if user.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Could not map user"})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	access_token, err := generateAccessToken(user)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create token"})

	}
	RegisteredClaims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, RegisteredClaims).SignedString(JwtKey)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create refresh token"})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   false, // Set to true in production (needs HTTPS)
		SameSite: "Strict",
	})

	user.RefreshToken = refreshToken
	db.DB.Save(&user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": access_token,
	})
}

func Register(c *fiber.Ctx) error {
	var input RegisterInput
	// Placeholder register handler

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if input.Name == "" || input.Email == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "All fields are required"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not hash password"})

	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create user"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered"})
}

func Me(c *fiber.Ctx) error {
	var u models.User

	user := c.Locals("user").(*Claims)
	if err := db.DB.First(&u, user.UserID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{
		"id":    u.ID,
		"name":  u.Name,
		"email": u.Email,
	})
}

func RefreshToken(c *fiber.Ctx) error {

	// Step 1: Get the refresh token from the request body or cookie
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing refresh token"})
	}

	// Step 2: Parse and validate the JWT
	token, err := jwt.ParseWithClaims(refreshToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token"})
	}

	// Step 3: Find the user in DB with this refresh token
	var user models.User
	if err := db.DB.First(&user, "refresh_token = ?", refreshToken).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Refresh token not found"})
	}

	// Step 4: Generate a new access token

	access_token, err := generateAccessToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create access token"})
	}

	// Step 5: Return the new access token
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"access_token": access_token})
}

func Logout(c *fiber.Ctx) error {

	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing refresh token"})
	}
	var user models.User
	if err := db.DB.First(&user, "refresh_token = ?", refreshToken).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Refresh token not found"})

	}

	db.DB.Model(&user).Update("refresh_token", "")
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now(),
		HTTPOnly: true,
		Secure:   false, // Set to true in production (needs HTTPS)
		SameSite: "Strict",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User logged out"})

}

func generateAccessToken(user models.User) (string, error) {
	claims := &Claims{
		UserID: int(user.ID),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return accessToken.SignedString(JwtKey)

}
