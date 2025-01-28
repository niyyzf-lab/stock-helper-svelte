<script lang="ts">
  import { fade } from 'svelte/transition'
  import { ChevronLeft, Maximize2, Minimize2 } from 'lucide-svelte'
  import type { StockSignal } from '../types/stock'
  import StockChart from './StockChart.svelte'

  export let fileName: string = ''
  export let onClose: () => void = () => {}
  export let isModal: boolean = false

  let loading = false
  let isFullscreen = false
  let record: {
    strategyId: number
    strategyName: string
    executionTime: string
    completionTime: string
    totalStocks: number
    processedStocks: number
    signals: StockSignal[]
  } | null = null

  let selectedSignal: StockSignal | null = null

  async function loadRecordDetail() {
    loading = true
    selectedSignal = null // 重置选中的股票
    try {
      record = await (window as any).go.main.App.GetExecutionRecord(fileName)
      console.log('记录详情:', record)
      if (record?.signals?.length > 0) {
        selectedSignal = record.signals[0]
      }
    } catch (err) {
      console.error('加载记录详情失败:', err)
      record = null
    } finally {
      loading = false
    }
  }

  function handleStockSelect(signal: StockSignal) {
    console.log('选中股票:', signal)
    selectedSignal = signal
  }

  function toggleFullscreen() {
    isFullscreen = !isFullscreen
  }

  // 监听 fileName 变化，重新加载数据
  $: {
    if (fileName) {
      loadRecordDetail()
    }
  }
</script>

