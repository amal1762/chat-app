package models

import "time"

type Message struct {
	ID         uint      `gorm:"primaryKey"`
	SenderID   uint      `json:"sender_id"`
	Sender     User      `gorm:"foreignKey:SenderID;references:ID" json:"sender"`
	ReceiverID uint      `json:"receiver_id"`
	Receiver   User      `gorm:"foreignKey:ReceiverID;references:ID" json:"receiver"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}
