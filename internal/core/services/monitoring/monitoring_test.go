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
	"testing"
	"time"
)

func TestMetricsCollector_NewCollector(t *testing.T) {
	collector := NewMetricsCollector()

	if collector == nil {
		t.Fatal("Expected collector to be created")
	}

	if collector.started {
		t.Error("Expected collector to not be started initially")
	}

	if len(collector.metrics) != 0 {
		t.Error("Expected empty metrics initially")
	}
}

func TestMetricsCollector_IncrementCounter(t *testing.T) {
	collector := NewMetricsCollector()

	// Test incrementing a new counter
	collector.IncrementCounter("test_counter", map[string]string{"label": "value"})

	metric, exists := collector.GetMetric("test_counter", map[string]string{"label": "value"})
	if !exists {
		t.Fatal("Expected metric to exist")
	}

	if metric.Value != 1 {
		t.Errorf("Expected counter value to be 1, got %f", metric.Value)
	}

	if metric.Type != Counter {
		t.Errorf("Expected metric type to be Counter, got %s", metric.Type)
	}

	// Test incrementing existing counter
	collector.IncrementCounter("test_counter", map[string]string{"label": "value"})

	metric, exists = collector.GetMetric("test_counter", map[string]string{"label": "value"})
	if !exists {
		t.Fatal("Expected metric to still exist")
	}

	if metric.Value != 2 {
		t.Errorf("Expected counter value to be 2, got %f", metric.Value)
	}
}

func TestMetricsCollector_SetGauge(t *testing.T) {
	collector := NewMetricsCollector()

	// Test setting a gauge
	collector.SetGauge("test_gauge", 42.5, map[string]string{"label": "value"})

	metric, exists := collector.GetMetric("test_gauge", map[string]string{"label": "value"})
	if !exists {
		t.Fatal("Expected metric to exist")
	}

	if metric.Value != 42.5 {
		t.Errorf("Expected gauge value to be 42.5, got %f", metric.Value)
	}

	if metric.Type != Gauge {
		t.Errorf("Expected metric type to be Gauge, got %s", metric.Type)
	}

	// Test updating gauge
	collector.SetGauge("test_gauge", 100.0, map[string]string{"label": "value"})

	metric, exists = collector.GetMetric("test_gauge", map[string]string{"label": "value"})
	if !exists {
		t.Fatal("Expected metric to still exist")
	}

	if metric.Value != 100.0 {
		t.Errorf("Expected gauge value to be 100.0, got %f", metric.Value)
	}
}

func TestMetricsCollector_RecordHistogram(t *testing.T) {
	collector := NewMetricsCollector()

	// Test recording histogram
	collector.RecordHistogram("test_histogram", 1.5, map[string]string{"label": "value"})

	metric, exists := collector.GetMetric("test_histogram", map[string]string{"label": "value"})
	if !exists {
		t.Fatal("Expected metric to exist")
	}

	if metric.Value != 1.5 {
		t.Errorf("Expected histogram value to be 1.5, got %f", metric.Value)
	}

	if metric.Type != Histogram {
		t.Errorf("Expected metric type to be Histogram, got %s", metric.Type)
	}
}

func TestMetricsCollector_RecordTimer(t *testing.T) {
	collector := NewMetricsCollector()

	// Test recording timer
	duration := 1500 * time.Millisecond
	collector.RecordTimer("test_timer", duration, map[string]string{"label": "value"})

	metric, exists := collector.GetMetric("test_timer", map[string]string{"label": "value"})
	if !exists {
		t.Fatal("Expected metric to exist")
	}

	expectedValue := 1.5 // 1500ms = 1.5 seconds
	if metric.Value != expectedValue {
		t.Errorf("Expected timer value to be %f, got %f", expectedValue, metric.Value)
	}

	if metric.Type != Histogram {
		t.Errorf("Expected metric type to be Histogram, got %s", metric.Type)
	}
}

