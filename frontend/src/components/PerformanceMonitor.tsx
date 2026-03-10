import React, { useState, useEffect, useCallback } from 'react'

interface PerformanceMetrics {
  total_requests: number
  slow_requests: number
  slow_percentage: number
  average_response_ms: number
  max_response_ms: number
  min_response_ms: number
  p95_response_ms: number
  p99_response_ms: number
  last_reset: string
  last_slow_request: string
  target_response_ms: number
  meets_target: boolean
}

interface PageLoadMetrics {
  domContentLoaded: number
  load: number
  firstPaint: number
  firstContentfulPaint: number
  largestContentfulPaint: number
  timeToInteractive: number
  cumulativeLayoutShift: number
  firstInputDelay: number
}

interface PerformanceMonitorProps {
  apiMetricsEndpoint?: string
  showInProduction?: boolean
  refreshInterval?: number
  position?: 'top-left' | 'top-right' | 'bottom-left' | 'bottom-right'
}

/**
 * PerformanceMonitor - Real-time performance monitoring component
 * Displays API response times and page load metrics
 * 
 * Features:
 * - Real-time API response time tracking
 * - Page load performance metrics (Web Vitals)
 * - Performance budget alerts
 * - Visual indicators for performance status
 * 
 * @example
 * // Basic usage
 * <PerformanceMonitor />
 * 
 * @example
 * // Custom configuration
 * <PerformanceMonitor 
 *   apiMetricsEndpoint="/api/performance/metrics"
 *   refreshInterval={5000}
 *   position="bottom-right"
 * />
 */
