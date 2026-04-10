package handler

import (
	"encoding/json"
	"net/http"
	"socialnet/internal/cache"
	"socialnet/internal/database"
	"socialnet/internal/middleware"
	"socialnet/internal/models"
	"socialnet/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

type UserRegisterRequest struct {
	ProfilePicturePath string `json:"profilePicturePath"`
	BannerPath         string `json:"bannerPath"`
	Email              string `json:"email"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	DocumentNumber     string `json:"documentNumber"`
	ProfilePicture     string `json:"profilePicture"`
}

func NewUserHandler(userRepo database.UserRepository, sessionRepo database.SessionRepository, cache *cache.Cache) *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(userRepo, sessionRepo, cache),
	}
}

func (uh *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var userRequest UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}
	user := models.NewUser(userRequest.Username, userRequest.Email, userRequest.Password, userRequest.ProfilePicture, userRequest.BannerPath)
	if err := uh.userService.CreateUser(r.Context(), user); err != nil {

	}

}

func (uh *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginReq service.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "wrong input for login request", http.StatusBadRequest)
		return
	}
	session, err := uh.userService.Login(r.Context(), &loginReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // false only in localhost dev
		SameSite: http.SameSiteStrictMode,
		Expires:  session.ExpiresAt,
	})
	w.WriteHeader(http.StatusOK)
}
func (uh *UserHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if err := uh.userService.Logout(r.Context(), userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})
	w.WriteHeader(http.StatusOK)
}
