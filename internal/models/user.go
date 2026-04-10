package models

type User struct {
	ProfilePicturePath string `db:"profile_picture_path" json:"profilePicturePath"`
	BannerPath         string `db:"banner_path" json:"bannerPath"`
	UserID             int64  `db:"user_id" json:"userId"`
	Username           string `db:"username" json:"username"`
	Email              string `db:"email" json:"email"`
	Password           string `json:"-"`
	HashedPassword     string `db:"hashed_password" json:"-"`
}

func NewUser(username string, email string, password string, profilePic, banner string) *User {
	return &User{
		Username:           username,
		ProfilePicturePath: profilePic,
		BannerPath:         banner,
		Email:              email,
		Password:           password,
	}
}
