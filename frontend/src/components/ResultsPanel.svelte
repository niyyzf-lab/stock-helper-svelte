<script lang="ts">
  import { fade, fly } from 'svelte/transition'
  import { onDestroy } from 'svelte'
  import StockChart from './StockChart.svelte'
  import { ChevronLeft, ChevronRight } from 'lucide-svelte'

  interface StockSignal {
    code: string;
    name: string;
    price: number;
    change: number;
    turnover: number;
    reason: string;
  }

  export let stockSignals: StockSignal[] = []
  export let totalStocks: number = 0
  export let matchedStocks: number = 0

  // 添加侧边栏状态
  let isCollapsed = true;
  let isHovered = false;
  let listWidth = isCollapsed ? 120 : 320;

  // 格式化数字
  function formatNumber(num: number, decimals: number = 2): string {
    if (num === undefined || num === null) return '0.00'
    return num.toFixed(decimals)
  }

  // 格式化百分比
  function formatPercent(num: number): string {
    if (num === undefined || num === null) return '0.00%'
    return (num).toFixed(2) + '%'
  }

  // 计算匹配率
  $: matchRate = totalStocks > 0 ? ((matchedStocks / totalStocks) * 100).toFixed(1) + '%' : '0.0%'

  // 安全获取数据
  $: safeSignals = stockSignals || []
  $: safeTotalStocks = typeof totalStocks === 'number' ? totalStocks : 0
  $: safeMatchedStocks = typeof matchedStocks === 'number' ? matchedStocks : 0

  // 选中的股票
  let selectedStock: StockSignal | null = null

  // 处理股票选择
  function handleStockSelect(stock: StockSignal) {
    selectedStock = stock
  }

  // 处理鼠标进入
  function handleMouseEnter() {
    if (isCollapsed) {
      isHovered = true;
    }
  }

  // 处理鼠标离开
  function handleMouseLeave() {
    isHovered = false;
  }

  // 切换侧边栏状态
  function toggleSidebar() {
    isCollapsed = !isCollapsed;
    listWidth = isCollapsed ? 120 : 320;
  }

  // 调试输出
  $: console.log('ResultsPanel state:', {
    totalStocks: safeTotalStocks,
    matchedStocks: safeMatchedStocks,
    matchRate,
    signalsCount: safeSignals.length
  })
</script>

