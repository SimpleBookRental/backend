package middleware

import (
	"bytes"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/SimpleBookRental/backend/internal/cache"
)

// CacheMiddleware returns a Gin middleware that caches GET responses in Redis
func CacheMiddleware(redisCache *cache.RedisCache, ttlSeconds int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.Next()
			return
		}

		// Create cache key based on userID, path and query
		cacheKey := "user:" + userID.(string) + ":cache:" + c.Request.URL.Path + "?" + c.Request.URL.RawQuery

		// Try to get cached response
		var cachedResponse []byte
		found, err := redisCache.Get(cacheKey, &cachedResponse)
		if err == nil && found {
			// log cache hit
			log.Println("[Cache HIT] Key:", cacheKey)
			c.Data(http.StatusOK, "application/json", cachedResponse)
			c.Abort()
			return
		} else if err != nil {
			log.Println("[Cache ERROR] Key:", cacheKey, "Error:", err.Error())
		} else {
			log.Println("[Cache MISS] Key:", cacheKey)
		}

		// Capture response body
		writer := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = writer

		c.Next()

		// Cache only if status is 200
		if c.Writer.Status() == http.StatusOK {
			err := redisCache.Set(cacheKey, writer.body.Bytes())
			if err != nil {
				log.Println("[Cache SET ERROR] Key:", cacheKey, "Error:", err.Error())
			} else {
				log.Println("[Cache SET] Key:", cacheKey)
			}
		}
	}
}

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
