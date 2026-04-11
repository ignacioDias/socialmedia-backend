package handler

import (
	"net/http"
	"socialnet/internal/service"
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
	followingService *service.FollowingService
}

func NewFollowingHandler() FollowingHandler {
	return nil
}

func (fh *implFollowingHandler) FollowUser(w http.ResponseWriter, r *http.Request)       {}
func (fh *implFollowingHandler) UnfollowUser(w http.ResponseWriter, r *http.Request)     {}
func (fh *implFollowingHandler) GetFollowers(w http.ResponseWriter, r *http.Request)     {}
func (fh *implFollowingHandler) GetFollowing(w http.ResponseWriter, r *http.Request)     {}
func (fh *implFollowingHandler) GetCantFollowers(w http.ResponseWriter, r *http.Request) {}
func (fh *implFollowingHandler) GetCantFollowing(w http.ResponseWriter, r *http.Request) {}
