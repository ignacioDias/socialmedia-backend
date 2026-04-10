package service

import "socialnet/internal/database"

type UserService struct {
	userRepo    *database.UserRepository
	sessionRepo *database.SessionRepository
}
