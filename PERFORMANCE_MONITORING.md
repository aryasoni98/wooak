# Performance Monitoring and Metrics Collection

This document describes the comprehensive performance monitoring and metrics collection system implemented in Wooak.

## Overview

Wooak now includes a built-in performance monitoring system with Prometheus-compatible metrics that tracks:
- File I/O operations (SSH config file operations)
- AI service requests and cache performance
- System resource usage (memory, CPU, goroutines)
- Application health checks

## Features

### 1. Prometheus-Compatible Metrics

All metrics are exported in Prometheus exposition format and can be accessed via HTTP endpoint.

**Metrics Endpoint**: `http://localhost:9090/metrics`

Example metrics output:
```
# TYPE http_requests_total counter
http_requests_total{method="GET",status="200"} 42 1728576000000

# TYPE go_memory_alloc_bytes gauge
go_memory_alloc_bytes 2048576 1728576000000

# TYPE operation_duration histogram
operation_duration{operation="ssh_config_load",status="success"} 0.025 1728576000000
```

### 2. File I/O Monitoring

All SSH config file operations are automatically monitored:

**Metrics Collected:**
- `ssh_config_load_total` - Total number of config file loads
- `ssh_config_save_total` - Total number of config file saves
- `ssh_config_backup_total` - Total number of backup operations
- `operation_duration{operation="ssh_config_load"}` - Time taken to load config
- `operation_duration{operation="ssh_config_save"}` - Time taken to save config
- `operation_duration{operation="ssh_config_backup"}` - Time taken to create backup
- `ssh_config_file_size_bytes{operation="load|save"}` - Size of config file
- `ssh_config_backup_count` - Number of backup files maintained

### 3. AI Service Monitoring

AI service requests and cache performance are tracked:

**Metrics Collected:**
- `ai_request_total{type="recommendation|search|security",status="success|failure"}` - Total AI requests
- `ai_request_duration_seconds{provider="ollama|openai",type="..."}` - AI request duration
- `ai_request_error_total{provider="...",type="..."}` - AI request errors
- `ai_tokens_used_total{provider="...",type="..."}` - Total tokens consumed
- `cache_hit{cache="ai_recommendation"}` - AI cache hits
- `cache_miss{cache="ai_recommendation"}` - AI cache misses
- `operation_duration{operation="ai_recommendation"}` - AI operation duration

### 4. Runtime Metrics

System and Go runtime metrics are automatically collected every 30 seconds:

**Metrics Collected:**
- `go_memory_alloc_bytes` - Bytes of allocated heap objects
- `go_memory_total_alloc_bytes` - Cumulative bytes allocated
- `go_memory_sys_bytes` - Total bytes obtained from system
- `go_memory_heap_alloc_bytes` - Bytes in heap
- `go_memory_heap_sys_bytes` - Heap system bytes
- `go_memory_heap_idle_bytes` - Idle heap bytes
- `go_memory_heap_inuse_bytes` - In-use heap bytes
- `go_memory_heap_released_bytes` - Released heap bytes
- `go_memory_heap_objects` - Number of heap objects
- `go_gc_cycles_total` - Number of GC cycles
- `go_gc_pause_total_ns` - Total GC pause time in nanoseconds
- `go_goroutines` - Number of goroutines
- `go_cpu_count` - Number of CPUs

### 5. Health Monitoring

Health check endpoint provides application status:

**Health Endpoint**: `http://localhost:9090/health`

Example response:
```json
{
  "status": "healthy",
  "timestamp": "2025-10-10T15:23:45Z"
}
```

Health statuses:
- `healthy` - All systems operational (HTTP 200)
- `degraded` - Some components degraded (HTTP 503)
- `unhealthy` - Critical failures (HTTP 503)

## Usage

### Starting the Application

When you start Wooak, the monitoring service automatically starts and begins collecting metrics:

```bash
./wooak
```

The metrics server listens on port 9090 by default.

### Accessing Metrics

**View metrics in Prometheus format:**
```bash
curl http://localhost:9090/metrics
```

**Check application health:**
```bash
curl http://localhost:9090/health
```

### Integration with Prometheus

Add this configuration to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'wooak'
    static_configs:
      - targets: ['localhost:9090']
    scrape_interval: 15s
```

### Grafana Dashboard

You can create Grafana dashboards using these metrics. Example queries:

**File I/O Performance:**
```promql
rate(ssh_config_load_total[5m])
histogram_quantile(0.95, rate(operation_duration_bucket{operation="ssh_config_load"}[5m]))
```

**Memory Usage Over Time:**
```promql
go_memory_heap_alloc_bytes
```

**AI Request Success Rate:**
```promql
rate(ai_request_total{status="success"}[5m]) / rate(ai_request_total[5m])
```

**Cache Hit Ratio:**
```promql
rate(cache_hit{cache="ai_recommendation"}[5m]) / 
(rate(cache_hit{cache="ai_recommendation"}[5m]) + rate(cache_miss{cache="ai_recommendation"}[5m]))
```

## Implementation Details

### Architecture

The monitoring system consists of:

1. **MetricsCollector** - Core metrics collection and storage
2. **MonitoringService** - Orchestrates all monitoring components
3. **HealthMonitor** - Health check management
4. **Profiler** - Runtime profiling capabilities

### Adding Custom Metrics

To add custom metrics in your code:

```go
// Get the monitoring service
monitoring := monitoringService.GetMetrics()

// Increment a counter
monitoring.IncrementCounter("my_metric_total", map[string]string{
    "label": "value",
})

// Set a gauge
monitoring.SetGauge("current_value", 42.0, map[string]string{
    "component": "my_component",
})

// Record operation duration
start := time.Now()
// ... do work ...
monitoringService.RecordOperation("my_operation", time.Since(start), true)
```

### Testing

Run monitoring tests:
```bash
go test ./internal/core/services/monitoring/... -v
```

All metrics and monitoring functionality is thoroughly tested with comprehensive unit tests.

## Performance Impact

The monitoring system is designed to have minimal performance impact:
- Metrics collection uses thread-safe concurrent maps
- Runtime metrics are collected every 30 seconds
- HTTP endpoints are served in a separate goroutine
- No blocking operations in metric recording

## Troubleshooting

**Port already in use:**
If port 9090 is already in use, you can modify the port in `cmd/main.go`:
```go
http.ListenAndServe(":9091", nil) // Change to any available port
```

**Metrics not showing:**
1. Ensure the application is running
2. Check that port 9090 is accessible
3. Verify no firewall is blocking the port
4. Check application logs for any errors

**High memory usage:**
The metrics collector stores all metrics in memory. For long-running applications, consider implementing metric rotation or using an external metrics database.

## Future Enhancements

Potential improvements:
- Metric persistence to disk
- Configurable metrics retention policy
- Additional metric types (summaries, etc.)
- Custom alerting based on metric thresholds
- Integration with additional monitoring systems (Datadog, New Relic, etc.)

## References

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
- [Go pprof Profiling](https://golang.org/pkg/runtime/pprof/)