<section class="results-panel" in:fade={{duration: 200}}>
  <!-- 头部统计信息 -->
  <div class="panel-header">
    <div class="title-group">
      <div class="title-with-control">
        <h2>筛选结果</h2>
        <button class="collapse-btn" on:click={toggleSidebar}>
          <svelte:component this={isCollapsed ? ChevronRight : ChevronLeft} size={16} strokeWidth={1.5} />
        </button>
      </div>
      <div class="stats">
        <span class="stat">
          <span class="label">总数</span>
          <span class="value">{safeTotalStocks}</span>
        </span>
        <span class="divider"></span>
        <span class="stat">
          <span class="label">匹配</span>
          <span class="value highlight">{safeMatchedStocks}</span>
        </span>
        <span class="divider"></span>
        <span class="stat">
          <span class="label">匹配率</span>
          <span class="value">{matchRate}</span>
        </span>
      </div>
    </div>
  </div>

  <!-- 主内容区域 -->
  <div class="panel-content">
    <!-- 左侧列表 -->
    {#if isHovered}
      <div 
        class="stock-list expanded"
        on:mouseleave={handleMouseLeave}
        transition:fly={{x: -200, duration: 200}}
      >
        {#each safeSignals as stock, index}
          <div 
            class="stock-item" 
            class:selected={selectedStock === stock}
            role="button"
            tabindex="0"
            on:click={() => handleStockSelect(stock)}
            on:keydown={e => e.key === 'Enter' && handleStockSelect(stock)}
          >
            <div class="stock-index">{(index + 1).toString().padStart(2, '0')}</div>
            <div class="stock-content">
              <div class="stock-info">
                <div class="stock-name">
                  <span class="code">{stock.code}</span>
                  <span class="name">{stock.name}</span>
                </div>
                <div class="stock-price">
                  <span class="price">{formatNumber(stock.price)}</span>
                  <span class="change" class:up={stock.change > 0} class:down={stock.change < 0}>
                    {stock.change > 0 ? '+' : ''}{formatPercent(stock.change)}
                  </span>
                </div>
              </div>
            </div>
          </div>
        {/each}
      </div>
    {/if}
    <div 
      class="stock-list"
      class:collapsed={isCollapsed}
      style="width: {listWidth}px"
      on:mouseenter={handleMouseEnter}
    >
      {#each safeSignals as stock, index}
        <div 
          class="stock-item" 
          class:selected={selectedStock === stock}
          role="button"
          tabindex="0"
          on:click={() => handleStockSelect(stock)}
          on:keydown={e => e.key === 'Enter' && handleStockSelect(stock)}
        >
          <div class="stock-index">{(index + 1).toString().padStart(2, '0')}</div>
          <div class="stock-content">
            <div class="stock-info">
              <div class="stock-name">
                <span class="code">{stock.code}</span>
                <span class="name">{stock.name}</span>
              </div>
              <div class="stock-price">
                <span class="price">{formatNumber(stock.price)}</span>
                <span class="change" class:up={stock.change > 0} class:down={stock.change < 0}>
                  {stock.change > 0 ? '+' : ''}{formatPercent(stock.change)}
                </span>
              </div>
            </div>
          </div>
        </div>
      {/each}
    </div>

    <!-- 右侧内容区域 -->
    <div class="stock-detail">
      {#if selectedStock}
        <div class="detail-content">
          <div class="detail-body">
            <StockChart 
              code={selectedStock.code} 
              defaultVisibleRange={90}
              rightOffsetBars={45}
            />
          </div>
        </div>
      {:else}
        <div class="empty-state">
          <span>选择股票查看详情</span>
        </div>
      {/if}
    </div>
  </div>
</section>

<style>
  .results-panel {
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    background: white;
    margin-top: 16px;
    height: calc(100vh - 90px);
    max-height: calc(100vh - 90px);
    display: flex;
    flex-direction: column;
  }

  .panel-header {
    padding: 16px;
    border-bottom: 1px solid #e5e7eb;
  }

  .title-group {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .title-with-control {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  h2 {
    font-size: 16px;
    font-weight: 600;
    color: #111827;
    margin: 0;
  }

  .collapse-btn {
    width: 24px;
    height: 24px;
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.2s ease;
    color: #6b7280;
  }

  .collapse-btn:hover {
    background: #f8fafc;
    color: #111827;
    transform: scale(1.1);
  }

  .stats {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .stat {
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .label {
    font-size: 13px;
    color: #6b7280;
  }

  .value {
    font-size: 13px;
    font-weight: 500;
    color: #111827;
  }

  .value.highlight {
    color: #2563eb;
  }

  .divider {
    width: 1px;
    height: 12px;
    background: #e5e7eb;
  }

  .panel-content {
    flex: 1;
    display: flex;
    overflow: hidden;
    position: relative;
  }

  .stock-list {
    border-right: 1px solid #e5e7eb;
    overflow-y: auto;
    background: #f1f5f9;
    transition: width 0.2s ease;
  }

  .stock-list.expanded {
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    width: 320px;
    box-shadow: 4px 0 12px rgba(0, 0, 0, 0.1);
    background: #f1f5f9;
    z-index: 20;
  }

  .stock-list.collapsed .stock-index,
  .stock-list.collapsed .stock-price {
    display: none;
  }

  .stock-list.collapsed .stock-name {
    text-align: center;
    width: 100%;
  }

  .stock-list.collapsed .name {
    font-size: 11px;
  }

  .stock-list.collapsed .code {
    font-size: 12px;
  }

  .stock-item {
    height: 48px;
    padding: 0 16px;
    border-bottom: 1px solid #e2e8f0;
    cursor: pointer;
    transition: all 0.15s ease;
    display: flex;
    align-items: center;
    gap: 10px;
    background: #f8fafc;
    margin: 1px 1px;
    border-radius: 4px;
  }

  .stock-item:hover {
    background: #ffffff;
    transform: translateX(2px);
    box-shadow: 0 2px 4px rgba(0,0,0,0.05);
  }

  .stock-item.selected {
    background: #f0f7ff;
    border-left: 2px solid #2563eb;
    padding-left: 14px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.05);
  }

  .stock-index {
    font-family: 'SF Mono', monospace;
    font-size: 11px;
    font-weight: 500;
    color: #9ca3af;
    min-width: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .stock-content {
    flex: 1;
    min-width: 0;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 8px;
  }

  .stock-info {
    display: flex;
    align-items: center;
    gap: 12px;
    min-width: 0;
    flex: 1;
  }

  .stock-name {
    display: flex;
    flex-direction: column;
    gap: 2px;
    min-width: 0;
    flex: 1;
  }

  .code {
    font-size: 13px;
    font-weight: 500;
    color: #111827;
  }

  .name {
    font-size: 12px;
    color: #6b7280;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .stock-price {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-left: auto;
  }

  .price {
    font-size: 13px;
    font-weight: 500;
    color: #111827;
    white-space: nowrap;
  }

  .change {
    font-size: 12px;
    font-weight: 500;
    padding: 2px 8px;
    border-radius: 3px;
    white-space: nowrap;
  }

  .up {
    color: #10b981;
    background: #ecfdf5;
  }

  .down {
    color: #ef4444;
    background: #fef2f2;
  }

  .stock-detail {
    flex: 1;
    overflow-y: auto;
    background: white;
  }

  .detail-content {
    height: 100%;
    display: flex;
    flex-direction: column;
  }

  .detail-body {
    flex: 1;
    min-height: 0;
  }

  .empty-state {
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #9ca3af;
    font-size: 14px;
  }

  :global(.results-panel ::-webkit-scrollbar) {
    width: 4px;
  }

  :global(.results-panel ::-webkit-scrollbar-track) {
    background: #e2e8f0;
  }

  :global(.results-panel ::-webkit-scrollbar-thumb) {
    background-color: #cbd5e1;
    border-radius: 2px;
  }
</style> 