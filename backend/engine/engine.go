package engine

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"stock-helper-svelte/backend/api"
)

// Engine 执行引擎实现
type Engine struct {
	config     ExecutionConfig
	apiClient  *api.Client
	metrics    *ExecutionMetrics
	workerPool *WorkerPool

	// 状态控制
	state     *engineState
	stateLock sync.RWMutex

	// 上下文控制
	ctx    context.Context
	cancel context.CancelFunc

	// 状态更新器
	statusUpdater StatusUpdater
}

// NewEngine 创建新的执行引擎
func NewEngine(config ExecutionConfig, statusUpdater StatusUpdater) (*Engine, error) {
	// 验证配置
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(config.Context, config.ExecutionTimeout)

	return &Engine{
		config:        config,
		apiClient:     config.APIClient,
		ctx:           ctx,
		cancel:        cancel,
		statusUpdater: statusUpdater,
		state:         newEngineState(),
	}, nil
}

// validateConfig 验证配置
func validateConfig(config *ExecutionConfig) error {
	if config.APIClient == nil {
		return NewInvalidConfigError("APIClient", fmt.Errorf("API client is required"))
	}

	if config.Context == nil {
		return NewInvalidConfigError("Context", fmt.Errorf("context is required"))
	}

	if config.WorkerPoolSize <= 0 {
		config.WorkerPoolSize = 16
	}
	if config.BatchSize <= 0 {
		config.BatchSize = 100
	}
	if config.RetryAttempts <= 0 {
		config.RetryAttempts = 3
	}
	if config.RetryDelay <= 0 {
		config.RetryDelay = time.Second
	}
	if config.ExecutionTimeout <= 0 {
		config.ExecutionTimeout = 24 * time.Hour
	}

	return nil
}

// GetStatus 获取当前执行状态
func (e *Engine) GetStatus() ExecutionStatus {
	e.stateLock.RLock()
	defer e.stateLock.RUnlock()
	return e.state.toExecutionStatus()
}

// updateState 更新内部状态
func (e *Engine) updateState(update func(*engineState)) {
	e.stateLock.Lock()
	update(e.state)
	status := e.state.toExecutionStatus()
	e.stateLock.Unlock()

	// 通知状态更新
	if e.statusUpdater != nil {
		e.statusUpdater.UpdateStatus(status)
	}
}

// Execute 执行策略
func (e *Engine) Execute(strategy *Strategy) error {
	// 检查是否已经在运行
	if e.IsRunning() {
		return NewEngineError(ErrEngineAlreadyRunning, "engine is already running", fmt.Errorf("engine is already running"))
	}

	// 获取股票列表
	stocks, err := e.apiClient.GetIndexList()
	if err != nil {
		return NewAPIRequestError("failed to get stock list", err)
	}

	// 过滤股票列表
	stocks = e.filterStocks(stocks)

	// 初始化 metrics
	e.metrics = NewExecutionMetrics(int32(len(stocks)))

	// 初始化状态
	e.updateState(func(s *engineState) {
		s.status = StatusRunning
		s.startTime = time.Now()
		s.totalStocks = int32(len(stocks))
		s.processedCount = 0
		s.currentStock = ""
		s.speed = 0
		s.error = ""
		s.strategyId = strategy.ID
		s.paused = false
		s.shouldStop = false
	})

	// 创建工作池
	pool, err := NewWorkerPool(e.config.WorkerPoolSize, strategy, e.metrics, e.ctx, e.apiClient, e.statusUpdater)
	if err != nil {
		e.updateState(func(s *engineState) {
			s.status = StatusError
			s.error = err.Error()
		})
		return NewEngineError(ErrWorkerPoolCreation, "failed to create worker pool", err)
	}

	// 保存工作池引用
	e.workerPool = pool

	// 确保在函数返回时清理工作池
	defer func() {
		if e.workerPool != nil {
			e.workerPool.Close()
			e.workerPool = nil
		}
	}()

	// 启动工作池
	pool.Start(e.ctx)

	// 启动状态更新协程
	go e.updateStatus()

	// 分批处理股票
	for i := 0; i < len(stocks); i += e.config.BatchSize {
		// 检查是否应该停止
		if e.state.shouldStop {
			e.updateState(func(s *engineState) {
				if s.status != StatusCompleted { // 如果不是正常完成,设置为停止状态
					s.status = StatusStopped
				}
			})
			return nil
		}

		// 检查是否暂停
		for e.IsPaused() && !e.state.shouldStop {
			time.Sleep(100 * time.Millisecond)
		}

		end := i + e.config.BatchSize
		if end > len(stocks) {
			end = len(stocks)
		}

		// 提交一批股票到工作池
		for _, stock := range stocks[i:end] {
			e.updateState(func(s *engineState) {
				s.currentStock = fmt.Sprintf("%s(%s)", stock.Name, stock.Code)
			})
			pool.Submit(stock)
		}
	}

	// 等待所有任务完成
	pool.Wait()

	// 确保最终状态正确
	e.updateState(func(s *engineState) {
		if s.status != StatusError && s.status != StatusStopped {
			s.status = StatusCompleted
			s.currentStock = ""
			s.shouldStop = true
		}
	})

	return nil
}

