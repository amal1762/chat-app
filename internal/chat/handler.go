package chat

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/amal1762/chat-app/pkg/models"

	"github.com/amal1762/chat-app/internal/auth"
	"github.com/amal1762/chat-app/internal/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type IncomingMessage struct {
	ReceiverID uint   `json:"receiver_id"`
	Content    string `json:"content"`
}

func SetupRoutes(r fiber.Router) {
	r.Get("/ws", auth.AuthMiddleware, websocket.New(ChatWebSocket))
	r.Get("/users", auth.AuthMiddleware, FetchUser)
	r.Get("/messages/:withUserID", auth.AuthMiddleware, FetchMessage)
}

func FetchUser(c *fiber.Ctx) error {
	claims := c.Locals("user").(*auth.Claims)

	var users []models.User
	db.DB.Where("id != ?", claims.UserID).Find(&users)
	return c.JSON(fiber.Map{
		"users": users,
	})

}

func FetchMessage(c *fiber.Ctx) error {
	withID, err := strconv.Atoi(c.Params("withUserID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "withuser not provided"})
	}

	claims := c.Locals("user").(*auth.Claims)
	userID := claims.UserID // convert int to uint
	var messages []models.Message
	if err := db.DB.
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", userID, withID, withID, userID).
		Order("created_at ASC").
		Find(&messages).Error; err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch messages",
		})
	}

	return c.Status(200).JSON(messages)

}
func ChatWebSocket(c *websocket.Conn) {
	log.Println("WebSocket endpoint hit")
	ctx := c.Locals("user")
	claims := ctx.(*auth.Claims)
	userID := uint(claims.UserID)

	HubInstance.mu.Lock()
	HubInstance.clients[userID] = c
	HubInstance.mu.Unlock()

	defer func() {
		HubInstance.mu.Lock()
		delete(HubInstance.clients, userID)
		HubInstance.mu.Unlock()
		c.Close()
	}()

	for {
		_, msgBytes, err := c.ReadMessage()
		if err != nil {
			break
		}

		var incoming IncomingMessage
		if err := json.Unmarshal(msgBytes, &incoming); err != nil {
			continue
		}

		message := models.Message{
			SenderID:   userID,
			ReceiverID: incoming.ReceiverID,
			Content:    incoming.Content,
		}

		if err := db.DB.Create(&message).Error; err != nil {
			log.Println("Failed to save message:", err)
		}

		msgToSend, _ := json.Marshal(message)

		HubInstance.mu.RLock()
		receiverConn, receiverOnline := HubInstance.clients[incoming.ReceiverID]
		HubInstance.mu.RUnlock()

		if receiverOnline && incoming.ReceiverID != userID {
			if err := receiverConn.WriteMessage(websocket.TextMessage, msgToSend); err != nil {
				log.Println("Error sending message to receiver:", err)
			}
		}
	}
}
