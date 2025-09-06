package app

import "yordamchi-dev-bot/internal/middleware"

// MetricsProvider combines metrics and cache middleware for command access
type MetricsProvider struct {
	metricsMiddleware *middleware.MetricsMiddleware
	cachingMiddleware *middleware.CachingMiddleware
}

// NewMetricsProvider creates a new metrics provider
func NewMetricsProvider(metricsMiddleware *middleware.MetricsMiddleware, cachingMiddleware *middleware.CachingMiddleware) *MetricsProvider {
	return &MetricsProvider{
		metricsMiddleware: metricsMiddleware,
		cachingMiddleware: cachingMiddleware,
	}
}

// GetMetrics returns performance metrics
func (mp *MetricsProvider) GetMetrics() map[string]interface{} {
	return mp.metricsMiddleware.GetMetrics()
}

// GetCacheStats returns cache statistics
func (mp *MetricsProvider) GetCacheStats() map[string]interface{} {
	return mp.cachingMiddleware.GetCacheStats()
}