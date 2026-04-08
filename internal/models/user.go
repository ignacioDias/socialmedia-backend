package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ProfilePicturePath string `db:"profile_picture_path" json:"profilePicturePath"`
	BannerPath         string `db:"banner_path" json:"bannerPath"`
	UserID             int64  `db:"user_id" json:"userId"`
	Username           string `db:"username" json:"username"`
	Email              string `db:"email" json:"email"`
	HashedPassword     string `db:"hashed_password" json:"hashedPassword"`
}

func NewUser(username string, email string, password string, profilePic, banner string) (*User, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}
	return &User{
		Username:           username,
		ProfilePicturePath: profilePic,
		BannerPath:         banner,
		Email:              email,
		HashedPassword:     hashedPassword,
	}, nil
}

func hashPassword(password string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(pass), err
}

func (u User) IsPasswordCorrect(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	if err == nil {
		return true, nil
	}
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}
	return false, err
}
