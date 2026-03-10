package middleware

import (
	"log"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// PerformanceMetrics holds performance statistics
type PerformanceMetrics struct {
	mu sync.RWMutex

	// Response time tracking
	TotalRequests    int64   `json:"total_requests"`
	SlowRequests     int64   `json:"slow_requests"` // > 200ms
	AverageResponse  float64 `json:"average_response_ms"`
	MaxResponse      float64 `json:"max_response_ms"`
	MinResponse      float64 `json:"min_response_ms"`
	P95Response      float64 `json:"p95_response_ms"`
	P99Response      float64 `json:"p99_response_ms"`

	// Time-based tracking
	LastReset        time.Time `json:"last_reset"`
	LastSlowRequest  time.Time `json:"last_slow_request"`

	// Response time histogram (for percentile calculation)
	responseTimes []float64
	slowThreshold time.Duration
}

// NewPerformanceMetrics creates a new PerformanceMetrics instance
func NewPerformanceMetrics() *PerformanceMetrics {
	return &PerformanceMetrics{
		LastReset:       time.Now(),
		responseTimes:   make([]float64, 0, 1000),
		slowThreshold:   200 * time.Millisecond,
		MinResponse:     math.MaxFloat64,
	}
}

// RecordResponse records a response time
func (pm *PerformanceMetrics) RecordResponse(duration time.Duration) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.TotalRequests++
	durationMs := float64(duration.Nanoseconds()) / 1e6

	// Track slow requests
	if duration > pm.slowThreshold {
		pm.SlowRequests++
		pm.LastSlowRequest = time.Now()
	}

	// Update min/max
	if durationMs < pm.MinResponse {
		pm.MinResponse = durationMs
	}
	if durationMs > pm.MaxResponse {
		pm.MaxResponse = durationMs
	}

	// Update average (running average to avoid storing all values)
	pm.AverageResponse = pm.AverageResponse + (durationMs-pm.AverageResponse)/float64(pm.TotalRequests)

	// Store recent response times for percentile calculation (keep last 1000)
	pm.responseTimes = append(pm.responseTimes, durationMs)
	if len(pm.responseTimes) > 1000 {
		pm.responseTimes = pm.responseTimes[1:]
	}

	// Calculate percentiles
	pm.calculatePercentiles()
}

// calculatePercentiles calculates P95 and P99 response times
func (pm *PerformanceMetrics) calculatePercentiles() {
	if len(pm.responseTimes) == 0 {
		return
	}

	// Sort response times
	sorted := make([]float64, len(pm.responseTimes))
	copy(sorted, pm.responseTimes)
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	// Calculate P95
	p95Index := int(float64(len(sorted)) * 0.95)
	if p95Index >= len(sorted) {
		p95Index = len(sorted) - 1
	}
	pm.P95Response = sorted[p95Index]

	// Calculate P99
	p99Index := int(float64(len(sorted)) * 0.99)
	if p99Index >= len(sorted) {
		p99Index = len(sorted) - 1
	}
	pm.P99Response = sorted[p99Index]
}

// Reset resets all metrics
func (pm *PerformanceMetrics) Reset() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.TotalRequests = 0
	pm.SlowRequests = 0
	pm.AverageResponse = 0
	pm.MaxResponse = 0
	pm.MinResponse = math.MaxFloat64
	pm.P95Response = 0
	pm.P99Response = 0
	pm.LastReset = time.Now()
	pm.responseTimes = pm.responseTimes[:0]
}

// GetMetrics returns a copy of current metrics
func (pm *PerformanceMetrics) GetMetrics() map[string]interface{} {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return map[string]interface{}{
		"total_requests":     pm.TotalRequests,
		"slow_requests":      pm.SlowRequests,
		"slow_percentage":   float64(pm.SlowRequests) / float64(pm.TotalRequests) * 100,
		"average_response_ms": pm.AverageResponse,
		"max_response_ms":     pm.MaxResponse,
		"min_response_ms":     pm.MinResponse,
		"p95_response_ms":     pm.P95Response,
		"p99_response_ms":     pm.P99Response,
		"last_reset":         pm.LastReset,
		"last_slow_request":  pm.LastSlowRequest,
		"target_response_ms":  200,
		"meets_target":       pm.AverageResponse < 200,
	}
}

// PerformanceMiddleware creates middleware for performance monitoring
func PerformanceMiddleware(metrics *PerformanceMetrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Record metrics
		metrics.RecordResponse(duration)

		// Add response time header
		c.Header("X-Response-Time", duration.String())

		// Log slow requests
		if duration > metrics.slowThreshold {
			log.Printf("Slow request: %s %s took %v", c.Request.Method, c.Request.URL.Path, duration)
		}
	}
}

// PerformanceMonitorHandler returns a handler for performance metrics endpoint
func PerformanceMonitorHandler(metrics *PerformanceMetrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    metrics.GetMetrics(),
		})
	}
}

// PerformanceBudgetMiddleware creates middleware that alerts when performance budget is exceeded
func PerformanceBudgetMiddleware(metrics *PerformanceMetrics, budgetMs time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		if duration > budgetMs {
			// Add warning header
			c.Header("X-Performance-Warning", "Response time exceeded budget")
			c.Header("X-Response-Time", duration.String())
			c.Header("X-Performance-Budget", budgetMs.String())
		}
	}
}

// Global metrics instance
var globalMetrics = NewPerformanceMetrics()

// GetGlobalMetrics returns the global metrics instance
func GetGlobalMetrics() *PerformanceMetrics {
	return globalMetrics
}