<div class="record-detail" class:fullscreen={isFullscreen} in:fade={{duration: 200}}>
  <div class="detail-content" class:fullscreen={isFullscreen}>
    {#if loading}
      <div class="loading-state">
        <div class="spinner"></div>
        <span>加载中...</span>
      </div>
    {:else if record}
      <div class="content-layout" class:fullscreen={isFullscreen}>
        <div class="signals-list">
          {#if record.signals && record.signals.length > 0}
            {#each record.signals as signal}
              <!-- svelte-ignore a11y-click-events-have-key-events -->
              <div 
                class="signal-item" 
                class:active={selectedSignal?.code === signal.code}
                on:click={() => handleStockSelect(signal)}
              >
                <div class="stock-info">
                  <span class="stock-code">{signal.code}</span>
                  <span class="stock-name">{signal.name}</span>
                </div>
                <div class="stock-price">
                  <span class="price">{signal.price.toFixed(2)}</span>
                  <span class="change" class:up={signal.change >= 0} class:down={signal.change < 0}>
                    {signal.change >= 0 ? '+' : ''}{signal.change.toFixed(2)}%
                  </span>
                </div>
              </div>
            {/each}
          {:else}
            <div class="empty-state">
              <span>暂无匹配的股票</span>
            </div>
          {/if}
        </div>

        {#if selectedSignal}
          <div class="chart-panel" class:fullscreen={isFullscreen}>
            <div class="panel-header">
              <div class="stock-info">
                <span class="stock-code">{selectedSignal.code}</span>
                <span class="stock-name">{selectedSignal.name}</span>
                <span class="price">{selectedSignal.price.toFixed(2)}</span>
                <span class="change" class:up={selectedSignal.change >= 0} class:down={selectedSignal.change < 0}>
                  {selectedSignal.change >= 0 ? '+' : ''}{selectedSignal.change.toFixed(2)}%
                </span>
              </div>
              <div class="stock-meta">
                <div class="meta-item">
                  <span class="label">换手率</span>
                  <span class="value">{selectedSignal.turnover.toFixed(2)}%</span>
                </div>
                <div class="meta-item">
                  <span class="label">选股理由</span>
                  <span class="value">{selectedSignal.reason}</span>
                </div>
              </div>
            </div>
            <div class="chart-container" class:fullscreen={isFullscreen}>
              <StockChart code={selectedSignal.code} />
              <button class="fullscreen-btn" on:click={toggleFullscreen}>
                {#if isFullscreen}
                  <Minimize2 size={20} strokeWidth={1.5} />
                {:else}
                  <Maximize2 size={20} strokeWidth={1.5} />
                {/if}
              </button>
            </div>
          </div>
        {/if}
      </div>
    {:else}
      <div class="error-state">
        <span>加载失败</span>
      </div>
    {/if}
  </div>
</div>

<style>
  .record-detail {
    display: flex;
    flex-direction: column;
    height: 100%;
    background: #fff;
    position: relative;
    overflow: hidden;
  }

  .record-detail.fullscreen {
    width: 100vw;
    height: 100vh;
    max-width: none;
    position: fixed;
    top: 0;
    left: 0;
    transform: none;
    border-radius: 0;
    z-index: 1100;
  }

  .detail-content {
    flex: 1;
    overflow: hidden;
    display: flex;
  }

  .content-layout {
    display: flex;
    width: 100%;
    height: 100%;
  }

  .content-layout.fullscreen {
    height: 100vh;
  }

  .signals-list {
    width: 240px;
    min-width: 240px;
    border-right: 1px solid #e5e7eb;
    overflow-y: auto;
    background: #fff;
    height: 100%;
  }

  .signal-item {
    padding: 12px;
    border-bottom: 1px solid #f3f4f6;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .signal-item:hover {
    background: #f9fafb;
  }

  .signal-item.active {
    background: #f3f4f6;
  }

  .signal-item .stock-info {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 8px;
  }

  .signal-item .stock-code {
    font-size: 13px;
    font-weight: 500;
    color: #374151;
  }

  .signal-item .stock-name {
    font-size: 13px;
    color: #6b7280;
  }

  .signal-item .stock-price {
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .signal-item .price {
    font-size: 13px;
    color: #374151;
  }

  .signal-item .change {
    font-size: 12px;
    padding: 2px 6px;
    border-radius: 4px;
  }

  .chart-panel {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    min-width: 0;
  }

  .chart-panel.fullscreen {
    width: 100vw;
  }

  .panel-header {
    padding: 16px 24px;
    background: #fff;
    border-bottom: 1px solid #e5e7eb;
  }

  .panel-header .stock-info {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 12px;
  }

  .panel-header .stock-code {
    font-size: 18px;
    font-weight: 500;
    color: #374151;
  }

  .panel-header .stock-name {
    font-size: 18px;
    color: #6b7280;
  }

  .panel-header .price {
    font-size: 18px;
    font-weight: 500;
    color: #374151;
  }

  .panel-header .change {
    font-size: 16px;
    font-weight: 500;
    padding: 4px 8px;
    border-radius: 4px;
  }

  .stock-meta {
    display: flex;
    gap: 24px;
  }

  .meta-item {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .meta-item .label {
    font-size: 14px;
    color: #6b7280;
  }

  .meta-item .value {
    font-size: 14px;
    color: #374151;
  }

  .change.up {
    color: #ef4444;
    background: rgba(239, 68, 68, 0.1);
  }

  .change.down {
    color: #10b981;
    background: rgba(16, 185, 129, 0.1);
  }

  .chart-container {
    flex: 1;
    min-height: 0;
    background: #fff;
    position: relative;
    overflow: hidden;
  }

  .chart-container.fullscreen {
    height: 100vh;
  }

  .loading-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    gap: 12px;
    color: #6b7280;
  }

  .spinner {
    width: 24px;
    height: 24px;
    border: 2px solid #e5e7eb;
    border-top-color: #2563eb;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  .empty-state, .error-state {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 200px;
    color: #6b7280;
    background: white;
    border-radius: 8px;
    border: 1px solid #e5e7eb;
  }

  .error-state {
    color: #ef4444;
  }

  .fullscreen-btn {
    position: absolute;
    top: 16px;
    right: 16px;
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 6px;
    padding: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    color: #374151;
    z-index: 100;
    transition: all 0.2s;
  }

  .fullscreen-btn:hover {
    background: #f9fafb;
    color: #1f2937;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }
</style>