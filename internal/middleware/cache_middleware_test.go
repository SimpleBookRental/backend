// Unit tests for CacheMiddleware using testify.
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/SimpleBookRental/backend/internal/cache"
	"github.com/SimpleBookRental/backend/internal/mocks"
)

func setupGinWithCache(cacheInstance cache.Cache) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	// Inject user_id into context BEFORE CacheMiddleware
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "1")
		c.Next()
	})
	r.Use(CacheMiddleware(cacheInstance, 60)) // Set TTL to 60 seconds
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": true})
	})
	// Add POST route for NonGetRequest test
	r.POST("/protected", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": true})
	})
	return r
}

func TestCacheMiddleware_CacheHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCache := mocks.NewMockCache(ctrl)

	// Expect Get to be called and return cached data
	mockCache.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		DoAndReturn(func(key string, dest interface{}) (bool, error) {
			if b, ok := dest.(*[]byte); ok {
				*b = []byte(`{"success": true}`)
			}
			return true, nil
		}).
		Times(1)

	router := setupGinWithCache(mockCache)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer validtoken")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestCacheMiddleware_CacheMiss(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCache := mocks.NewMockCache(ctrl)

	// Expect Get to be called and return nil (cache miss)
	mockCache.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		DoAndReturn(func(key string, dest interface{}) (bool, error) {
			return false, nil
		}).
		Times(1)
	// Expect Set to be called to store the response
	mockCache.EXPECT().
		Set(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	router := setupGinWithCache(mockCache)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer validtoken")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestCacheMiddleware_NonGetRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCache := mocks.NewMockCache(ctrl)

	// Should not call Get or Set for non-GET requests

	router := setupGinWithCache(mockCache)

	req, _ := http.NewRequest("POST", "/protected", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code) // Should proceed without caching
}