export function PerformanceMonitor({
  apiMetricsEndpoint = '/api/performance/metrics',
  showInProduction = false,
  refreshInterval = 5000,
  position = 'bottom-right',
}: PerformanceMonitorProps) {
  const [isVisible, setIsVisible] = useState(false)
  const [apiMetrics, setApiMetrics] = useState<PerformanceMetrics | null>(null)
  const [pageMetrics, setPageMetrics] = useState<PageLoadMetrics | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // Check if we should show in production
  const shouldShow = process.env.NODE_ENV === 'development' || showInProduction

  // Position styles
  const positionStyles = {
    'top-left': { top: '20px', left: '20px' },
    'top-right': { top: '20px', right: '20px' },
    'bottom-left': { bottom: '20px', left: '20px' },
    'bottom-right': { bottom: '20px', right: '20px' },
  }

  // Fetch API performance metrics
  const fetchApiMetrics = useCallback(async () => {
    if (!apiMetricsEndpoint) return
    
    setIsLoading(true)
    try {
      const response = await fetch(apiMetricsEndpoint)
      if (!response.ok) throw new Error('Failed to fetch metrics')
      const data = await response.json()
      setApiMetrics(data.data || data)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load metrics')
    } finally {
      setIsLoading(false)
    }
  }, [apiMetricsEndpoint])

  // Collect page load metrics using Performance API
  const collectPageMetrics = useCallback(() => {
    if (typeof performance === 'undefined') return

    const timing = performance.timing
    const navigation = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming
    
    // Calculate various timing metrics
    const domContentLoaded = timing.domContentLoadedEventEnd - timing.navigationStart
    const load = timing.loadEventEnd - timing.navigationStart
    
    // Get paint timing
    const paintEntries = performance.getEntriesByType('paint')
    const firstPaint = paintEntries.find(e => e.name === 'first-paint')?.startTime || 0
    const firstContentfulPaint = paintEntries.find(e => e.name === 'first-contentful-paint')?.startTime || 0

    // Get LCP (Largest Contentful Paint)
    let largestContentfulPaint = 0
    const lcpEntries = performance.getEntriesByType('largest-contentful-paint')
    if (lcpEntries.length > 0) {
      largestContentfulPaint = lcpEntries[lcpEntries.length - 1].startTime
    }

    // Get CLS (Cumulative Layout Shift)
    let cumulativeLayoutShift = 0
    const clsEntries = performance.getEntriesByType('layout-shift')
    clsEntries.forEach(entry => {
      cumulativeLayoutShift += (entry as any).value
    })

    // Get FID (First Input Delay) - using FCP as proxy since FID is event-based
    const firstInputDelay = firstContentfulPaint

    // Time to Interactive (approximation)
    const timeToInteractive = domContentLoaded + (firstContentfulPaint - domContentLoaded) * 1.5

    setPageMetrics({
      domContentLoaded,
      load,
      firstPaint,
      firstContentfulPaint,
      largestContentfulPaint,
      timeToInteractive,
      cumulativeLayoutShift,
      firstInputDelay,
    })
  }, [])

  // Initial load
  useEffect(() => {
    if (!shouldShow) return

    collectPageMetrics()
    fetchApiMetrics()

    // Set up periodic refresh
    const interval = setInterval(fetchApiMetrics, refreshInterval)
    return () => clearInterval(interval)
  }, [shouldShow, refreshInterval, fetchApiMetrics, collectPageMetrics])

  // Keyboard shortcut to toggle visibility (Shift + P)
  useEffect(() => {
    const handleKeyPress = (e: KeyboardEvent) => {
      if (e.shiftKey && e.key === 'P') {
        setIsVisible(v => !v)
      }
    }

    window.addEventListener('keydown', handleKeyPress)
    return () => window.removeEventListener('keydown', handleKeyPress)
  }, [])

  // Get performance status color
  const getStatusColor = (value: number, target: number) => {
    if (value <= target * 0.5) return '#22c55e' // Green - Excellent
    if (value <= target) return '#eab308' // Yellow - Good
    if (value <= target * 1.5) return '#f97316' // Orange - Needs Improvement
    return '#ef4444' // Red - Poor
  }

  // Format milliseconds
  const formatMs = (ms: number) => {
    if (ms < 1) return `${ms.toFixed(2)}ms`
    if (ms < 1000) return `${ms.toFixed(0)}ms`
    return `${(ms / 1000).toFixed(2)}s`
  }

  // Calculate overall performance score
  const calculateScore = () => {
    if (!apiMetrics && !pageMetrics) return null

    let score = 100
    const penalties = []

    // API Response time penalty
    if (apiMetrics) {
      if (apiMetrics.average_response_ms > 200) {
        penalties.push(`API: ${apiMetrics.average_response_ms.toFixed(0)}ms`)
        score -= Math.min(30, (apiMetrics.average_response_ms - 200) / 10)
      }
    }

    // Page load penalties
    if (pageMetrics) {
      if (pageMetrics.largestContentfulPaint > 2500) {
        penalties.push(`LCP: ${pageMetrics.largestContentfulPaint.toFixed(0)}ms`)
        score -= 20
      }
      if (pageMetrics.firstContentfulPaint > 1800) {
        penalties.push(`FCP: ${pageMetrics.firstContentfulPaint.toFixed(0)}ms`)
        score -= 15
      }
      if (pageMetrics.cumulativeLayoutShift > 0.1) {
        penalties.push(`CLS: ${pageMetrics.cumulativeLayoutShift.toFixed(3)}`)
        score -= 15
      }
    }

    return {
      score: Math.max(0, Math.round(score)),
      penalties,
    }
  }

  const performanceScore = calculateScore()

  if (!shouldShow) return null

  return (
    <>
      {/* Toggle button */}
      <button
        onClick={() => setIsVisible(v => !v)}
        style={{
          position: 'fixed',
          ...positionStyles[position],
          zIndex: 9999,
          padding: '8px 12px',
          background: '#1a1a1a',
          color: '#fff',
          border: 'none',
          borderRadius: '6px',
          cursor: 'pointer',
          fontSize: '12px',
          fontFamily: 'system-ui, -apple-system, sans-serif',
          boxShadow: '0 2px 8px rgba(0,0,0,0.2)',
          transition: 'all 0.2s',
        }}
        title="Shift+P to toggle"
      >
        📊 Performance
      </button>

      {/* Monitor panel */}
      {isVisible && (
        <div
          style={{
            position: 'fixed',
            ...positionStyles[position],
            zIndex: 9998,
            marginTop: '50px',
            background: 'rgba(26, 26, 26, 0.95)',
            backdropFilter: 'blur(10px)',
            borderRadius: '12px',
            padding: '16px',
            color: '#fff',
            fontSize: '12px',
            fontFamily: 'system-ui, -apple-system, sans-serif',
            boxShadow: '0 4px 20px rgba(0,0,0,0.3)',
            maxWidth: '400px',
            width: '380px',
            maxHeight: '80vh',
            overflow: 'auto',
          }}
        >
          {/* Header */}
          <div style={{ 
            display: 'flex', 
            justifyContent: 'space-between', 
            alignItems: 'center',
            marginBottom: '12px',
            paddingBottom: '8px',
            borderBottom: '1px solid rgba(255,255,255,0.1)',
          }}>
            <h3 style={{ margin: 0, fontSize: '14px', fontWeight: 600 }}>
              Performance Monitor
            </h3>
            <button
              onClick={() => setIsVisible(false)}
              style={{
                background: 'transparent',
                border: 'none',
                color: '#fff',
                cursor: 'pointer',
                fontSize: '16px',
                padding: '0',
                lineHeight: 1,
              }}
            >
              ×
            </button>
          </div>

          {/* Performance Score */}
          {performanceScore && (
            <div style={{ 
              marginBottom: '16px', 
              padding: '12px', 
              background: 'rgba(255,255,255,0.05)',
              borderRadius: '8px',
            }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <span style={{ color: '#aaa' }}>Performance Score</span>
                <span style={{ 
                  fontSize: '24px', 
                  fontWeight: 'bold',
                  color: performanceScore.score >= 90 ? '#22c55e' : 
                         performanceScore.score >= 70 ? '#eab308' : '#ef4444',
                }}>
                  {performanceScore.score}
                </span>
              </div>
              {performanceScore.penalties.length > 0 && (
                <div style={{ marginTop: '8px', color: '#f97316', fontSize: '11px' }}>
                  {performanceScore.penalties.map((p, i) => (
                    <div key={i}>⚠️ {p}</div>
                  ))}
                </div>
              )}
            </div>
          )}

          {/* API Metrics */}
          {apiMetrics && (
            <div style={{ marginBottom: '16px' }}>
              <h4 style={{ 
                margin: '0 0 8px 0', 
                fontSize: '13px', 
                fontWeight: 600,
                color: '#60a5fa',
              }}>
                API Performance
              </h4>
              
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '8px' }}>
                <MetricItem 
                  label="Avg Response"
                  value={formatMs(apiMetrics.average_response_ms)}
                  target={apiMetrics.target_response_ms}
                  actual={apiMetrics.average_response_ms}
                />
                <MetricItem 
                  label="P95 Response"
                  value={formatMs(apiMetrics.p95_response_ms)}
                  target={apiMetrics.target_response_ms}
                  actual={apiMetrics.p95_response_ms}
                />
                <MetricItem 
                  label="P99 Response"
                  value={formatMs(apiMetrics.p99_response_ms)}
                  target={apiMetrics.target_response_ms * 2}
                  actual={apiMetrics.p99_response_ms}
                />
                <MetricItem 
                  label="Max Response"
                  value={formatMs(apiMetrics.max_response_ms)}
                />
              </div>

              <div style={{ 
                marginTop: '8px', 
                paddingTop: '8px', 
                borderTop: '1px solid rgba(255,255,255,0.1)',
                display: 'flex',
                justifyContent: 'space-between',
              }}>
                <span style={{ color: '#aaa' }}>Total Requests</span>
                <span>{apiMetrics.total_requests.toLocaleString()}</span>
              </div>
              
              <div style={{ 
                display: 'flex',
                justifyContent: 'space-between',
                marginTop: '4px',
              }}>
                <span style={{ color: '#aaa' }}>Slow Requests</span>
                <span style={{ color: apiMetrics.slow_percentage > 5 ? '#ef4444' : '#22c55e' }}>
                  {apiMetrics.slow_requests} ({apiMetrics.slow_percentage.toFixed(1)}%)
                </span>
              </div>

              {apiMetrics.meets_target && (
                <div style={{ 
                  marginTop: '8px', 
                  padding: '6px', 
                  background: 'rgba(34, 197, 94, 0.2)',
                  borderRadius: '4px',
                  color: '#22c55e',
                  textAlign: 'center',
                  fontSize: '11px',
                }}>
                  ✓ Meeting target (< {apiMetrics.target_response_ms}ms)
                </div>
              )}
            </div>
          )}

          {/* Page Load Metrics */}
          {pageMetrics && (
            <div style={{ marginBottom: '16px' }}>
              <h4 style={{ 
                margin: '0 0 8px 0', 
                fontSize: '13px', 
                fontWeight: 600,
                color: '#a78bfa',
              }}>
                Page Load Metrics
              </h4>
              
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '8px' }}>
                <MetricItem 
                  label="FCP"
                  value={formatMs(pageMetrics.firstContentfulPaint)}
                  target={1800}
                  actual={pageMetrics.firstContentfulPaint}
                />
                <MetricItem 
                  label="LCP"
                  value={formatMs(pageMetrics.largestContentfulPaint)}
                  target={2500}
                  actual={pageMetrics.largestContentfulPaint}
                />
                <MetricItem 
                  label="TTI"
                  value={formatMs(pageMetrics.timeToInteractive)}
                  target={3800}
                  actual={pageMetrics.timeToInteractive}
                />
                <MetricItem 
                  label="CLS"
                  value={pageMetrics.cumulativeLayoutShift.toFixed(3)}
                  target={0.1}
                  actual={pageMetrics.cumulativeLayoutShift}
                  isLowerBetter
                />
              </div>
            </div>
          )}

          {/* Loading/Error states */}
          {isLoading && (
            <div style={{ textAlign: 'center', padding: '12px', color: '#aaa' }}>
              Loading metrics...
            </div>
          )}
          
          {error && (
            <div style={{ 
              padding: '8px', 
              background: 'rgba(239, 68, 68, 0.2)', 
              borderRadius: '4px',
              color: '#ef4444',
              fontSize: '11px',
            }}>
              Error: {error}
            </div>
          )}

          {/* Footer */}
          <div style={{ 
            marginTop: '12px',
            paddingTop: '8px',
            borderTop: '1px solid rgba(255,255,255,0.1)',
            color: '#666',
            fontSize: '10px',
            textAlign: 'center',
          }}>
            Press Shift+P to toggle • Auto-refresh {refreshInterval / 1000}s
          </div>
        </div>
      )}
    </>
  )
}

