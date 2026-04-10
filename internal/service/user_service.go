package service

import (
	"context"
	"errors"
	"regexp"
	"socialnet/internal/cache"
	"socialnet/internal/database"
	"socialnet/internal/models"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) error
	Logout(ctx context.Context, userID int64) error
	Login(ctx context.Context, loginReq *UserLoginRequest) (*models.Session, error)
}

type userServiceImpl struct {
	userRepo    database.UserRepository
	sessionRepo database.SessionRepository
	cache       *cache.Cache
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func NewUserService(userRepo database.UserRepository, sessionRepo database.SessionRepository, cache *cache.Cache) UserService {
	return &userServiceImpl{
		userRepo:    userRepo,
		cache:       cache,
		sessionRepo: sessionRepo,
	}
}

func (s *userServiceImpl) CreateUser(ctx context.Context, user *models.User) error {
	if !isValidEmail(user.Email) {
		return errors.New("invalid email format")
	}
	if !isValidPassword(user.Password) {
		return errors.New("invalid password: must be 8–32 characters long and include at least one number, one special character, one uppercase letter, and one lowercase letter")
	}
	password, err := hashPassword(user.Password)
	if err != nil {
		return err
	}
	user.HashedPassword = string(password)

	return s.userRepo.CreateUser(ctx, user)
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
		return nil, err
	}
	return session, nil
}

func (s *userServiceImpl) Logout(ctx context.Context, userID int64) error {
	err := s.sessionRepo.DeleteSessionsByUserID(ctx, userID)
	return err
}

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
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

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 12)
}

func comparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
