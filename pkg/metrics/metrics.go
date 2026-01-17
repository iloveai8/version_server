// Package metrics 提供Prometheus监控指标
//
// 该包实现了：
// - HTTP请求指标（请求数、延迟、状态码）
// - Redis操作指标
// - 业务指标（版本配置查询等）
//
// 使用说明：
// 需要先安装依赖：go get github.com/prometheus/client_golang/prometheus
package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// ==================== HTTP指标 ====================

var (
	// HTTP请求总数
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "HTTP请求总数",
		},
		[]string{"method", "path", "status"},
	)

	// HTTP请求延迟（直方图）
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP请求延迟（秒）",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// HTTP请求进行中（Gauge）
	httpRequestsInProgress = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_in_progress",
			Help: "当前正在处理的HTTP请求数",
		},
		[]string{"method", "path"},
	)

	// HTTP响应大小（字节）
	httpResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP响应大小（字节）",
			Buckets: []float64{100, 1000, 10000, 100000, 1000000},
		},
		[]string{"method", "path"},
	)
)

// ==================== Redis指标 ====================

var (
	// Redis命令执行总数
	redisCommandsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "redis_commands_total",
			Help: "Redis命令执行总数",
		},
		[]string{"command", "status"}, // status: success/error
	)

	// Redis命令延迟
	redisCommandDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "redis_command_duration_seconds",
			Help:    "Redis命令执行延迟（秒）",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"command"},
	)

	// Redis连接数
	redisConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "redis_connections",
			Help: "Redis当前连接数",
		},
	)
)

// ==================== 业务指标 ====================

var (
	// 版本配置查询总数
	versionQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "version_queries_total",
			Help: "版本配置查询总数",
		},
		[]string{"env", "status"}, // env: dev/pro, status: hit/miss
	)

	// GM配置变更总数
	gmConfigChangesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gm_config_changes_total",
			Help: "GM配置变更总数",
		},
		[]string{"type"}, // type: toggle_gm, toggle_block, update_config
	)

	// IP白名单检查总数
	ipWhitelistChecksTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ip_whitelist_checks_total",
			Help: "IP白名单检查总数",
		},
		[]string{"status"}, // status: allowed/denied
	)

	// 活跃连接数
	activeConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "当前活跃连接数",
		},
	)
)

// ==================== HTTP指标记录函数 ====================

// RecordHTTPRequest 记录HTTP请求
//
// 参数:
//   method: HTTP方法
//   path: 请求路径
//   statusCode: HTTP状态码
//   duration: 请求耗时（秒）
//   responseSize: 响应大小（字节）
func RecordHTTPRequest(method, path string, statusCode int, duration time.Duration, responseSize int) {
	status := strconv.Itoa(statusCode)

	// 记录请求总数
	httpRequestsTotal.WithLabelValues(method, path, status).Inc()

	// 记录请求延迟
	httpRequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())

	// 记录响应大小
	if responseSize > 0 {
		httpResponseSize.WithLabelValues(method, path).Observe(float64(responseSize))
	}
}

// IncHTTPRequestsInProgress 增加进行中的HTTP请求数
//
// 参数:
//   method: HTTP方法
//   path: 请求路径
func IncHTTPRequestsInProgress(method, path string) {
	httpRequestsInProgress.WithLabelValues(method, path).Inc()
}

// DecHTTPRequestsInProgress 减少进行中的HTTP请求数
//
// 参数:
//   method: HTTP方法
//   path: 请求路径
func DecHTTPRequestsInProgress(method, path string) {
	httpRequestsInProgress.WithLabelValues(method, path).Dec()
}

// ==================== Redis指标记录函数 ====================

// RecordRedisCommand 记录Redis命令执行
//
// 参数:
//   command: Redis命令名
//   duration: 执行耗时（秒）
//   err: 错误信息
func RecordRedisCommand(command string, duration time.Duration, err error) {
	status := "success"
	if err != nil {
		status = "error"
	}

	// 记录命令总数
	redisCommandsTotal.WithLabelValues(command, status).Inc()

	// 记录命令延迟
	redisCommandDuration.WithLabelValues(command).Observe(duration.Seconds())
}

// SetRedisConnections 设置Redis连接数
//
// 参数:
//   count: 连接数
func SetRedisConnections(count float64) {
	redisConnections.Set(count)
}

// ==================== 业务指标记录函数 ====================

// RecordVersionQuery 记录版本配置查询
//
// 参数:
//   env: 环境（dev/pro）
//   hit: 是否命中
func RecordVersionQuery(env string, hit bool) {
	status := "miss"
	if hit {
		status = "hit"
	}

	versionQueriesTotal.WithLabelValues(env, status).Inc()
}

// RecordGMConfigChange 记录GM配置变更
//
// 参数:
//   changeType: 变更类型
func RecordGMConfigChange(changeType string) {
	gmConfigChangesTotal.WithLabelValues(changeType).Inc()
}

// RecordIPWhitelistCheck 记录IP白名单检查
//
// 参数:
//   allowed: 是否允许
func RecordIPWhitelistCheck(allowed bool) {
	status := "denied"
	if allowed {
		status = "allowed"
	}

	ipWhitelistChecksTotal.WithLabelValues(status).Inc()
}

// IncActiveConnections 增加活跃连接数
func IncActiveConnections() {
	activeConnections.Inc()
}

// DecActiveConnections 减少活跃连接数
func DecActiveConnections() {
	activeConnections.Dec()
}

// SetActiveConnections 设置活跃连接数
//
// 参数:
//   count: 连接数
func SetActiveConnections(count float64) {
	activeConnections.Set(count)
}
