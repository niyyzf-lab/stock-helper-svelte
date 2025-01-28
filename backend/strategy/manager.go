package strategy

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"stock-helper-svelte/backend/api"
	"stock-helper-svelte/backend/engine"
	"stock-helper-svelte/backend/types"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ExecutionUpdater 定义了执行状态更新的接口
type ExecutionUpdater interface {
	UpdateProgress(processedStocks int, currentStock string)
	AddSignal(types.StockSignal)
}

// Manager 策略管理器
type Manager struct {
	basePath  string              // 策略文件基础路径
	apiClient *api.Client         // API客户端
	ctx       context.Context     // 上下文
	engine    *engine.Engine      // 执行引擎
	mutex     sync.RWMutex        // 读写锁
	signals   []types.StockSignal // 当前执行的信号
}

// statusUpdater 实现 engine.StatusUpdater 接口
type statusUpdater struct {
	ctx     context.Context
	manager *Manager
}

func (s *statusUpdater) UpdateStatus(status engine.ExecutionStatus) {
	runtime.EventsEmit(s.ctx, "engine:status", status)
}

func (s *statusUpdater) UpdateProgress(processedStocks int, currentStock string) {
	// 这个方法可能不需要实现,因为 UpdateStatus 已经包含了进度信息
}

func (s *statusUpdater) AddSignal(signal types.StockSignal) {
	s.manager.mutex.Lock()
	s.manager.signals = append(s.manager.signals, signal)
	s.manager.mutex.Unlock()
	runtime.EventsEmit(s.ctx, "engine:signal", signal)
}

// NewManager 创建新的策略管理器
func NewManager(basePath string, apiClient *api.Client, ctx context.Context) *Manager {
	manager := &Manager{
		basePath:  basePath,
		apiClient: apiClient,
		ctx:       ctx,
		signals:   make([]types.StockSignal, 0),
	}

	// 创建状态更新器
	updater := &statusUpdater{
		ctx:     ctx,
		manager: manager,
	}

	// 创建执行引擎配置
	config := engine.ExecutionConfig{
		WorkerPoolSize:   16,
		BatchSize:        100,
		RetryAttempts:    3,
		RetryDelay:       time.Second,
		ExecutionTimeout: 24 * time.Hour,
		APIClient:        apiClient,
		Context:          ctx,
	}

	// 创建执行引擎
	eng, err := engine.NewEngine(config, updater)
	if err != nil {
		fmt.Printf("Warning: Failed to create engine: %v\n", err)
	}
	manager.engine = eng

	return manager
}

// GetCurrentStatus 获取当前执行状态
func (m *Manager) GetCurrentStatus() engine.ExecutionStatus {
	if m.engine == nil {
		return engine.ExecutionStatus{
			Status: engine.StatusIdle,
		}
	}
	return m.engine.GetStatus()
}

// ExecuteStrategy 执行策略
func (m *Manager) ExecuteStrategy(strategy *engine.Strategy) error {
	if m.engine == nil {
		return fmt.Errorf("engine not initialized")
	}

	// 重置信号列表
	m.mutex.Lock()
	m.signals = make([]types.StockSignal, 0)
	m.mutex.Unlock()

	// 执行策略
	err := m.engine.Execute(strategy)

	// 获取当前状态
	status := m.engine.GetStatus()
	fmt.Printf("策略执行完成，状态: %s\n", status.Status)

	// 如果执行完成或被中止，保存结果
	if status.Status == engine.StatusCompleted || status.Status == engine.StatusStopped {
		fmt.Printf("准备保存执行结果...\n")

		// 获取当前信号列表的副本
		m.mutex.RLock()
		signals := make([]types.StockSignal, len(m.signals))
		copy(signals, m.signals)
		m.mutex.RUnlock()

		// 创建执行结果
		result := &engine.ExecutionResult{
			StrategyID:      strategy.ID,
			StrategyName:    strategy.Name,
			ExecutionTime:   status.StartTime,
			CompletionTime:  time.Now(),
			TotalStocks:     status.TotalStocks,
			ProcessedStocks: status.ProcessedCount,
			Signals:         signals,
		}

		// 生成安全的文件名
		safeName := regexp.MustCompile(`[^a-zA-Z0-9_-]`).ReplaceAllString(strategy.Name, "_")
		fileName := fmt.Sprintf("strategy_%s_%s.json",
			safeName,
			status.StartTime.Format("20060102_150405"))

		// 保存结果
		recordDir, err := m.getRecordDir()
		if err != nil {
			fmt.Printf("获取记录目录失败: %v\n", err)
			return fmt.Errorf("保存结果失败: %v", err)
		}
		fmt.Printf("记录目录: %s\n", recordDir)

		filePath := filepath.Join(recordDir, fileName)
		fmt.Printf("准备写入文件: %s\n", filePath)

		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			fmt.Printf("序列化结果失败: %v\n", err)
			return fmt.Errorf("序列化结果失败: %v", err)
		}

		if err := os.WriteFile(filePath, data, 0644); err != nil {
			fmt.Printf("写入文件失败: %v\n", err)
			return fmt.Errorf("写入结果文件失败: %v", err)
		}
		fmt.Printf("执行结果已保存到: %s\n", filePath)
	}

	return err
}

