package database

import (
	"context"
	"errors"
	"socialnet/internal/models"

	"github.com/jmoiron/sqlx"
)

var ErrRepostNotFound = errors.New("Repost not found")

type RepostRepository interface {
	CreateRepost(ctx context.Context, repost *models.Repost) error
	GetRepostsCountFromPost(ctx context.Context, postID int64) (int64, error)
	DeleteRepost(ctx context.Context, repost *models.Repost) error
}

type repostRepository struct {
	db *sqlx.DB
}

func NewRepostRepository(db *sqlx.DB) RepostRepository {
	return &repostRepository{
		db: db,
	}
}

func (rr *repostRepository) CreateRepost(ctx context.Context, repost *models.Repost) error {
	query := `INSERT INTO reposts (post_id, user_id, content) VALUES ($1, $2, $3)`
	_, err := rr.db.ExecContext(ctx, query, repost.PostID, repost.UserID, repost.Content)
	return err
}

func (rr *repostRepository) GetRepostsCountFromPost(ctx context.Context, postID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM reposts WHERE post_id = $1`
	var cantReposts int64
	err := rr.db.GetContext(ctx, &cantReposts, query, postID)
	return cantReposts, err
}

func (rr *repostRepository) DeleteRepost(ctx context.Context, repost *models.Repost) error {
	query := `DELETE FROM reposts WHERE user_id = $1 AND post_id = $2`
	result, err := rr.db.ExecContext(ctx, query, repost.UserID, repost.PostID)
	return CheckErrResult(result, err, ErrRepostNotFound)
}
