package engine

import (
	"runtime"
	"sync/atomic"
	"time"
)

// ExecutionMetrics 执行指标收集器
type ExecutionMetrics struct {
	startTime      time.Time
	processedCount atomic.Int32
	errorCount     atomic.Int32
	totalStocks    int32
	lastUpdateTime time.Time
	lastProcessed  int32
}

// NewExecutionMetrics 创建新的指标收集器
func NewExecutionMetrics(totalStocks int32) *ExecutionMetrics {
	return &ExecutionMetrics{
		startTime:      time.Now(),
		totalStocks:    totalStocks,
		lastUpdateTime: time.Now(),
	}
}

// IncrementProcessed 增加处理计数
func (m *ExecutionMetrics) IncrementProcessed() {
	m.processedCount.Add(1)
}

// IncrementErrors 增加错误计数
func (m *ExecutionMetrics) IncrementErrors() {
	m.errorCount.Add(1)
}

// GetStats 获取当前统计信息
func (m *ExecutionMetrics) GetStats() ExecutionStats {
	currentTime := time.Now()
	processed := m.processedCount.Load()

	// 计算处理速度
	timeDiff := currentTime.Sub(m.lastUpdateTime).Seconds()
	processedDiff := processed - m.lastProcessed
	speed := float64(processedDiff) / timeDiff

	// 更新最后的统计时间和处理数量
	m.lastUpdateTime = currentTime
	m.lastProcessed = processed

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return ExecutionStats{
		StartTime:       m.startTime,
		ProcessedStocks: processed,
		TotalStocks:     m.totalStocks,
		ErrorCount:      m.errorCount.Load(),
		CurrentSpeed:    speed,
		MemoryUsage:     memStats.Alloc,
		GoroutineCount:  runtime.NumGoroutine(),
	}
}

// Reset 重置指标
func (m *ExecutionMetrics) Reset() {
	m.startTime = time.Now()
	m.processedCount.Store(0)
	m.errorCount.Store(0)
	m.lastUpdateTime = time.Now()
	m.lastProcessed = 0
}