func TestMetricsCollector_GetAllMetrics(t *testing.T) {
	collector := NewMetricsCollector()

	// Add some metrics
	collector.IncrementCounter("counter1", nil)
	collector.SetGauge("gauge1", 10.0, nil)
	collector.RecordHistogram("histogram1", 5.0, nil)

	allMetrics := collector.GetAllMetrics()

	if len(allMetrics) != 3 {
		t.Errorf("Expected 3 metrics, got %d", len(allMetrics))
	}
}

func TestMetricsCollector_GetMetricsByType(t *testing.T) {
	collector := NewMetricsCollector()

	// Add metrics of different types
	collector.IncrementCounter("counter1", nil)
	collector.IncrementCounter("counter2", nil)
	collector.SetGauge("gauge1", 10.0, nil)

	counters := collector.GetMetricsByType(Counter)
	gauges := collector.GetMetricsByType(Gauge)

	if len(counters) != 2 {
		t.Errorf("Expected 2 counters, got %d", len(counters))
	}

	if len(gauges) != 1 {
		t.Errorf("Expected 1 gauge, got %d", len(gauges))
	}
}

func TestMetricsCollector_GetStats(t *testing.T) {
	collector := NewMetricsCollector()

	// Add some metrics
	collector.IncrementCounter("counter1", nil)
	collector.SetGauge("gauge1", 10.0, nil)

	stats := collector.GetStats()

	if stats["total_metrics"] != 2 {
		t.Errorf("Expected 2 total metrics, got %v", stats["total_metrics"])
	}

	if stats["started"] != false {
		t.Error("Expected collector to not be started")
	}
}

func TestHealthMonitor_RegisterHealthCheck(t *testing.T) {
	monitor := NewHealthMonitor(NewMetricsCollector())

	check := NewAlwaysHealthyCheck("test_check")
	monitor.RegisterHealthCheck(check)

	// Check if registered
	checks := monitor.CheckHealth(context.Background())
	if len(checks) != 1 {
		t.Errorf("Expected 1 health check, got %d", len(checks))
	}

	if checks["test_check"] == nil {
		t.Fatal("Expected test_check to be registered")
	}

	if checks["test_check"].Status != Healthy {
		t.Errorf("Expected status to be Healthy, got %s", checks["test_check"].Status)
	}
}

func TestHealthMonitor_UnregisterHealthCheck(t *testing.T) {
	monitor := NewHealthMonitor(NewMetricsCollector())

	check := NewAlwaysHealthyCheck("test_check")
	monitor.RegisterHealthCheck(check)

	// Unregister
	monitor.UnregisterHealthCheck("test_check")

	// Check if unregistered
	checks := monitor.CheckHealth(context.Background())
	if len(checks) != 0 {
		t.Errorf("Expected 0 health checks, got %d", len(checks))
	}
}

func TestHealthMonitor_GetOverallHealth(t *testing.T) {
	monitor := NewHealthMonitor(NewMetricsCollector())

	// Add some health checks
	monitor.RegisterHealthCheck(NewAlwaysHealthyCheck("check1"))
	monitor.RegisterHealthCheck(NewAlwaysHealthyCheck("check2"))

	overall := monitor.GetOverallHealth(context.Background())

	if overall.Status != Healthy {
		t.Errorf("Expected overall status to be Healthy, got %s", overall.Status)
	}

	if overall.Name != "overall" {
		t.Errorf("Expected name to be 'overall', got %s", overall.Name)
	}
}

func TestHealthMonitor_GetHealthSummary(t *testing.T) {
	monitor := NewHealthMonitor(NewMetricsCollector())

	// Add some health checks
	monitor.RegisterHealthCheck(NewAlwaysHealthyCheck("check1"))
	monitor.RegisterHealthCheck(NewAlwaysHealthyCheck("check2"))

	summary := monitor.GetHealthSummary(context.Background())

	if summary["overall_status"] != Healthy {
		t.Errorf("Expected overall status to be Healthy, got %v", summary["overall_status"])
	}

	if summary["total_checks"] != 2 {
		t.Errorf("Expected 2 total checks, got %v", summary["total_checks"])
	}
}

