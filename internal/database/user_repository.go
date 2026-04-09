package database

import (
	"context"
	"database/sql"
	"errors"
	"socialnet/internal/models"

	"github.com/jmoiron/sqlx"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, userID int64) (*models.User, error)
	UpdateUserInfo(ctx context.Context, user *models.User) error
	UpdateUserProfilePicture(ctx context.Context, profilePicturePath string, userID int64) error
	UpdateUserBanner(ctx context.Context, bannerPath string, userID int64) error
	UpdateUserHashedPassword(ctx context.Context, hashedPassword string, userID int64) error
	DeleteUserByID(ctx context.Context, userID int64) error
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

func (ur *userRepository) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	query := `SELECT * FROM users WHERE user_id = $1`
	var user models.User
	if err := ur.db.GetContext(ctx, &user, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) UpdateUserInfo(ctx context.Context, user *models.User) error {
	query := `UPDATE users SET username = $1, email = $2 WHERE user_id = $3`
	result, err := ur.db.ExecContext(ctx, query, user.Username, user.Email, user.UserID)
	return CheckErrResult(result, err, ErrUserNotFound)
}

func (ur *userRepository) UpdateUserProfilePicture(ctx context.Context, profilePicturePath string, userID int64) error {
	query := `UPDATE users SET profile_picture_path = $1 WHERE user_id = $2`
	result, err := ur.db.ExecContext(ctx, query, profilePicturePath, userID)
	return CheckErrResult(result, err, ErrUserNotFound)
}

func (ur *userRepository) UpdateUserBanner(ctx context.Context, bannerPath string, userID int64) error {
	query := `UPDATE users SET banner_path = $1 WHERE user_id = $2`
	result, err := ur.db.ExecContext(ctx, query, bannerPath, userID)
	return CheckErrResult(result, err, ErrUserNotFound)
}

func (ur *userRepository) UpdateUserHashedPassword(ctx context.Context, hashedPassword string, userID int64) error {
	query := `UPDATE users SET hashed_password = $1 WHERE user_id = $2`
	result, err := ur.db.ExecContext(ctx, query, hashedPassword, userID)
	return CheckErrResult(result, err, ErrUserNotFound)
}

func (ur *userRepository) DeleteUserByID(ctx context.Context, userID int64) error {
	query := `DELETE FROM users WHERE user_id = $1`
	result, err := ur.db.ExecContext(ctx, query, userID)
	return CheckErrResult(result, err, ErrUserNotFound)
}
