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
	"fmt"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

// MetricType represents the type of metric
type MetricType string

const (
	Counter   MetricType = "counter"
	Gauge     MetricType = "gauge"
	Histogram MetricType = "histogram"
	Timer     MetricType = "timer"
)

// Metric represents a performance metric
type Metric struct {
	Name        string                 `json:"name"`
	Type        MetricType             `json:"type"`
	Value       float64                `json:"value"`
	Labels      map[string]string      `json:"labels,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
	Description string                 `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// MetricsCollector collects and manages performance metrics
type MetricsCollector struct {
	metrics            map[string]*Metric
	mutex              sync.RWMutex
	started            bool
	stopCh             chan struct{}
	collectionInterval time.Duration
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		metrics:            make(map[string]*Metric),
		stopCh:             make(chan struct{}),
		collectionInterval: 30 * time.Second, // Default to 30 seconds, aligned with typical Prometheus scrape interval
	}
}

// SetCollectionInterval sets the interval for runtime metrics collection
func (mc *MetricsCollector) SetCollectionInterval(interval time.Duration) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.collectionInterval = interval
}

// Start starts the metrics collector
func (mc *MetricsCollector) Start() {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	if mc.started {
		return
	}

	mc.started = true
	go mc.collectSystemMetrics()
}

// Stop stops the metrics collector
func (mc *MetricsCollector) Stop() {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	if !mc.started {
		return
	}

	close(mc.stopCh)
	mc.started = false
}

// IncrementCounter increments a counter metric
func (mc *MetricsCollector) IncrementCounter(name string, labels map[string]string) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	key := mc.getMetricKey(name, labels)
	if metric, exists := mc.metrics[key]; exists {
		metric.Value++
		metric.Timestamp = time.Now()
	} else {
		mc.metrics[key] = &Metric{
			Name:      name,
			Type:      Counter,
			Value:     1,
			Labels:    labels,
			Timestamp: time.Now(),
		}
	}
}

// AddToCounter adds a custom value to a counter metric
func (mc *MetricsCollector) AddToCounter(name string, value float64, labels map[string]string) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	key := mc.getMetricKey(name, labels)
	if metric, exists := mc.metrics[key]; exists {
		metric.Value += value
		metric.Timestamp = time.Now()
	} else {
		mc.metrics[key] = &Metric{
			Name:      name,
			Type:      Counter,
			Value:     value,
			Labels:    labels,
			Timestamp: time.Now(),
		}
	}
}

// SetCounter sets a counter metric to an absolute value (for cumulative metrics)
func (mc *MetricsCollector) SetCounter(name string, value float64, labels map[string]string) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	key := mc.getMetricKey(name, labels)
	mc.metrics[key] = &Metric{
		Name:      name,
		Type:      Counter,
		Value:     value,
		Labels:    labels,
		Timestamp: time.Now(),
	}
}

// SetGauge sets a gauge metric value
func (mc *MetricsCollector) SetGauge(name string, value float64, labels map[string]string) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	key := mc.getMetricKey(name, labels)
	mc.metrics[key] = &Metric{
		Name:      name,
		Type:      Gauge,
		Value:     value,
		Labels:    labels,
		Timestamp: time.Now(),
	}
}

// RecordHistogram records a histogram value
func (mc *MetricsCollector) RecordHistogram(name string, value float64, labels map[string]string) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	key := mc.getMetricKey(name, labels)
	if metric, exists := mc.metrics[key]; exists {
		// For simplicity, we'll store the latest value
		// In a real implementation, you'd want to store buckets
		metric.Value = value
		metric.Timestamp = time.Now()
	} else {
		mc.metrics[key] = &Metric{
			Name:      name,
			Type:      Histogram,
			Value:     value,
			Labels:    labels,
			Timestamp: time.Now(),
		}
	}
}

// RecordTimer records a timer value
func (mc *MetricsCollector) RecordTimer(name string, duration time.Duration, labels map[string]string) {
	mc.RecordHistogram(name, float64(duration.Nanoseconds())/1e9, labels) // Convert to seconds
}

// GetMetric retrieves a metric by name and labels
func (mc *MetricsCollector) GetMetric(name string, labels map[string]string) (*Metric, bool) {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	key := mc.getMetricKey(name, labels)
	metric, exists := mc.metrics[key]
	return metric, exists
}

// GetAllMetrics returns all collected metrics
func (mc *MetricsCollector) GetAllMetrics() map[string]*Metric {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	// Return a copy to avoid race conditions
	result := make(map[string]*Metric)
	for key, metric := range mc.metrics {
		result[key] = &Metric{
			Name:        metric.Name,
			Type:        metric.Type,
			Value:       metric.Value,
			Labels:      metric.Labels,
			Timestamp:   metric.Timestamp,
			Description: metric.Description,
			Metadata:    metric.Metadata,
		}
	}
	return result
}

// GetMetricsByType returns metrics filtered by type
func (mc *MetricsCollector) GetMetricsByType(metricType MetricType) map[string]*Metric {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	result := make(map[string]*Metric)
	for key, metric := range mc.metrics {
		if metric.Type == metricType {
			result[key] = &Metric{
				Name:        metric.Name,
				Type:        metric.Type,
				Value:       metric.Value,
				Labels:      metric.Labels,
				Timestamp:   metric.Timestamp,
				Description: metric.Description,
				Metadata:    metric.Metadata,
			}
		}
	}
	return result
}

