package models

import "time"

type Like struct {
	TargetID   int64      `db:"target_id" json:"targetId"`
	TargetType TargetType `db:"target_type" json:"targetType"`
	UserID     int64      `db:"user_id" json:"userId"`
	CreatedAt  time.Time  `db:"created_at" json:"createdAt"`
}

type Repost struct {
	PostID    int64     `db:"post_id" json:"postId"`
	UserID    int64     `db:"user_id" json:"userId"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	Content   string    `db:"content" json:"content"`
	ImagePath string    `db:"image_path" json:"imagePath"`
}

type Bookmark struct {
	PostID    int64     `db:"post_id" json:"postId"`
	UserID    int64     `db:"user_id" json:"userId"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

func NewLike(targetID int64, targetType TargetType, userID int64) *Like {
	return &Like{
		TargetID:   targetID,
		TargetType: targetType,
		UserID:     userID,
		CreatedAt:  time.Now().UTC(),
	}
}

func NewRepost(postID int64, userID int64, imagePath, content string) *Repost {
	return &Repost{
		PostID:    postID,
		UserID:    userID,
		ImagePath: imagePath,
		Content:   content,
		CreatedAt: time.Now().UTC(),
	}
}

func NewBookmark(postID int64, userID int64) *Bookmark {
	return &Bookmark{
		PostID:    postID,
		UserID:    userID,
		CreatedAt: time.Now().UTC(),
	}
}
