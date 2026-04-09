package models

import "time"

type Follow struct {
	FollowingID int64     `db:"following_id" json:"followingId"`
	FollowerID  int64     `db:"follower_id" json:"followerId"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
}

func NewFollow(followingID, followerID int64) *Follow {
	return &Follow{
		FollowingID: followingID,
		FollowerID:  followerID,
		CreatedAt:   time.Now(),
	}
}
