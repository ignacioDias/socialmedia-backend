package models

import "time"

type Conversation struct {
	ConversationID int64     `db:"conversation_id" json:"conversationId"`
	CreatedAt      time.Time `db:"created_at" json:"createdAt"`
	LastMessageAt  time.Time `db:"last_message_at" json:"lastMessageAt"`
}

type ConversationMember struct {
	ConversationID int64     `db:"conversation_id" json:"conversationId"`
	UserID         int64     `db:"user_id" json:"userId"`
	JoinedAt       time.Time `db:"joined_at" json:"joinedAt"`
}

type Message struct {
	MessageID      int64     `db:"message_id" json:"messageId"`
	ConversationID int64     `db:"conversation_id" json:"conversationId"`
	SenderID       int64     `db:"sender_id" json:"senderId"`
	ImagePath      string    `db:"image_path" json:"imagePath"`
	Content        string    `db:"content" json:"content"`
	CreatedAt      time.Time `db:"created_at" json:"createdAt"`
}

func NewConversation() Conversation {
	now := time.Now().UTC()
	return Conversation{
		CreatedAt:     now,
		LastMessageAt: now,
	}
}

func NewConversationMember(conversationID, userID int64) ConversationMember {
	return ConversationMember{
		ConversationID: conversationID,
		UserID:         userID,
		JoinedAt:       time.Now().UTC(),
	}
}

func NewMessage(conversationID, senderID int64, content, imagePath string) Message {
	return Message{
		ConversationID: conversationID,
		SenderID:       senderID,
		Content:        content,
		ImagePath:      imagePath,
		CreatedAt:      time.Now().UTC(),
	}
}
