package cache

import (
	"context"
	"encoding/json"
	"time"

	"go-user-service/internal/models"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
    Client *redis.Client
}

// NewRedisCache creates a new RedisCache instance
func NewRedisCache(addr, password string, db int) *RedisCache {
    client := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })

    return &RedisCache{Client: client}
}

// GetUser tries to retrieve a cached user by username. Returns (nil, nil) if not found.
func (r *RedisCache) GetUser(ctx context.Context, username string) (*models.User, error) {
    key := "user:" + username
    s, err := r.Client.Get(ctx, key).Result()
    if err == redis.Nil {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }

    var u models.User
    if err := json.Unmarshal([]byte(s), &u); err != nil {
        return nil, err
    }
    return &u, nil
}

// SetUser caches the user with given TTL
func (r *RedisCache) SetUser(ctx context.Context, user *models.User, ttl time.Duration) error {
    key := "user:" + user.Username
    b, err := json.Marshal(user)
    if err != nil {
        return err
    }
    return r.Client.Set(ctx, key, b, ttl).Err()
}

// DeleteUser removes a user cache entry
func (r *RedisCache) DeleteUser(ctx context.Context, username string) error {
    key := "user:" + username
    return r.Client.Del(ctx, key).Err()
}