// updateStatus 更新执行状态
func (e *Engine) updateStatus() {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			stats := e.metrics.GetStats()
			e.updateState(func(s *engineState) {
				s.processedCount = int32(stats.ProcessedStocks)
				s.speed = stats.CurrentSpeed
			})
		case <-e.ctx.Done():
			return
		}
	}
}

// getStatusString 获取当前状态字符串
func (e *Engine) getStatusString() string {
	if e.state.shouldStop {
		return "stopped"
	}
	if e.state.paused {
		return "paused"
	}
	if e.state.status == StatusRunning {
		return "running"
	}
	return "idle"
}

// filterStocks 过滤股票列表
func (e *Engine) filterStocks(stocks []api.Index) []api.Index {
	filtered := make([]api.Index, 0, len(stocks))
	for _, stock := range stocks {
		// 过滤ST股票
		if strings.Contains(strings.ToUpper(stock.Name), "ST") {
			continue
		}

		// 过滤退市股票
		if strings.Contains(stock.Name, "退") || strings.Contains(stock.Name, "退市") {
			continue
		}

		// 只保留沪深主板、创业板、科创板的股票
		if strings.HasPrefix(stock.Code, "00") || // 深证主板
			strings.HasPrefix(stock.Code, "30") || // 创业板
			strings.HasPrefix(stock.Code, "60") || // 上证主板
			strings.HasPrefix(stock.Code, "68") { // 科创板
			filtered = append(filtered, stock)
		}
	}
	return filtered
}

// Pause 暂停执行
func (e *Engine) Pause() {
	e.updateState(func(s *engineState) {
		s.paused = true
		s.status = StatusPaused
	})
}

// Resume 恢复执行
func (e *Engine) Resume() {
	e.updateState(func(s *engineState) {
		s.paused = false
		s.status = StatusRunning
	})
}

// Stop 停止执行
func (e *Engine) Stop() {
	e.updateState(func(s *engineState) {
		s.shouldStop = true
		s.status = StatusStopped
	})
}

// IsPaused 是否暂停
func (e *Engine) IsPaused() bool {
	e.stateLock.RLock()
	defer e.stateLock.RUnlock()
	return e.state.paused
}

// IsRunning 是否正在运行
func (e *Engine) IsRunning() bool {
	e.stateLock.RLock()
	defer e.stateLock.RUnlock()
	return e.state.status == StatusRunning
}

// GetExecutionStats 获取执行统计信息
func (e *Engine) GetExecutionStats() ExecutionStats {
	if e.metrics == nil {
		return ExecutionStats{}
	}
	return e.metrics.GetStats()
}

// Close 关闭执行引擎
func (e *Engine) Close() error {
	e.Stop()
	if e.cancel != nil {
		e.cancel()
	}
	if e.workerPool != nil {
		e.workerPool.Close()
	}
	return nil
}