func TestProfiler_GetRuntimeStats(t *testing.T) {
	profiler := NewProfiler(DefaultProfilerConfig(), NewMetricsCollector())

	stats := profiler.GetRuntimeStats()

	if stats["goroutines"] == nil {
		t.Error("Expected goroutines stat to be present")
	}

	if stats["memory_allocated"] == nil {
		t.Error("Expected memory_allocated stat to be present")
	}
}

func TestProfiler_ProfileFunction(t *testing.T) {
	profiler := NewProfiler(DefaultProfilerConfig(), NewMetricsCollector())

	// Test profiling a function
	err := profiler.ProfileFunction(context.Background(), "test_function", func() error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestMonitoringService_NewService(t *testing.T) {
	service := NewMonitoringService()

	if service == nil {
		t.Fatal("Expected service to be created")
	}

	if service.started {
		t.Error("Expected service to not be started initially")
	}

	if service.metrics == nil {
		t.Error("Expected metrics collector to be created")
	}

	if service.profiler == nil {
		t.Error("Expected profiler to be created")
	}

	if service.health == nil {
		t.Error("Expected health monitor to be created")
	}
}

func TestMonitoringService_StartStop(t *testing.T) {
	service := NewMonitoringService()

	// Start service
	service.Start()

	if !service.IsStarted() {
		t.Error("Expected service to be started")
	}

	// Stop service
	service.Stop()

	if service.IsStarted() {
		t.Error("Expected service to be stopped")
	}
}

func TestMonitoringService_RecordOperation(t *testing.T) {
	service := NewMonitoringService()
	service.Start()
	defer service.Stop()

	// Record successful operation
	service.RecordOperation("test_op", 100*time.Millisecond, true)

	// Record failed operation
	service.RecordOperation("test_op", 50*time.Millisecond, false)

	// Check metrics
	metrics := service.GetMetrics().GetAllMetrics()

	// Should have operation metrics
	found := false
	for _, metric := range metrics {
		if metric.Name == "operation_duration" || metric.Name == "operation_total" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected operation metrics to be recorded")
	}
}

func TestMonitoringService_RecordCacheHitMiss(t *testing.T) {
	service := NewMonitoringService()
	service.Start()
	defer service.Stop()

	// Record cache hit
	service.RecordCacheHit("test_cache")

	// Record cache miss
	service.RecordCacheMiss("test_cache")

	// Check metrics
	metrics := service.GetMetrics().GetAllMetrics()

	// Should have cache metrics
	found := false
	for _, metric := range metrics {
		if metric.Name == "cache_hit" || metric.Name == "cache_miss" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected cache metrics to be recorded")
	}
}

func TestMonitoringService_GetSystemStats(t *testing.T) {
	service := NewMonitoringService()
	service.Start()
	defer service.Stop()

	stats := service.GetSystemStats()

	if stats["runtime"] == nil {
		t.Error("Expected runtime stats to be present")
	}

	if stats["metrics"] == nil {
		t.Error("Expected metrics stats to be present")
	}

	if stats["health"] == nil {
		t.Error("Expected health stats to be present")
	}
}

func TestMonitoringService_GetPerformanceReport(t *testing.T) {
	service := NewMonitoringService()
	service.Start()
	defer service.Stop()

	report := service.GetPerformanceReport()

	if report["system"] == nil {
		t.Error("Expected system stats to be present")
	}

	if report["metrics"] == nil {
		t.Error("Expected metrics to be present")
	}

	if report["counters"] == nil {
		t.Error("Expected counters to be present")
	}

	if report["gauges"] == nil {
		t.Error("Expected gauges to be present")
	}

	if report["histograms"] == nil {
		t.Error("Expected histograms to be present")
	}

	if report["timers"] == nil {
		t.Error("Expected timers to be present")
	}
}
