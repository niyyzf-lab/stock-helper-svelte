<script lang="ts">
  import { fade, fly } from 'svelte/transition'
  import { quintOut } from 'svelte/easing'
  import { Target } from 'lucide-svelte'
  import StockChart from '../StockChart.svelte'
  import type { TestResult } from 'src/types/stock'

  export let currentStock: any
  export let currentStockIndex: number
  export let stocks: any[]
  export let showActions: boolean
  export let onAnswer: (direction: 'up' | 'down' | 'shock') => void

  // 格式化日期
  function formatDate(dateStr: string, timeStr?: string) {
    if (!dateStr) return '未知时间'
    try {
      const date = new Date(dateStr)
      let result = date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit'
      })
      if (timeStr) {
        result += ' ' + timeStr
      }
      return result
    } catch {
      return '无效日期'
    }
  }
</script>

<div class="test-container">
  <div class="chart-wrapper">
    <StockChart 
      code={currentStock.dm}
      endTime={currentStock.testDate}
      freq="dh"
    />

    <!-- 添加悬浮操作按钮 -->
    {#if showActions}
      <div class="floating-actions" in:fly={{y: -20, duration: 400, easing: quintOut}}>
        <div class="action-buttons">
          <button class="action-btn up" on:click={() => onAnswer('up')}>
            <span class="btn-icon">↑</span>
            <span class="btn-label">看涨</span>
          </button>
          <button class="action-btn shock" on:click={() => onAnswer('shock')}>
            <span class="btn-icon">↔</span>
            <span class="btn-label">震荡</span>
          </button>
          <button class="action-btn down" on:click={() => onAnswer('down')}>
            <span class="btn-icon">↓</span>
            <span class="btn-label">看跌</span>
          </button>
        </div>
      </div>
    {/if}
  </div>
  
  <div class="test-footer">
    <div class="footer-content">
      <!-- 左侧股票信息 -->
      <div class="stock-info">
        <div class="stock-basic">
          <span class="stock-code">{currentStock.dm}</span>
          <span class="divider">·</span>
          <span class="stock-name">{currentStock.mc}</span>
        </div>
        <span class="test-time">
          {formatDate(currentStock.testDate)}
        </span>
      </div>

      <!-- 右侧进度信息 -->
      <div class="progress-info">
        <Target size={14} strokeWidth={1.5} />
        <span class="progress-text">
          <span class="current">{currentStockIndex + 1}</span>
          <span class="divider">/</span>
          <span class="total">{stocks.length}</span>
        </span>
      </div>
    </div>
  </div>
</div>

<style>
  .test-container {
    display: flex;
    flex-direction: column;
    height: 100%;
    gap: 1px;
    background: var(--neutral-100);
    border-radius: var(--radius-lg);
    overflow: hidden;
  }

  .chart-wrapper {
    position: relative;
    flex: 1;
    min-height: 0;
    background: white;
    overflow: hidden;
  }

  .test-footer {
    background: white;
    padding: 12px 16px;
    border-top: 1px solid var(--border-color);
    border-bottom-left-radius: var(--radius-lg);
    border-bottom-right-radius: var(--radius-lg);
  }

  .footer-content {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .stock-info {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .stock-basic {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .stock-code {
    font-family: var(--font-mono);
    font-size: 15px;
    font-weight: 600;
    color: var(--primary-600);
  }

  .stock-name {
    font-size: 15px;
    font-weight: 500;
    color: var(--text-primary);
  }

  .test-time {
    font-size: 13px;
    color: var(--text-secondary);
  }

  .divider {
    color: var(--neutral-300);
  }

  .progress-info {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 6px 12px;
    background: var(--neutral-50);
    border: 1px solid var(--border-color);
    border-radius: 24px;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  }

  .progress-text {
    font-size: 14px;
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .progress-info .current {
    color: var(--primary-600);
    font-weight: 600;
  }

  .progress-info .total {
    color: var(--text-secondary);
  }

  .floating-actions {
    position: absolute;
    top: 24px;
    right: 24px;
    z-index: 10;
  }

  .action-buttons {
    display: flex;
    gap: 8px;
    padding: 8px;
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(12px);
    -webkit-backdrop-filter: blur(12px);
    border: 1px solid var(--border-color);
    border-radius: 16px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  }

  .action-btn {
    min-width: 80px;
    height: 36px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 6px;
    border: none;
    border-radius: 10px;
    font-size: 14px;
    font-weight: 500;
    color: white;
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  }

  .action-btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.15);
  }

  .action-btn:active {
    transform: translateY(0) scale(0.98);
  }

  .action-btn.up {
    background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
  }

  .action-btn.down {
    background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  }

  .action-btn.shock {
    background: linear-gradient(135deg, #6366f1 0%, #4f46e5 100%);
  }

  .btn-icon {
    font-size: 16px;
    font-weight: 600;
  }
</style> 