// ClearMetrics clears all metrics
func (mc *MetricsCollector) ClearMetrics() {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.metrics = make(map[string]*Metric)
}

// GetStats returns statistics about the metrics collector
func (mc *MetricsCollector) GetStats() map[string]interface{} {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	stats := map[string]interface{}{
		"total_metrics": len(mc.metrics),
		"started":       mc.started,
	}

	// Count metrics by type
	typeCounts := make(map[MetricType]int)
	for _, metric := range mc.metrics {
		typeCounts[metric.Type]++
	}
	stats["metrics_by_type"] = typeCounts

	return stats
}

// getMetricKey creates a unique key for a metric
func (mc *MetricsCollector) getMetricKey(name string, labels map[string]string) string {
	key := name
	if len(labels) > 0 {
		// Sort label keys for deterministic key generation
		keys := make([]string, 0, len(labels))
		for k := range labels {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			key += ":" + k + "=" + labels[k]
		}
	}
	return key
}

// collectSystemMetrics collects system-level metrics
func (mc *MetricsCollector) collectSystemMetrics() {
	mc.mutex.RLock()
	interval := mc.collectionInterval
	mc.mutex.RUnlock()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mc.collectRuntimeMetrics()
		case <-mc.stopCh:
			return
		}
	}
}

// collectRuntimeMetrics collects Go runtime metrics
func (mc *MetricsCollector) collectRuntimeMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Memory metrics
	mc.SetGauge("go_memory_alloc_bytes", float64(m.Alloc), nil)
	mc.SetGauge("go_memory_total_alloc_bytes", float64(m.TotalAlloc), nil)
	mc.SetGauge("go_memory_sys_bytes", float64(m.Sys), nil)
	mc.SetGauge("go_memory_heap_alloc_bytes", float64(m.HeapAlloc), nil)
	mc.SetGauge("go_memory_heap_sys_bytes", float64(m.HeapSys), nil)
	mc.SetGauge("go_memory_heap_idle_bytes", float64(m.HeapIdle), nil)
	mc.SetGauge("go_memory_heap_inuse_bytes", float64(m.HeapInuse), nil)
	mc.SetGauge("go_memory_heap_released_bytes", float64(m.HeapReleased), nil)
	mc.SetGauge("go_memory_heap_objects", float64(m.HeapObjects), nil)

	// GC metrics - use counters for cumulative totals
	mc.SetCounter("go_gc_cycles_total", float64(m.NumGC), nil)
	mc.SetCounter("go_gc_pause_total_ns", float64(m.PauseTotalNs), nil)

	// Goroutines
	mc.SetGauge("go_goroutines", float64(runtime.NumGoroutine()), nil)

	// CPU count
	mc.SetGauge("go_cpu_count", float64(runtime.NumCPU()), nil)
}

// ToPrometheusFormat exports metrics in Prometheus exposition format
func (mc *MetricsCollector) ToPrometheusFormat() string {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	var builder strings.Builder

	// Group metrics by name
	metricsByName := make(map[string][]*Metric)
	for _, metric := range mc.metrics {
		metricsByName[metric.Name] = append(metricsByName[metric.Name], metric)
	}

	// Sort metric names for consistent output
	names := make([]string, 0, len(metricsByName))
	for name := range metricsByName {
		names = append(names, name)
	}
	sort.Strings(names)

	// Write each metric
	for _, name := range names {
		metrics := metricsByName[name]
		if len(metrics) == 0 {
			continue
		}

		// Write HELP and TYPE comments
		if metrics[0].Description != "" {
			builder.WriteString(fmt.Sprintf("# HELP %s %s\n", name, metrics[0].Description))
		}
		metricType := metrics[0].Type
		prometheusType := mc.getPrometheusType(metricType)
		builder.WriteString(fmt.Sprintf("# TYPE %s %s\n", name, prometheusType))

		// Write metric values without timestamps (Prometheus uses scrape time)
		for _, metric := range metrics {
			if len(metric.Labels) > 0 {
				labelStr := mc.formatLabels(metric.Labels)
				builder.WriteString(fmt.Sprintf("%s{%s} %g\n", name, labelStr, metric.Value))
			} else {
				builder.WriteString(fmt.Sprintf("%s %g\n", name, metric.Value))
			}
		}

		builder.WriteString("\n")
	}

	return builder.String()
}

// getPrometheusType converts internal metric type to Prometheus type
func (mc *MetricsCollector) getPrometheusType(metricType MetricType) string {
	switch metricType {
	case Counter:
		return "counter"
	case Gauge:
		return "gauge"
	case Histogram:
		return "histogram"
	case Timer:
		return "histogram"
	default:
		return "untyped"
	}
}

// escapePrometheusLabelValue escapes special characters in label values according to Prometheus format
func escapePrometheusLabelValue(value string) string {
	// Escape backslashes first to avoid double-escaping
	value = strings.ReplaceAll(value, "\\", "\\\\")
	// Escape double quotes
	value = strings.ReplaceAll(value, "\"", "\\\"")
	// Escape newlines
	value = strings.ReplaceAll(value, "\n", "\\n")
	return value
}

// formatLabels formats labels for Prometheus exposition format
func (mc *MetricsCollector) formatLabels(labels map[string]string) string {
	if len(labels) == 0 {
		return ""
	}

	// Sort labels for consistent output
	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(labels))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%q", k, escapePrometheusLabelValue(labels[k])))
	}

	return strings.Join(parts, ",")
}
