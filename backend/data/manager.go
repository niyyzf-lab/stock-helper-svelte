package data

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"stock-helper-svelte/backend/api"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// UpdateStatus 更新状态
type UpdateStatus struct {
	IsUpdating    bool      `json:"isUpdating"`    // 是否正在更新
	StartTime     time.Time `json:"startTime"`     // 开始时间
	Total         int       `json:"total"`         // 总数
	Completed     int       `json:"completed"`     // 已完成数
	Current       string    `json:"current"`       // 当前处理的股票
	Progress      float64   `json:"progress"`      // 进度百分比
	Speed         float64   `json:"speed"`         // 速度(个/秒)
	EstimateTime  int       `json:"estimateTime"`  // 预计剩余时间(秒)
	ErrorCount    int       `json:"errorCount"`    // 错误数量
	LastError     string    `json:"lastError"`     // 最后一个错误
	LastUpdateStr string    `json:"lastUpdateStr"` // 最后更新时间
}

const (
	MaxConcurrent = 64 // 最大并发数
)

// Manager 数据管理器
type Manager struct {
	apiClient *api.Client
	ctx       context.Context
	mutex     sync.RWMutex
	status    UpdateStatus
}

// NewManager 创建新的数据管理器
func NewManager(apiClient *api.Client) *Manager {
	return &Manager{
		apiClient: apiClient,
	}
}

// SetContext 设置上下文
func (m *Manager) SetContext(ctx context.Context) {
	m.ctx = ctx
}

// GetStatus 获取更新状态
func (m *Manager) GetStatus() UpdateStatus {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.status
}

// emitStatus 发送状态更新事件
func (m *Manager) emitStatus() {
	if m.ctx != nil {
		m.mutex.RLock()
		status := m.status
		m.mutex.RUnlock()
		runtime.EventsEmit(m.ctx, "update:status", status)
	}
}

// GetStockData 获取股票数据
func (m *Manager) GetStockData(code string) ([]api.KLineData, []api.HistoricalTransaction, error) {
	// 并发请求K线数据和历史成交数据
	var wg sync.WaitGroup
	errChan := make(chan error, 2)
	var klineData []api.KLineData
	var transData []api.HistoricalTransaction
	var klineErr, transErr error

	wg.Add(2)

	// 获取K线数据
	go func() {
		defer wg.Done()
		var err error
		klineData, err = m.apiClient.GetKLineData(code, "dh")
		if err != nil {
			errChan <- fmt.Errorf("获取K线数据失败: %v", err)
			klineErr = err
		}
	}()

	// 获取历史成交数据
	go func() {
		defer wg.Done()
		var err error
		transData, err = m.apiClient.GetHistoricalTransactions(code)
		if err != nil {
			errChan <- fmt.Errorf("获取历史成交数据失败: %v", err)
			transErr = err
		}
	}()

	// 等待所有请求完成
	wg.Wait()

	// 检查错误
	if klineErr != nil {
		return nil, nil, klineErr
	}
	if transErr != nil {
		return nil, nil, transErr
	}

	// 确保两个数据都获取成功
	if len(klineData) == 0 {
		return nil, nil, fmt.Errorf("获取K线数据为空")
	}

	return klineData, transData, nil
}

// filterStocks 过滤掉ST股票和退市股票
func (m *Manager) filterStocks(stocks []api.Index) []api.Index {
	filtered := make([]api.Index, 0, len(stocks))
	for _, stock := range stocks {
		// 过滤掉ST股票
		if strings.Contains(strings.ToUpper(stock.Name), "ST") {
			continue
		}

		// 过滤掉退市股票
		if strings.Contains(stock.Name, "退") {
			continue
		}

		// 过滤掉含有 "退市" 的股票
		if strings.Contains(stock.Name, "退市") {
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

// UpdateAllStocks 更新所有股票数据
func (m *Manager) UpdateAllStocks(ctx context.Context) error {
	m.ctx = ctx
	m.mutex.Lock()
	if m.status.IsUpdating {
		m.mutex.Unlock()
		return fmt.Errorf("数据更新正在进行中")
	}

	// 初始化状态
	m.status = UpdateStatus{
		IsUpdating: true,
		StartTime:  time.Now(),
		Completed:  0,
		Total:      0,
	}
	m.mutex.Unlock()
	m.emitStatus() // 立即发送初始状态

	// 确保结束时更新状态
	defer func() {
		m.mutex.Lock()
		m.status.IsUpdating = false
		m.status.LastUpdateStr = time.Now().Format("2006-01-02 15:04:05")
		m.mutex.Unlock()
		m.emitStatus()
	}()

	// 获取股票列表
	indices, err := m.apiClient.GetIndexList()
	if err != nil {
		return fmt.Errorf("failed to get index list: %v", err)
	}

	// 过滤股票列表
	indices = m.filterStocks(indices)
	log.Printf("过滤后的股票数量: %d\n", len(indices))

	m.mutex.Lock()
	m.status.Total = len(indices)
	m.mutex.Unlock()
	m.emitStatus() // 发送总数更新

	// 创建工作池
	sem := make(chan struct{}, MaxConcurrent)
	var wg sync.WaitGroup
	completed := 0

	// 遍历所有股票
	for _, index := range indices {
		wg.Add(1)
		sem <- struct{}{} // 获取信号量

		go func(idx api.Index) {
			defer func() {
				<-sem // 释放信号量
				wg.Done()
			}()

			// 获取数据，缓存由API客户端管理
			_, _, err := m.GetStockData(idx.Code)

			m.mutex.Lock()
			completed++
			m.status.Completed = completed
			m.status.Current = fmt.Sprintf("%s(%s) (%d/%d)", idx.Name, idx.Code, completed, m.status.Total)

			// 更新进度
			elapsed := time.Since(m.status.StartTime).Seconds()
			m.status.Speed = float64(completed) / elapsed
			m.status.Progress = float64(completed) / float64(m.status.Total) * 100
			if m.status.Speed > 0 {
				remaining := int((float64(m.status.Total-completed) / m.status.Speed))
				m.status.EstimateTime = remaining
			}

			// 只有在所有接口都失败时才计入错误
			if err != nil && strings.Contains(err.Error(), "所有接口请求失败") {
				m.status.ErrorCount++
				m.status.LastError = fmt.Sprintf("%s: %v", idx.Code, err)
			}
			m.mutex.Unlock()
			m.emitStatus()
		}(index)
	}

	// 等待所有更新完成
	wg.Wait()
	return nil
}

// GetLastUpdateTime 获取最后更新时间
func (m *Manager) GetLastUpdateTime() (time.Time, error) {
	// 这里可以直接返回当前时间，因为缓存由API客户端管理
	return time.Now(), nil
}
