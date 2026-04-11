package handler

import (
	"net/http"
	"socialnet/internal/cache"
	"socialnet/internal/database"
	"socialnet/internal/middleware"
	"socialnet/internal/models"
	"socialnet/internal/service"
	"strconv"
)

type FollowingHandler interface {
	FollowUser(w http.ResponseWriter, r *http.Request)
	UnfollowUser(w http.ResponseWriter, r *http.Request)
	GetFollowers(w http.ResponseWriter, r *http.Request)
	GetFollowing(w http.ResponseWriter, r *http.Request)
	GetCantFollowers(w http.ResponseWriter, r *http.Request)
	GetCantFollowing(w http.ResponseWriter, r *http.Request)
}

type implFollowingHandler struct {
	followingService service.FollowingService
}

func NewFollowingHandler(followingRepo database.FollowRepository, cache *cache.Cache) FollowingHandler {
	return &implFollowingHandler{
		followingService: service.NewFollowingService(followingRepo, cache),
	}
}

func (fh *implFollowingHandler) getFollowFromRequest(w http.ResponseWriter, r *http.Request) (*models.Follow, bool) {
	followerID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return &models.Follow{}, false
	}
	idPath := r.PathValue("user_id")
	followingID, err := strconv.ParseInt(idPath, 10, 64)
	if err != nil {
		http.Error(w, "wrong id", http.StatusBadRequest)
		return &models.Follow{}, false
	}
	return models.NewFollow(followingID, followerID), true
}

func (fh *implFollowingHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	follow, ok := fh.getFollowFromRequest(w, r)
	if !ok {
		return
	}
	if err := fh.followingService.CreateFollow(r.Context(), follow); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, follow, http.StatusCreated)
}

func (fh *implFollowingHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	follow, ok := fh.getFollowFromRequest(w, r)
	if !ok {
		return
	}
	if err := fh.followingService.DeleteFollow(r.Context(), follow); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
func (fh *implFollowingHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	idPath := r.PathValue("user_id")
	userID, err := strconv.ParseInt(idPath, 10, 64)
	if err != nil {
		http.Error(w, "wrong id", http.StatusBadRequest)
		return
	}
	limit, err := ParseQueryInt(r, "limit", DEFAULT_LIMIT)
	if err != nil {
		http.Error(w, "invalid format for limit ", http.StatusBadRequest)
		return
	}
	offset, err := ParseQueryInt(r, "offset", DEFAULT_OFFSET)
	if err != nil {
		http.Error(w, "invalid format for offset ", http.StatusBadRequest)
		return
	}
	users, err := fh.followingService.GetFollowersFromID(r.Context(), userID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, users, http.StatusOK)
}

func (fh *implFollowingHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	idPath := r.PathValue("user_id")
	userID, err := strconv.ParseInt(idPath, 10, 64)
	if err != nil {
		http.Error(w, "wrong id", http.StatusBadRequest)
		return
	}
	limit, err := ParseQueryInt(r, "limit", DEFAULT_LIMIT)
	if err != nil {
		http.Error(w, "invalid format for limit ", http.StatusBadRequest)
		return
	}
	offset, err := ParseQueryInt(r, "offset", DEFAULT_OFFSET)
	if err != nil {
		http.Error(w, "invalid format for offset ", http.StatusBadRequest)
		return
	}
	users, err := fh.followingService.GetFollowingFromID(r.Context(), userID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, users, http.StatusOK)
}
func (fh *implFollowingHandler) GetCantFollowers(w http.ResponseWriter, r *http.Request) {
	idPath := r.PathValue("user_id")
	userID, err := strconv.ParseInt(idPath, 10, 64)
	if err != nil {
		http.Error(w, "wrong id", http.StatusBadRequest)
		return
	}
	followers, err := fh.followingService.GetCountFollowers(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, followers, http.StatusOK)
}
func (fh *implFollowingHandler) GetCantFollowing(w http.ResponseWriter, r *http.Request) {
	idPath := r.PathValue("user_id")
	userID, err := strconv.ParseInt(idPath, 10, 64)
	if err != nil {
		http.Error(w, "wrong id", http.StatusBadRequest)
		return
	}
	following, err := fh.followingService.GetCountFollowing(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteResponseWithEncoder(w, following, http.StatusOK)
}
