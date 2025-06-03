package main

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/amal1762/chat-app/internal/router"
	"github.com/amal1762/chat-app/internal/db"
	"fmt"
)

func main() {
	if err := db.Init(); err != nil {
		log.Fatalf("DB connection failed: %v", err)
    }
    fmt.Println("Connected to PostgreSQL!")
	app := fiber.New()
	router.SetupRoutes(app)


	log.Fatal(app.Listen(":8000"))
}