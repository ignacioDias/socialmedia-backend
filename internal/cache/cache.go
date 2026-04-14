package cache

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
	ctx    context.Context
}

func NewCache(addr string) *Cache {
	redis := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	return &Cache{
		client: redis,
		ctx:    context.Background(),
	}
}

func GenerateCacheKey(parts ...string) string {
	h := sha256.New()
	for _, p := range parts {
		h.Write([]byte(p))
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (c *Cache) Get(key string, dest any) error {
	val, err := c.client.Get(c.ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (c *Cache) Set(key string, value any, expiration time.Duration) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(c.ctx, key, val, expiration).Err()
}

func (c *Cache) Delete(key string) error {
	return c.client.Del(c.ctx, key).Err()
}

func (c *Cache) HealthCheck() error {
	return c.client.Ping(c.ctx).Err()
}

func (c *Cache) Close() error {
	return c.client.Close()
}
