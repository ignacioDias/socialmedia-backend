package database

import (
	"context"
	"socialnet/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (profile_picture_path, banner_path, username, email, hashed_password) VALUES ($1, $2, $3, $4, $5) RETURNING user_id`
	return ur.db.QueryRowContext(ctx, query, user.ProfilePicturePath, user.BannerPath, user.Username, user.Email, user.HashedPassword).Scan(&user.UserID)
}
