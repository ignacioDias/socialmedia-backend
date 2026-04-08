package models

import "time"

type Post struct {
	PostID    int64     `db:"post_id" json:"postId"`
	UserID    int64     `db:"user_id" json:"userId"`
	Title     string    `db:"title" json:"title"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	ImagePath string    `db:"image_path" json:"imagePath"`
}

func NewPost(userID int64, title string, imagePath, content string) *Post {
	return &Post{
		UserID:    userID,
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
		ImagePath: imagePath,
	}
}
