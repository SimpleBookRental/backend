// Package cache defines the cache interface for pluggable cache implementations.
package cache

// Cache interface for cache implementations (e.g., RedisCache).
type Cache interface {
	Get(key string, dest interface{}) (bool, error)
	Set(key string, value interface{}) error
	Delete(key string) error
}
