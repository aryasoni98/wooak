// Copyright 2025.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package monitoring

import (
	"context"
	"sync"
	"time"
)

// MonitoringService provides comprehensive monitoring capabilities
type MonitoringService struct {
	metrics  *MetricsCollector
	profiler *Profiler
	health   *HealthMonitor
	started  bool
	mutex    sync.RWMutex
}

// NewMonitoringService creates a new monitoring service
func NewMonitoringService() *MonitoringService {
	metrics := NewMetricsCollector()
	profiler := NewProfiler(DefaultProfilerConfig(), metrics)
	health := NewHealthMonitor(metrics)

	return &MonitoringService{
		metrics:  metrics,
		profiler: profiler,
		health:   health,
	}
}

// Start starts the monitoring service
func (ms *MonitoringService) Start() {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	if ms.started {
		return
	}

	ms.started = true
	ms.metrics.Start()

	// Register default health checks
	ms.registerDefaultHealthChecks()
}

// Stop stops the monitoring service
func (ms *MonitoringService) Stop() {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	if !ms.started {
		return
	}

	ms.metrics.Stop()
	ms.started = false
}

// GetMetrics returns the metrics collector
func (ms *MonitoringService) GetMetrics() *MetricsCollector {
	return ms.metrics
}

// GetProfiler returns the profiler
func (ms *MonitoringService) GetProfiler() *Profiler {
	return ms.profiler
}

// GetHealthMonitor returns the health monitor
func (ms *MonitoringService) GetHealthMonitor() *HealthMonitor {
	return ms.health
}

// RecordOperation records an operation metric
func (ms *MonitoringService) RecordOperation(operation string, duration time.Duration, success bool) {
	status := "success"
	if !success {
		status = "failure"
	}

	ms.metrics.RecordTimer("operation_duration", duration, map[string]string{
		"operation": operation,
		"status":    status,
	})

	ms.metrics.IncrementCounter("operation_total", map[string]string{
		"operation": operation,
		"status":    status,
	})
}

// RecordCacheHit records a cache hit
func (ms *MonitoringService) RecordCacheHit(cacheName string) {
	ms.metrics.IncrementCounter("cache_hit", map[string]string{
		"cache": cacheName,
	})
}

// RecordCacheMiss records a cache miss
func (ms *MonitoringService) RecordCacheMiss(cacheName string) {
	ms.metrics.IncrementCounter("cache_miss", map[string]string{
		"cache": cacheName,
	})
}

// RecordMemoryUsage records memory usage
func (ms *MonitoringService) RecordMemoryUsage(component string, bytes int64) {
	ms.metrics.SetGauge("memory_usage_bytes", float64(bytes), map[string]string{
		"component": component,
	})
}

// RecordGoroutineCount records the number of goroutines
func (ms *MonitoringService) RecordGoroutineCount() {
	// This would be called periodically to record goroutine count
	ms.metrics.SetGauge("goroutines", float64(0), map[string]string{}) // Placeholder
}

// GetSystemStats returns comprehensive system statistics
func (ms *MonitoringService) GetSystemStats() map[string]interface{} {
	stats := make(map[string]interface{})

	// Runtime stats
	stats["runtime"] = ms.profiler.GetRuntimeStats()

	// Metrics stats
	stats["metrics"] = ms.metrics.GetStats()

	// Health summary
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stats["health"] = ms.health.GetHealthSummary(ctx)

	return stats
}

// GetPerformanceReport returns a comprehensive performance report
func (ms *MonitoringService) GetPerformanceReport() map[string]interface{} {
	report := make(map[string]interface{})

	// System statistics
	report["system"] = ms.GetSystemStats()

	// All metrics
	report["metrics"] = ms.metrics.GetAllMetrics()

	// Performance metrics by type
	report["counters"] = ms.metrics.GetMetricsByType(Counter)
	report["gauges"] = ms.metrics.GetMetricsByType(Gauge)
	report["histograms"] = ms.metrics.GetMetricsByType(Histogram)
	report["timers"] = ms.metrics.GetMetricsByType(Timer)

	return report
}

// registerDefaultHealthChecks registers default health checks
func (ms *MonitoringService) registerDefaultHealthChecks() {
	// Application health check
	ms.health.RegisterHealthCheck(NewAlwaysHealthyCheck("application"))

	// Metrics health check
	ms.health.RegisterHealthCheck(NewSimpleHealthCheck("metrics", func(ctx context.Context) (HealthStatus, string, error) {
		stats := ms.metrics.GetStats()
		if stats["total_metrics"].(int) > 0 {
			return Healthy, "Metrics collection is working", nil
		}
		return Degraded, "No metrics collected yet", nil
	}))

	// Profiler health check
	ms.health.RegisterHealthCheck(NewSimpleHealthCheck("profiler", func(ctx context.Context) (HealthStatus, string, error) {
		// Check if profiler is available
		profiles := ms.profiler.GetProfileNames()
		if len(profiles) > 0 {
			return Healthy, "Profiler is available", nil
		}
		return Degraded, "Profiler not available", nil
	}))
}

// IsStarted returns whether the monitoring service is started
func (ms *MonitoringService) IsStarted() bool {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()
	return ms.started
}
