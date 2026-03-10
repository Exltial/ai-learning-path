package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheConfig holds cache configuration
type CacheConfig struct {
	DefaultTTL time.Duration
	CourseTTL  time.Duration
	UserTTL    time.Duration
	KeyPrefix  string
}

// DefaultCacheConfig returns default cache configuration
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		DefaultTTL: 30 * time.Minute,
		CourseTTL:  15 * time.Minute,
		UserTTL:    10 * time.Minute,
		KeyPrefix:  "ai-learning:",
	}
}

// CacheService provides Redis caching for frequently accessed data
type CacheService struct {
	client *redis.Client
	config *CacheConfig
}

// NewCacheService creates a new CacheService
func NewCacheService(client *redis.Client, config *CacheConfig) *CacheService {
	if config == nil {
		config = DefaultCacheConfig()
	}
	return &CacheService{
		client: client,
		config: config,
	}
}

// CacheKey generates a cache key with prefix
func (s *CacheService) CacheKey(entityType, id string) string {
	return fmt.Sprintf("%s%s:%s", s.config.KeyPrefix, entityType, id)
}

// Get retrieves data from cache
func (s *CacheService) Get(ctx context.Context, key string, result interface{}) error {
	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, result)
}

// Set stores data in cache with default TTL
func (s *CacheService) Set(ctx context.Context, key string, value interface{}) error {
	return s.SetWithTTL(ctx, key, value, s.config.DefaultTTL)
}

// SetWithTTL stores data in cache with custom TTL
func (s *CacheService) SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	return s.client.Set(ctx, key, data, ttl).Err()
}

// Delete removes data from cache
func (s *CacheService) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}

// DeletePattern removes all keys matching a pattern
func (s *CacheService) DeletePattern(ctx context.Context, pattern string) error {
	iter := s.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := s.client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// GetOrSet retrieves from cache or sets if not exists
// Returns true if data was from cache, false if it was set
func (s *CacheService) GetOrSet(ctx context.Context, key string, result interface{}, fetchFunc func() (interface{}, error), ttl time.Duration) (bool, error) {
	// Try to get from cache
	if err := s.Get(ctx, key, result); err == nil {
		return true, nil
	}

	// Fetch from source
	data, err := fetchFunc()
	if err != nil {
		return false, err
	}

	// Set in cache
	if err := s.SetWithTTL(ctx, key, data, ttl); err != nil {
		return false, err
	}

	// Copy data to result
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal(dataBytes, result); err != nil {
		return false, err
	}

	return false, nil
}

// Course Cache Methods

// GetCourse retrieves a course from cache
func (s *CacheService) GetCourse(ctx context.Context, courseID string) (interface{}, error) {
	key := s.CacheKey("course", courseID)
	var course interface{}
	err := s.Get(ctx, key, &course)
	return course, err
}

// SetCourse stores a course in cache
func (s *CacheService) SetCourse(ctx context.Context, courseID string, course interface{}) error {
	key := s.CacheKey("course", courseID)
	return s.SetWithTTL(ctx, key, course, s.config.CourseTTL)
}

// InvalidateCourse removes a course from cache
func (s *CacheService) InvalidateCourse(ctx context.Context, courseID string) error {
	key := s.CacheKey("course", courseID)
	return s.Delete(ctx, key)
}

// User Cache Methods

// GetUser retrieves a user from cache
func (s *CacheService) GetUser(ctx context.Context, userID string) (interface{}, error) {
	key := s.CacheKey("user", userID)
	var user interface{}
	err := s.Get(ctx, key, &user)
	return user, err
}

// SetUser stores a user in cache
func (s *CacheService) SetUser(ctx context.Context, userID string, user interface{}) error {
	key := s.CacheKey("user", userID)
	return s.SetWithTTL(ctx, key, user, s.config.UserTTL)
}

// InvalidateUser removes a user from cache
func (s *CacheService) InvalidateUser(ctx context.Context, userID string) error {
	key := s.CacheKey("user", userID)
	return s.Delete(ctx, key)
}

// List Cache Methods

// GetCourseList retrieves a cached course list
func (s *CacheService) GetCourseList(ctx context.Context, category, difficulty string, page, limit int) (interface{}, error) {
	key := s.CacheKey("courses:list", fmt.Sprintf("%s:%s:%d:%d", category, difficulty, page, limit))
	var courses interface{}
	err := s.Get(ctx, key, &courses)
	return courses, err
}

// SetCourseList stores a course list in cache
func (s *CacheService) SetCourseList(ctx context.Context, category, difficulty string, page, limit int, courses interface{}) error {
	key := s.CacheKey("courses:list", fmt.Sprintf("%s:%s:%d:%d", category, difficulty, page, limit))
	return s.SetWithTTL(ctx, key, courses, 5*time.Minute)
}

// InvalidateCourseList removes all course list caches
func (s *CacheService) InvalidateCourseList(ctx context.Context) error {
	return s.DeletePattern(ctx, s.CacheKey("courses:list", "*"))
}

// Stats and Monitoring

// GetCacheStats returns cache statistics
func (s *CacheService) GetCacheStats(ctx context.Context) (map[string]interface{}, error) {
	info, err := s.client.Info(ctx, "stats").Result()
	if err != nil {
		return nil, err
	}

	// Parse Redis info (simplified)
	stats := map[string]interface{}{
		"connected": true,
		"info":      info,
	}

	return stats, nil
}

// WarmCache pre-warms cache for frequently accessed data
func (s *CacheService) WarmCache(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return s.SetWithTTL(ctx, key, value, ttl)
}

// HealthCheck checks if Redis connection is healthy
func (s *CacheService) HealthCheck(ctx context.Context) bool {
	_, err := s.client.Ping(ctx).Result()
	return err == nil
}
