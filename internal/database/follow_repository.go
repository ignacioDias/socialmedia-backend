package database

import (
	"context"
	"errors"
	"fmt"
	"socialnet/internal/models"

	"github.com/jmoiron/sqlx"
)

var ErrFollowNotFound = errors.New("follow not found")

type FollowRepository interface {
	CreateFollow(ctx context.Context, follow *models.Follow) error
}

type followRepository struct {
	db *sqlx.DB
}

func NewFollowRepository(db *sqlx.DB) FollowRepository {
	return &followRepository{db: db}
}

func (fr *followRepository) CreateFollow(ctx context.Context, follow *models.Follow) error {
	query := `INSERT INTO follows (following_id, follower_id) VALUES ($1, $2)`
	_, err := fr.db.ExecContext(ctx, query, follow.FollowingID, follow.FollowerID)
	return err
}

func (fr *followRepository) GetFollowersFromUserID(ctx context.Context, userID int64, limit, offset int) ([]models.User, error) {
	return fr.getUsersFromFollowRelation(ctx, userID, limit, offset, "f.follower_id = u.user_id", "f.following_id")
}

func (fr *followRepository) GetFollowingFromUserID(ctx context.Context, userID int64, limit, offset int) ([]models.User, error) {
	return fr.getUsersFromFollowRelation(ctx, userID, limit, offset, "f.following_id = u.user_id", "f.follower_id")
}

func (fr *followRepository) getUsersFromFollowRelation(ctx context.Context, userID int64, limit, offset int, joinCondition, whereColumn string) ([]models.User, error) {
	query := fmt.Sprintf(`
		SELECT u.profile_picture_path, u.banner_path, u.user_id, u.username, u.email, u.hashed_password
		FROM users u
		INNER JOIN follows f ON %s
		WHERE %s = $1
		LIMIT $2 OFFSET $3`, joinCondition, whereColumn)

	var users []models.User
	if err := fr.db.SelectContext(ctx, &users, query, userID, limit, offset); err != nil {
		return nil, err
	}
	return users, nil
}

func (fr *followRepository) GetCantFollowersFromUserID(ctx context.Context, userID int64) (int, error) {
	query := `SELECT COUNT(follower_id) FROM follows WHERE following_id = $1`
	return fr.countFromQuery(ctx, query, userID)
}

func (fr *followRepository) GetCantFollowingsFromUserID(ctx context.Context, userID int64) (int, error) {
	query := `SELECT COUNT(following_id) FROM follows WHERE follower_id = $1`
	return fr.countFromQuery(ctx, query, userID)
}

func (fr *followRepository) countFromQuery(ctx context.Context, query string, id int64) (int, error) {
	var cant int
	err := fr.db.GetContext(ctx, &cant, query, id)
	return cant, err
}

func (fr *followRepository) DeleteFollow(ctx context.Context, follow *models.Follow) error {
	query := `DELETE FROM follows WHERE follower_id = $1 AND following_id = $2`
	result, err := fr.db.ExecContext(ctx, query, follow.FollowerID, follow.FollowingID)
	return CheckErrResult(result, err, ErrFollowNotFound)
}
