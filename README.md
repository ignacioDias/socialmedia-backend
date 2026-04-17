# SocialNet (Go REST API)

SocialNet is a modular REST API for a social platform built with Go. It includes authentication, user management, posting, comments, likes, bookmarks, follows, direct messages, and image uploads.

A simple clone of twitter.

[![bender](https://imgflip.com/s/meme/Bender.jpg)](https://imgflip.com/s/meme/Bender.jpg)


## Features

- Cookie-based sessions (`session_id`)
- User registration, login, logout, profile updates
- Follow and follower graph
- Posts and feed from followed users
- Nested comments
- Likes for posts and comments
- Bookmarks
- Direct chats/messages
- Image upload endpoint (JPEG/PNG, up to 10 MB)
- Redis-backed cache abstraction
- SQL repository layer with transactional table initialization
- Unit tests for repository and service layers

## Tech Stack

- Go 1.26+
- `net/http` + Go 1.22 route patterns
- PostgreSQL via `github.com/jmoiron/sqlx`
- Redis via `github.com/redis/go-redis/v9`
- Password hashing via `golang.org/x/crypto`
- Tests with `go-sqlmock`

## Project Layout

```
cmd/socialnet/            # intended app entrypoint (currently empty)
internal/cache/           # Redis cache client and helpers
internal/database/        # repositories + DB schema init
internal/handler/         # HTTP handlers
internal/middleware/      # auth and rate-limit middleware
internal/models/          # domain models
internal/router/          # route registration
internal/server/          # HTTP server lifecycle
uploads/images/           # uploaded files storage
```

## API Base Path

All endpoints are mounted under:

`/api/v1`

## Authentication

- Login sets an HTTP-only cookie named `session_id`.
- Protected routes require that cookie.
- Middleware resolves user identity from the session repository and injects user ID into request context.

## Rate Limiting

Applied on selected routes through middleware:

- Token bucket per IP
- Default bucket: 3 tokens max, refill rate `0.05` tokens/sec (about 3 requests/min sustained)

## Pagination

Routes that support pagination use query params:

- `limit` (default `20`, max `100`)
- `offset` (default `0`)

## Uploads

`POST /api/v1/upload`

- Protected route
- Multipart field name: `myFile`
- Max file size: 10 MB
- Allowed MIME types: `image/jpeg`, `image/png`
- Stored in `uploads/images/`
- Response contains created file path

## Endpoints

### Auth / Session

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `DELETE /api/v1/auth/logout` (auth)

### Users

- `GET /api/v1/users/me` (auth)
- `PUT /api/v1/users/me/info` (auth)
- `PUT /api/v1/users/me/password` (auth)
- `PUT /api/v1/users/me/profile_picture` (auth)
- `PUT /api/v1/users/me/banner` (auth)
- `DELETE /api/v1/users/me` (auth)
- `GET /api/v1/users/{user_id}`

### Following

- `POST /api/v1/following/follow/{user_id}` (auth)
- `DELETE /api/v1/following/unfollow/{user_id}` (auth)
- `GET /api/v1/users/{user_id}/followers`
- `GET /api/v1/users/{user_id}/followers/count`
- `GET /api/v1/users/{user_id}/following`
- `GET /api/v1/users/{user_id}/following/count`

### Posts / Feed

- `POST /api/v1/posts` (auth)
- `GET /api/v1/posts/{post_id}`
- `DELETE /api/v1/posts/{post_id}` (auth)
- `GET /api/v1/users/{user_id}/posts`
- `GET /api/v1/feed` (auth)

### Comments

- `POST /api/v1/posts/{post_id}/comments` (auth)
- `POST /api/v1/comments/{comment_id}/comments` (auth)
- `GET /api/v1/comments/{comment_id}`
- `GET /api/v1/posts/{post_id}/comments`
- `GET /api/v1/comments/{comment_id}/comments`
- `DELETE /api/v1/comments/{comment_id}` (auth)
- `GET /api/v1/users/me/comments` (auth)

### Likes

- `POST /api/v1/posts/{post_id}/like` (auth)
- `POST /api/v1/comments/{comment_id}/like` (auth)
- `DELETE /api/v1/posts/{post_id}/like` (auth)
- `DELETE /api/v1/comments/{comment_id}/like` (auth)
- `GET /api/v1/users/me/likes` (auth)
- `GET /api/v1/posts/{post_id}/likes/count`
- `GET /api/v1/comments/{comment_id}/likes/count`

### Bookmarks

- `POST /api/v1/posts/{post_id}/bookmark` (auth)
- `DELETE /api/v1/posts/{post_id}/bookmark` (auth)
- `GET /api/v1/users/me/bookmarks` (auth)

### Chats / Messages

- `POST /api/v1/chats` (auth)
- `GET /api/v1/chats/me` (auth)
- `GET /api/v1/chats/{chat_id}/messages` (auth)
- `POST /api/v1/chats/{chat_id}/messages` (auth)
- `DELETE /api/v1/chats/messages/{message_id}` (auth)

## Running Tests

```bash
go test ./...
```

## Build Check

```bash
go build ./...
```
