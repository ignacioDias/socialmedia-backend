package handler

import (
	"encoding/json"
	"net/http"
	"socialnet/internal/cache"
	"socialnet/internal/database"
	"socialnet/internal/middleware"
	"socialnet/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userRepo database.UserRepository, sessionRepo database.SessionRepository, cache *cache.Cache) *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(userRepo, sessionRepo, cache),
	}
}

func (uh *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	//TODO: Profile picture and banner
	var userRequest service.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}
	user, err := uh.userService.CreateUser(r.Context(), &userRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteResponseWithEncoder(w, user, http.StatusCreated)

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

func (uh *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized access", http.StatusUnauthorized)
		return
	}
	user, err := uh.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, user, http.StatusOK)
}

func (uh *UserHandler) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized access", http.StatusUnauthorized)
		return
	}
	var updateInfoReq service.UpdateInfoReq
	if err := json.NewDecoder(r.Body).Decode(&updateInfoReq); err != nil {
		http.Error(w, "invalid format request", http.StatusBadRequest)
		return
	}
	user, err := uh.userService.UpdateInfo(r.Context(), userID, &updateInfoReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteResponseWithEncoder(w, user, http.StatusOK)
}

func (uh *UserHandler) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized access", http.StatusUnauthorized)
		return
	}
	var updatePassReq service.UpdatePasswordReq
	if err := json.NewDecoder(r.Body).Decode(&updatePassReq); err != nil {
		http.Error(w, "invalid format request", http.StatusBadRequest)
		return
	}
	if err := uh.userService.UpdatePassword(r.Context(), userID, &updatePassReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (uh *UserHandler) DeleteCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized access", http.StatusUnauthorized)
		return
	}
	if err := uh.userService.DeleteUserByID(r.Context(), userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)
}
