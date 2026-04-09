package models

import "time"

type Comment struct {
	CommentID  int64      `db:"comment_id" json:"commentId"`
	TargetID   int64      `db:"target_id" json:"targetId"`
	TargetType TargetType `db:"target_type" json:"targetType"`
	UserID     int64      `db:"user_id" json:"userId"`
	Content    string     `db:"content" json:"content"`
	ImagePath  string     `db:"image_path" json:"imagePath"`
	CreatedAt  time.Time  `db:"created_at" json:"createdAt"`
}

func NewComment(targetID int64, targetType TargetType, userID int64, content string, imagePath string) *Comment {
	return &Comment{
		TargetID:   targetID,
		TargetType: targetType,
		UserID:     userID,
		Content:    content,
		CreatedAt:  time.Now().UTC(),
		ImagePath:  imagePath,
	}
}
