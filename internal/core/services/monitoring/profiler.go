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
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

// ProfilerType represents the type of profiling
type ProfilerType string

const (
	CPUProfiler    ProfilerType = "cpu"
	MemoryProfiler ProfilerType = "memory"
	BlockProfiler  ProfilerType = "block"
	MutexProfiler  ProfilerType = "mutex"
)

// ProfilerConfig defines configuration for profiling
type ProfilerConfig struct {
	Enabled       bool          `json:"enabled"`
	Duration      time.Duration `json:"duration"`
	OutputDir     string        `json:"output_dir"`
	CPUProfile    bool          `json:"cpu_profile"`
	MemoryProfile bool          `json:"memory_profile"`
	BlockProfile  bool          `json:"block_profile"`
	MutexProfile  bool          `json:"mutex_profile"`
}

// DefaultProfilerConfig returns a default profiler configuration
func DefaultProfilerConfig() *ProfilerConfig {
	return &ProfilerConfig{
		Enabled:       true,
		Duration:      30 * time.Second,
		OutputDir:     "./profiles",
		CPUProfile:    true,
		MemoryProfile: true,
		BlockProfile:  false,
		MutexProfile:  false,
	}
}

// Profiler provides runtime profiling capabilities
type Profiler struct {
	config    *ProfilerConfig
	metrics   *MetricsCollector
	profiling bool
	mutex     sync.RWMutex
}

// NewProfiler creates a new profiler
func NewProfiler(config *ProfilerConfig, metrics *MetricsCollector) *Profiler {
	if config == nil {
		config = DefaultProfilerConfig()
	}

	return &Profiler{
		config:  config,
		metrics: metrics,
	}
}

// StartProfiling starts profiling with the specified type
func (p *Profiler) StartProfiling(profilerType ProfilerType) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.config.Enabled {
		return fmt.Errorf("profiling is disabled")
	}

	if p.profiling {
		return fmt.Errorf("profiling is already running")
	}

	switch profilerType {
	case CPUProfiler:
		return p.startCPUProfiling()
	case MemoryProfiler:
		return p.startMemoryProfiling()
	case BlockProfiler:
		return p.startBlockProfiling()
	case MutexProfiler:
		return p.startMutexProfiling()
	default:
		return fmt.Errorf("unsupported profiler type: %s", profilerType)
	}
}

// StopProfiling stops profiling and saves the results
func (p *Profiler) StopProfiling(profilerType ProfilerType) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.profiling {
		return fmt.Errorf("profiling is not running")
	}

	switch profilerType {
	case CPUProfiler:
		p.stopCPUProfiling()
		return nil
	case MemoryProfiler:
		p.stopMemoryProfiling()
		return nil
	case BlockProfiler:
		p.stopBlockProfiling()
		return nil
	case MutexProfiler:
		p.stopMutexProfiling()
		return nil
	default:
		return fmt.Errorf("unsupported profiler type: %s", profilerType)
	}
}

// ProfileFunction profiles a function execution
func (p *Profiler) ProfileFunction(ctx context.Context, name string, fn func() error) error {
	if !p.config.Enabled {
		return fn()
	}

	start := time.Now()
	defer func() {
		duration := time.Since(start)
		p.metrics.RecordTimer("function_execution_time", duration, map[string]string{
			"function": name,
		})
	}()

	return fn()
}

// GetRuntimeStats returns current runtime statistics
func (p *Profiler) GetRuntimeStats() map[string]interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return map[string]interface{}{
		"goroutines":           runtime.NumGoroutine(),
		"memory_allocated":     m.Alloc,
		"memory_total":         m.TotalAlloc,
		"memory_system":        m.Sys,
		"memory_heap":          m.HeapAlloc,
		"memory_heap_system":   m.HeapSys,
		"memory_heap_idle":     m.HeapIdle,
		"memory_heap_inuse":    m.HeapInuse,
		"memory_heap_released": m.HeapReleased,
		"memory_heap_objects":  m.HeapObjects,
		"gc_cycles":            m.NumGC,
		"gc_pause_total":       m.PauseTotalNs,
		"gc_pause_last":        m.PauseNs[(m.NumGC+255)%256],
	}
}

// startCPUProfiling starts CPU profiling
func (p *Profiler) startCPUProfiling() error {
	if !p.config.CPUProfile {
		return fmt.Errorf("CPU profiling is disabled")
	}

	// In a real implementation, you would start CPU profiling here
	p.profiling = true
	p.metrics.IncrementCounter("profiling_started", map[string]string{
		"type": "cpu",
	})

	return nil
}

// stopCPUProfiling stops CPU profiling
func (p *Profiler) stopCPUProfiling() {
	// In a real implementation, you would stop CPU profiling here
	p.profiling = false
	p.metrics.IncrementCounter("profiling_stopped", map[string]string{
		"type": "cpu",
	})
}

// startMemoryProfiling starts memory profiling
func (p *Profiler) startMemoryProfiling() error {
	if !p.config.MemoryProfile {
		return fmt.Errorf("memory profiling is disabled")
	}

	// In a real implementation, you would start memory profiling here
	p.profiling = true
	p.metrics.IncrementCounter("profiling_started", map[string]string{
		"type": "memory",
	})

	return nil
}

// stopMemoryProfiling stops memory profiling
func (p *Profiler) stopMemoryProfiling() {
	// In a real implementation, you would stop memory profiling here
	p.profiling = false
	p.metrics.IncrementCounter("profiling_stopped", map[string]string{
		"type": "memory",
	})
}

// startBlockProfiling starts block profiling
func (p *Profiler) startBlockProfiling() error {
	if !p.config.BlockProfile {
		return fmt.Errorf("block profiling is disabled")
	}

	runtime.SetBlockProfileRate(1)
	p.profiling = true
	p.metrics.IncrementCounter("profiling_started", map[string]string{
		"type": "block",
	})

	return nil
}

// stopBlockProfiling stops block profiling
func (p *Profiler) stopBlockProfiling() {
	runtime.SetBlockProfileRate(0)
	p.profiling = false
	p.metrics.IncrementCounter("profiling_stopped", map[string]string{
		"type": "block",
	})
}

// startMutexProfiling starts mutex profiling
func (p *Profiler) startMutexProfiling() error {
	if !p.config.MutexProfile {
		return fmt.Errorf("mutex profiling is disabled")
	}

	runtime.SetMutexProfileFraction(1)
	p.profiling = true
	p.metrics.IncrementCounter("profiling_started", map[string]string{
		"type": "mutex",
	})

	return nil
}

// stopMutexProfiling stops mutex profiling
func (p *Profiler) stopMutexProfiling() {
	runtime.SetMutexProfileFraction(0)
	p.profiling = false
	p.metrics.IncrementCounter("profiling_stopped", map[string]string{
		"type": "mutex",
	})
}

// GetProfileNames returns available profile names
func (p *Profiler) GetProfileNames() []string {
	// Return common profile names
	return []string{"goroutine", "heap", "allocs", "block", "mutex", "cpu"}
}

// GetProfile returns a profile by name
func (p *Profiler) GetProfile(name string) *pprof.Profile {
	return pprof.Lookup(name)
}
