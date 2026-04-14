package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"socialnet/internal/cache"
	"socialnet/internal/database"
	"socialnet/internal/models"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterRequest struct {
	Email          string `json:"email"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	ProfilePicture string `json:"profilePicture"`
}

type UpdateInfoReq struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
}
type UpdatePasswordReq struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type UserService interface {
	CreateUser(ctx context.Context, userReq *UserRegisterRequest) (*models.User, error)
	Logout(ctx context.Context, userID int64) error
	Login(ctx context.Context, loginReq *UserLoginRequest) (*models.Session, error)
	GetUserByID(ctx context.Context, userID int64) (*models.User, error)
	DeleteUserByID(ctx context.Context, userID int64) error
	UpdatePassword(ctx context.Context, userID int64, updatePassReq *UpdatePasswordReq) error
	UpdateInfo(ctx context.Context, userID int64, updateInfoReq *UpdateInfoReq) (*models.User, error)
}

type userServiceImpl struct {
	userRepo    database.UserRepository
	sessionRepo database.SessionRepository
	cache       *cache.Cache
}

const DEFAULT_PROFILE_PICTURE = ""
const DEFAULT_BANNER = ""

var ErrInvalidEmail = errors.New("invalid email format")
var ErrInvalidUsername = errors.New("invalid username format")

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func NewUserService(userRepo database.UserRepository, sessionRepo database.SessionRepository, cache *cache.Cache) UserService {
	return &userServiceImpl{
		userRepo:    userRepo,
		cache:       cache,
		sessionRepo: sessionRepo,
	}
}

func (s *userServiceImpl) CreateUser(ctx context.Context, userReq *UserRegisterRequest) (*models.User, error) {
	if !isValidEmail(userReq.Email) {
		return nil, errors.New("invalid email format")
	}
	if !isValidPassword(userReq.Password) {
		return nil, errors.New("invalid password: must be 8/32 characters long and include at least one number, one special character, one uppercase letter, and one lowercase letter")
	}
	if userReq.Username == "" {
		return nil, errors.New("invalid username")
	}
	password, err := hashPassword(userReq.Password)
	if err != nil {
		return nil, err
	}
	user := models.NewUser(userReq.Username, userReq.Email, string(password), DEFAULT_PROFILE_PICTURE, DEFAULT_BANNER)
	return user, s.userRepo.CreateUser(ctx, user)
}
func (s *userServiceImpl) Login(ctx context.Context, loginReq *UserLoginRequest) (*models.Session, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, loginReq.Email)
	if err != nil {
		return nil, err
	}
	if !comparePasswords(user.HashedPassword, loginReq.Password) {
		return nil, errors.New("invalid credentials")
	}
	session := models.NewSession(user.UserID)
	if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}
	return session, nil
}

func (s *userServiceImpl) Logout(ctx context.Context, userID int64) error {
	return s.sessionRepo.DeleteSessionsByUserID(ctx, userID)
}

func (s *userServiceImpl) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	var user *models.User
	key := fmt.Sprintf("user:%d", userID)
	if err := s.cache.Get(key, &user); err == nil {
		return user, nil
	}

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Set(key, user, 2*time.Hour)
	return user, nil
}

func (s *userServiceImpl) UpdatePassword(ctx context.Context, userID int64, updatePassReq *UpdatePasswordReq) error {
	if !isValidPassword(updatePassReq.NewPassword) {
		return errors.New("invalid password: must be 8/32 characters long and include at least one number, one special character, one uppercase letter, and one lowercase letter")
	}
	if ok, err := s.isPasswordCorrect(ctx, userID, updatePassReq.OldPassword); err != nil {
		return err
	} else if !ok {
		return errors.New("invalid password")
	}
	hashedPassword, err := hashPassword(updatePassReq.NewPassword)
	if err != nil {
		return err
	}

	if err := s.userRepo.UpdateUserHashedPassword(ctx, string(hashedPassword), userID); err != nil {
		return err
	}
	key := fmt.Sprintf("user:%d", userID)
	_ = s.cache.Delete(key)
	return nil
}

func (s *userServiceImpl) UpdateInfo(ctx context.Context, userID int64, updateInfoReq *UpdateInfoReq) (*models.User, error) {
	if updateInfoReq.Email != nil && !isValidEmail(*updateInfoReq.Email) {
		return nil, ErrInvalidEmail
	}
	if updateInfoReq.Username != nil && *updateInfoReq.Username == "" {
		return nil, ErrInvalidUsername
	}
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if updateInfoReq.Email != nil {
		user.Email = *updateInfoReq.Email
	}
	if updateInfoReq.Username != nil {
		user.Username = *updateInfoReq.Username
	}
	updatedUser, err := s.userRepo.UpdateUserInfo(ctx, user)
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("user:%d", userID)
	_ = s.cache.Delete(key)
	return updatedUser, nil
}

func (s *userServiceImpl) DeleteUserByID(ctx context.Context, userID int64) error {
	if err := s.userRepo.DeleteUserByID(ctx, userID); err != nil {
		return err
	}
	key := fmt.Sprintf("user:%d", userID)
	_ = s.cache.Delete(key)
	return nil
}

func (s *userServiceImpl) isPasswordCorrect(ctx context.Context, userID int64, oldPassword string) (bool, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return false, err
	}
	return comparePasswords(user.HashedPassword, oldPassword), nil
}

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 12)
}
func isValidPassword(password string) bool {
	if len(password) < 8 || len(password) >= 32 {
		return false
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool

	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func comparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
