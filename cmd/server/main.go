package main

import (
	"fmt"
	"log"

	"github.com/amal1762/chat-app/internal/db"
	"github.com/amal1762/chat-app/internal/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	if err := db.Init(); err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	fmt.Println("Connected to PostgreSQL!")
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://127.0.0.1:3000",
		AllowCredentials: true,
	}))

	router.SetupRoutes(app)

	log.Fatal(app.Listen(":8000"))
}