// Pause 暂停执行
func (m *Manager) Pause() {
	if m.engine != nil {
		m.engine.Pause()
	}
}

// Resume 恢复执行
func (m *Manager) Resume() {
	if m.engine != nil {
		m.engine.Resume()
	}
}

// Stop 停止执行
func (m *Manager) Stop() {
	if m.engine != nil {
		m.engine.Stop()

		// 等待状态变为已停止
		for {
			status := m.engine.GetStatus()
			if status.Status == engine.StatusStopped {
				// 获取当前信号列表的副本
				m.mutex.RLock()
				signals := make([]types.StockSignal, len(m.signals))
				copy(signals, m.signals)
				m.mutex.RUnlock()

				// 创建执行结果
				result := &engine.ExecutionResult{
					StrategyID:      status.StrategyId,
					StrategyName:    m.getStrategyName(status.StrategyId),
					ExecutionTime:   status.StartTime,
					CompletionTime:  time.Now(),
					TotalStocks:     status.TotalStocks,
					ProcessedStocks: status.ProcessedCount,
					Signals:         signals,
				}

				// 生成安全的文件名
				safeName := regexp.MustCompile(`[^a-zA-Z0-9_-]`).ReplaceAllString(result.StrategyName, "_")
				fileName := fmt.Sprintf("strategy_%s_%s.json",
					safeName,
					status.StartTime.Format("20060102_150405"))

				// 保存结果
				recordDir, err := m.getRecordDir()
				if err != nil {
					fmt.Printf("获取记录目录失败: %v\n", err)
					return
				}

				filePath := filepath.Join(recordDir, fileName)
				data, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					fmt.Printf("序列化结果失败: %v\n", err)
					return
				}

				if err := os.WriteFile(filePath, data, 0644); err != nil {
					fmt.Printf("写入结果文件失败: %v\n", err)
					return
				}

				fmt.Printf("执行结果已保存到: %s\n", filePath)
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// getStrategyName 根据策略ID获取策略名称
func (m *Manager) getStrategyName(id int) string {
	strategy, err := m.GetStrategyByID(id)
	if err != nil {
		return fmt.Sprintf("Strategy_%d", id)
	}
	return strategy.Name
}

// GetCurrentSignals 获取当前执行收集到的信号
func (m *Manager) GetCurrentSignals() []types.StockSignal {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.signals
}

// Close 关闭管理器
func (m *Manager) Close() error {
	if m.engine != nil {
		return m.engine.Close()
	}
	return nil
}

// 以下是策略文件管理相关的方法

// parseStrategyMeta 从Lua文件中解析元数据
func (m *Manager) parseStrategyMeta(filePath string) (*engine.StrategyMeta, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var meta engine.StrategyMeta
	scanner := bufio.NewScanner(file)

	// 用于匹配元数据的正则表达式
	idRegex := regexp.MustCompile(`--\s*@id:\s*(\d+)`)
	nameRegex := regexp.MustCompile(`--\s*@name:\s*(.+)`)
	descRegex := regexp.MustCompile(`--\s*@description:\s*(.+)`)

	for scanner.Scan() {
		line := scanner.Text()

		// 跳过非注释行
		if !strings.HasPrefix(strings.TrimSpace(line), "--") {
			break
		}

		// 解析ID
		if matches := idRegex.FindStringSubmatch(line); len(matches) > 1 {
			id, err := strconv.Atoi(matches[1])
			if err != nil {
				return nil, fmt.Errorf("invalid strategy ID: %s", matches[1])
			}
			meta.ID = id
			continue
		}

		// 解析名称
		if matches := nameRegex.FindStringSubmatch(line); len(matches) > 1 {
			meta.Name = strings.TrimSpace(matches[1])
			continue
		}

		// 解析描述
		if matches := descRegex.FindStringSubmatch(line); len(matches) > 1 {
			meta.Description = strings.TrimSpace(matches[1])
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// 验证必要字段
	if meta.ID == 0 || meta.Name == "" {
		return nil, fmt.Errorf("missing required metadata in strategy file: %s", filePath)
	}

	return &meta, nil
}

// GetStrategies 获取所有策略
func (m *Manager) GetStrategies() []engine.Strategy {
	pattern := filepath.Join(m.basePath, "*.lua")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil
	}

	var strategies []engine.Strategy
	for _, file := range files {
		meta, err := m.parseStrategyMeta(file)
		if err != nil {
			continue
		}

		strategies = append(strategies, engine.Strategy{
			ID:          meta.ID,
			Name:        meta.Name,
			Description: meta.Description,
			FilePath:    file,
		})
	}

	return strategies
}

// GetStrategyByID 根据ID获取策略
func (m *Manager) GetStrategyByID(id int) (*engine.Strategy, error) {
	strategies := m.GetStrategies()
	for _, s := range strategies {
		if s.ID == id {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("strategy not found: %d", id)
}

// 执行记录相关方法

// getRecordDir 获取记录存储目录
func (m *Manager) getRecordDir() (string, error) {
	appDataDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("无法获取用户配置目录: %v", err)
	}

	recordDir := filepath.Join(appDataDir, "stock-helper-svelte.exe", "records")
	if err := os.MkdirAll(recordDir, 0755); err != nil {
		return "", fmt.Errorf("无法创建记录目录: %v", err)
	}

	return recordDir, nil
}

// GetExecutionRecords 获取执行记录列表
func (m *Manager) GetExecutionRecords() ([]engine.ExecutionRecord, error) {
	recordDir, err := m.getRecordDir()
	if err != nil {
		return nil, err
	}

	files, err := filepath.Glob(filepath.Join(recordDir, "strategy_*.json"))
	if err != nil {
		return nil, fmt.Errorf("无法读取记录文件: %v", err)
	}

	var records []engine.ExecutionRecord
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Warning: 无法读取文件 %s: %v\n", file, err)
			continue
		}

		var result engine.ExecutionResult
		if err := json.Unmarshal(data, &result); err != nil {
			fmt.Printf("Warning: 无法解析文件 %s: %v\n", file, err)
			continue
		}

		record := engine.ExecutionRecord{
			FileName:       filepath.Base(file),
			StrategyID:     result.StrategyID,
			StrategyName:   result.StrategyName,
			ExecutionTime:  result.ExecutionTime,
			SignalCount:    len(result.Signals),
			ProcessedCount: result.ProcessedStocks,
			TotalStocks:    result.TotalStocks,
		}
		records = append(records, record)
	}

	// 按执行时间降序排序
	sort.Slice(records, func(i, j int) bool {
		return records[i].ExecutionTime.After(records[j].ExecutionTime)
	})

	return records, nil
}

// GetExecutionRecord 获取具体执行记录内容
func (m *Manager) GetExecutionRecord(fileName string) (*engine.ExecutionResult, error) {
	recordDir, err := m.getRecordDir()
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(fileName, "strategy_") || !strings.HasSuffix(fileName, ".json") {
		return nil, fmt.Errorf("无效的文件名格式")
	}

	filePath := filepath.Join(recordDir, fileName)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法读取记录文件: %v", err)
	}

	var result engine.ExecutionResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("无法解析记录文件: %v", err)
	}

	return &result, nil
}

// DeleteExecutionRecord 删除执行记录
func (m *Manager) DeleteExecutionRecord(fileName string) error {
	recordDir, err := m.getRecordDir()
	if err != nil {
		return err
	}

	if !strings.HasPrefix(fileName, "strategy_") || !strings.HasSuffix(fileName, ".json") {
		return fmt.Errorf("无效的文件名格式")
	}

	filePath := filepath.Join(recordDir, fileName)
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("无法删除记录文件: %v", err)
	}

	return nil
}
