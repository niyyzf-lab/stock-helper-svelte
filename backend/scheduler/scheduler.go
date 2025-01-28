package scheduler

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"stock-helper-svelte/backend/data"

	"github.com/robfig/cron/v3"
)

// Scheduler 定时任务调度器
type Scheduler struct {
	cron        *cron.Cron
	dataManager *data.Manager
	ctx         context.Context
	mutex       sync.RWMutex
	isRunning   bool
	lastRun     time.Time
	updating    bool // 添加更新状态标志
}

// NewScheduler 创建新的调度器
func NewScheduler(ctx context.Context, dataManager *data.Manager) *Scheduler {
	// 创建一个支持秒级别的定时任务调度器
	c := cron.New(cron.WithSeconds())

	return &Scheduler{
		cron:        c,
		dataManager: dataManager,
		ctx:         ctx,
		isRunning:   false,
		updating:    false,
	}
}

// Start 启动调度器
func (s *Scheduler) Start() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.isRunning {
		return fmt.Errorf("scheduler is already running")
	}

	// 添加每小时更新数据的任务
	// "0 0 * * * *" 表示每小时的开始（分钟和秒都是0）
	_, err := s.cron.AddFunc("0 0 * * * *", func() {
		s.mutex.Lock()
		if s.updating {
			log.Printf("检测到正在进行数据更新，跳过本次自动更新\n")
			s.mutex.Unlock()
			return
		}
		s.updating = true
		s.mutex.Unlock()

		log.Printf("开始执行定时数据更新任务: %s\n", time.Now().Format("2006-01-02 15:04:05"))

		if err := s.dataManager.UpdateAllStocks(s.ctx); err != nil {
			log.Printf("数据更新失败: %v\n", err)
		} else {
			s.lastRun = time.Now()
			log.Printf("数据更新完成: %s\n", s.lastRun.Format("2006-01-02 15:04:05"))
		}

		s.mutex.Lock()
		s.updating = false
		s.mutex.Unlock()
	})

	if err != nil {
		return fmt.Errorf("failed to add cron job: %v", err)
	}

	// 启动定时任务
	s.cron.Start()
	s.isRunning = true
	log.Println("定时任务调度器已启动")

	return nil
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.isRunning {
		s.cron.Stop()
		s.isRunning = false
		log.Println("定时任务调度器已停止")
	}
}

// GetStatus 获取调度器状态
func (s *Scheduler) GetStatus() map[string]interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return map[string]interface{}{
		"isRunning": s.isRunning,
		"lastRun":   s.lastRun.Format("2006-01-02 15:04:05"),
		"updating":  s.updating, // 添加更新状态到返回值
	}
}

// RunNow 立即执行一次数据更新
func (s *Scheduler) RunNow() error {
	s.mutex.Lock()
	if s.updating {
		s.mutex.Unlock()
		return fmt.Errorf("数据更新正在进行中，请稍后再试")
	}
	s.updating = true
	s.mutex.Unlock()

	defer func() {
		s.mutex.Lock()
		s.updating = false
		s.mutex.Unlock()
	}()

	log.Printf("手动触发数据更新任务: %s\n", time.Now().Format("2006-01-02 15:04:05"))

	if err := s.dataManager.UpdateAllStocks(s.ctx); err != nil {
		return fmt.Errorf("数据更新失败: %v", err)
	}

	s.lastRun = time.Now()
	log.Printf("手动数据更新完成: %s\n", s.lastRun.Format("2006-01-02 15:04:05"))
	return nil
}
