package models

import "time"

type Chat struct {
	ChatID        int64     `db:"chat_id" json:"chatId"`
	User1ID       int64     `db:"user1_id" json:"user1Id"`
	User2ID       int64     `db:"user2_id" json:"user2Id"`
	LastMessageAt time.Time `db:"last_message_at" json:"lastMessageAt"`
}

type Message struct {
	MessageID int64     `db:"message_id" json:"messageId"`
	ChatID    int64     `db:"chat_id" json:"chatId"`
	SenderID  int64     `db:"sender_id" json:"senderId"`
	ImagePath string    `db:"image_path" json:"imagePath"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

func NewChat(userAID, userBID int64) *Chat {
	if userAID > userBID {
		userAID, userBID = userBID, userAID
	}
	return &Chat{User1ID: userAID, User2ID: userBID}
}

func NewMessage(chatID, senderID int64, content, imagePath string) Message {
	return Message{
		ChatID:    chatID,
		SenderID:  senderID,
		Content:   content,
		ImagePath: imagePath,
		CreatedAt: time.Now().UTC(),
	}
}
