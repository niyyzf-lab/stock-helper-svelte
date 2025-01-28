package data

import (
	"context"
	"time"

	"stock-helper-svelte/backend/api"
)

type Updater struct {
	apiClient *api.Client
	manager   *Manager
	ctx       context.Context
}

func NewUpdater(apiClient *api.Client, ctx context.Context) *Updater {
	manager := NewManager(apiClient)
	manager.SetContext(ctx)
	return &Updater{
		apiClient: apiClient,
		manager:   manager,
		ctx:       ctx,
	}
}

// UpdateData 更新数据
func (u *Updater) UpdateData() error {
	return u.manager.UpdateAllStocks(u.ctx)
}

// GetLastUpdateTime 获取最后更新时间
func (u *Updater) GetLastUpdateTime() (time.Time, error) {
	return u.manager.GetLastUpdateTime()
}

// GetStatus 获取更新状态
func (u *Updater) GetStatus() UpdateStatus {
	return u.manager.GetStatus()
}
