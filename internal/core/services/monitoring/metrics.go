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
	metrics map[string]*Metric
	mutex   sync.RWMutex
	started bool
	stopCh  chan struct{}
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		metrics: make(map[string]*Metric),
		stopCh:  make(chan struct{}),
	}
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
	for k, v := range labels {
		key += ":" + k + "=" + v
	}
	return key
}

// collectSystemMetrics collects system-level metrics
func (mc *MetricsCollector) collectSystemMetrics() {
	ticker := time.NewTicker(30 * time.Second)
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
	// This would collect actual runtime metrics
	// For now, we'll just record that we're collecting
	mc.IncrementCounter("system_metrics_collected", map[string]string{
		"source": "runtime",
	})
}
