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
	"fmt"
	"sync"
	"time"
)

// HealthStatus represents the health status of a component
type HealthStatus string

const (
	Healthy   HealthStatus = "healthy"
	Unhealthy HealthStatus = "unhealthy"
	Degraded  HealthStatus = "degraded"
	Unknown   HealthStatus = "unknown"
)

// HealthCheck represents a health check
type HealthCheck struct {
	Name        string                 `json:"name"`
	Status      HealthStatus           `json:"status"`
	Message     string                 `json:"message,omitempty"`
	LastChecked time.Time              `json:"last_checked"`
	Duration    time.Duration          `json:"duration"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// HealthChecker defines the interface for health checks
type HealthChecker interface {
	Check(ctx context.Context) *HealthCheck
	GetName() string
}

// HealthMonitor monitors the health of various components
type HealthMonitor struct {
	checkers map[string]HealthChecker
	mutex    sync.RWMutex
	metrics  *MetricsCollector
}

// NewHealthMonitor creates a new health monitor
func NewHealthMonitor(metrics *MetricsCollector) *HealthMonitor {
	return &HealthMonitor{
		checkers: make(map[string]HealthChecker),
		metrics:  metrics,
	}
}

// RegisterHealthCheck registers a health checker
func (hm *HealthMonitor) RegisterHealthCheck(checker HealthChecker) {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	hm.checkers[checker.GetName()] = checker
}

// UnregisterHealthCheck unregisters a health checker
func (hm *HealthMonitor) UnregisterHealthCheck(name string) {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	delete(hm.checkers, name)
}

// CheckHealth performs all registered health checks
func (hm *HealthMonitor) CheckHealth(ctx context.Context) map[string]*HealthCheck {
	hm.mutex.RLock()
	checkers := make(map[string]HealthChecker)
	for name, checker := range hm.checkers {
		checkers[name] = checker
	}
	hm.mutex.RUnlock()

	results := make(map[string]*HealthCheck)
	for name, checker := range checkers {
		check := checker.Check(ctx)
		results[name] = check

		// Record metrics
		hm.metrics.RecordTimer("health_check_duration", check.Duration, map[string]string{
			"check":  name,
			"status": string(check.Status),
		})

		hm.metrics.IncrementCounter("health_check_total", map[string]string{
			"check":  name,
			"status": string(check.Status),
		})
	}

	return results
}

// GetOverallHealth returns the overall health status
func (hm *HealthMonitor) GetOverallHealth(ctx context.Context) *HealthCheck {
	checks := hm.CheckHealth(ctx)

	if len(checks) == 0 {
		return &HealthCheck{
			Name:        "overall",
			Status:      Unknown,
			Message:     "No health checks registered",
			LastChecked: time.Now(),
		}
	}

	healthyCount := 0
	unhealthyCount := 0
	degradedCount := 0
	unknownCount := 0

	var totalDuration time.Duration
	for _, check := range checks {
		totalDuration += check.Duration
		switch check.Status {
		case Healthy:
			healthyCount++
		case Unhealthy:
			unhealthyCount++
		case Degraded:
			degradedCount++
		case Unknown:
			unknownCount++
		}
	}

	overallStatus := Healthy
	message := "All systems healthy"

	switch {
	case unhealthyCount > 0:
		overallStatus = Unhealthy
		message = fmt.Sprintf("%d unhealthy components", unhealthyCount)
	case degradedCount > 0:
		overallStatus = Degraded
		message = fmt.Sprintf("%d degraded components", degradedCount)
	case unknownCount > 0:
		overallStatus = Unknown
		message = fmt.Sprintf("%d unknown components", unknownCount)
	}

	return &HealthCheck{
		Name:        "overall",
		Status:      overallStatus,
		Message:     message,
		LastChecked: time.Now(),
		Duration:    totalDuration,
		Metadata: map[string]interface{}{
			"total_checks":    len(checks),
			"healthy_count":   healthyCount,
			"unhealthy_count": unhealthyCount,
			"degraded_count":  degradedCount,
			"unknown_count":   unknownCount,
		},
	}
}

// GetHealthSummary returns a summary of health checks
func (hm *HealthMonitor) GetHealthSummary(ctx context.Context) map[string]interface{} {
	checks := hm.CheckHealth(ctx)
	overall := hm.GetOverallHealth(ctx)

	summary := map[string]interface{}{
		"overall_status":  overall.Status,
		"overall_message": overall.Message,
		"last_checked":    overall.LastChecked,
		"total_checks":    len(checks),
		"checks":          checks,
	}

	// Count by status
	statusCounts := make(map[HealthStatus]int)
	for _, check := range checks {
		statusCounts[check.Status]++
	}
	summary["status_counts"] = statusCounts

	return summary
}

// SimpleHealthCheck is a simple implementation of HealthChecker
type SimpleHealthCheck struct {
	name    string
	checkFn func(ctx context.Context) (HealthStatus, string, error)
}

// NewSimpleHealthCheck creates a new simple health check
func NewSimpleHealthCheck(name string, checkFn func(ctx context.Context) (HealthStatus, string, error)) *SimpleHealthCheck {
	return &SimpleHealthCheck{
		name:    name,
		checkFn: checkFn,
	}
}

// Check performs the health check
func (shc *SimpleHealthCheck) Check(ctx context.Context) *HealthCheck {
	start := time.Now()
	status, message, err := shc.checkFn(ctx)
	duration := time.Since(start)

	if err != nil {
		status = Unhealthy
		message = err.Error()
	}

	return &HealthCheck{
		Name:        shc.name,
		Status:      status,
		Message:     message,
		LastChecked: time.Now(),
		Duration:    duration,
	}
}

// GetName returns the health check name
func (shc *SimpleHealthCheck) GetName() string {
	return shc.name
}

// AlwaysHealthyCheck is a health check that always returns healthy
type AlwaysHealthyCheck struct {
	name string
}

// NewAlwaysHealthyCheck creates a new always healthy check
func NewAlwaysHealthyCheck(name string) *AlwaysHealthyCheck {
	return &AlwaysHealthyCheck{name: name}
}

// Check always returns healthy
func (ahc *AlwaysHealthyCheck) Check(ctx context.Context) *HealthCheck {
	return &HealthCheck{
		Name:        ahc.name,
		Status:      Healthy,
		Message:     "Component is healthy",
		LastChecked: time.Now(),
		Duration:    0,
	}
}

// GetName returns the health check name
func (ahc *AlwaysHealthyCheck) GetName() string {
	return ahc.name
}
