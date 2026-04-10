package router

import (
	"net/http"
	"socialnet/internal/cache"
	"socialnet/internal/database"
	"socialnet/internal/handler"
	"socialnet/internal/middleware"
)

type Router struct {
	mux              *http.ServeMux
	userHandler      handler.UserHandler
	authenticationMw *middleware.AuthenticationMiddleware
	rateLimit        *middleware.RateLimitMiddleware
}

func NewRouter(db *database.Database, cache *cache.Cache) *Router {
	return &Router{
		mux:              http.NewServeMux(),
		authenticationMw: middleware.NewAuthenticationMiddleware(db.SessionRepo),
		rateLimit:        middleware.NewRateLimitMiddleware(),
	}
}

func (r *Router) SetupRoutes() *http.ServeMux {
	//session
	r.mux.HandleFunc("POST /api/v1/auth/register", r.rateLimit.RateLimit(r.userHandler.RegisterUser))
	r.mux.HandleFunc("POST /api/v1/auth/login", r.rateLimit.RateLimit(r.userHandler.LoginUser))
	r.mux.HandleFunc("DELETE /api/v1/auth/logout", r.authenticationMw.AuthenticationMiddleware(r.userHandler.LogoutUser))

	//user
	r.mux.HandleFunc("GET /api/v1/users/me", r.authenticationMw.AuthenticationMiddleware(r.userHandler.GetCurrentUser))
	r.mux.HandleFunc("PUT /api/v1/users/me", r.authenticationMw.AuthenticationMiddleware(r.userHandler.UpdateUser))
	r.mux.HandleFunc("DELETE /api/v1/users/me", r.authenticationMw.AuthenticationMiddleware(r.userHandler.DeleteMe))
	r.mux.HandleFunc("DELETE /api/v1/users/{user_id}", r.authCheck(r.userHandler.DeleteUser))

	return r.mux
}