// Helper component for metric items
function MetricItem({ 
  label, 
  value, 
  target, 
  actual,
  isLowerBetter = true,
}: { 
  label: string
  value: string
  target?: number
  actual?: number
  isLowerBetter?: boolean
}) {
  const statusColor = target && actual !== undefined 
    ? getStatusColor(actual, target, isLowerBetter)
    : '#666'

  return (
    <div style={{ 
      padding: '8px', 
      background: 'rgba(255,255,255,0.05)', 
      borderRadius: '6px',
    }}>
      <div style={{ color: '#aaa', fontSize: '10px', marginBottom: '4px' }}>{label}</div>
      <div style={{ 
        fontSize: '14px', 
        fontWeight: 600,
        color: statusColor,
      }}>
        {value}
      </div>
    </div>
  )
}

function getStatusColor(value: number, target: number, isLowerBetter: boolean) {
  const ratio = value / target
  if (isLowerBetter) {
    if (ratio <= 0.5) return '#22c55e'
    if (ratio <= 1) return '#eab308'
    if (ratio <= 1.5) return '#f97316'
    return '#ef4444'
  } else {
    if (ratio >= 1.5) return '#22c55e'
    if (ratio >= 1) return '#eab308'
    if (ratio >= 0.5) return '#f97316'
    return '#ef4444'
  }
}

export default PerformanceMonitor
