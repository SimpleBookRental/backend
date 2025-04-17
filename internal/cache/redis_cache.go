package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisCache wraps a Redis client
type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
	ctx    context.Context
}

// NewRedisCache creates a new RedisCache instance
func NewRedisCache(addr string, ttlSeconds int) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisCache{
		client: rdb,
		ttl:    time.Duration(ttlSeconds) * time.Second,
		ctx:    context.Background(),
	}
}

// Get gets a value from cache by key and unmarshals into dest
func (r *RedisCache) Get(key string, dest interface{}) (bool, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return false, nil // cache miss
	} else if err != nil {
		return false, err
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Set sets a value in cache with TTL
func (r *RedisCache) Set(key string, value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(r.ctx, key, bytes, r.ttl).Err()
}

// Delete deletes a key from cache
func (r *RedisCache) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}
