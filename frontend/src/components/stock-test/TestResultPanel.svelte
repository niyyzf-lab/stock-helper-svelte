<script lang="ts">
  import { fly } from 'svelte/transition'
  import { quintOut } from 'svelte/easing'
  import type { TestResult } from 'src/types/stock'
  import { Trophy, Target, TrendingUp, TrendingDown, ChevronRight } from 'lucide-svelte'

  export let results: TestResult[] = []

  // 计算正确率
  $: accuracy = results.length > 0 
    ? (results.filter(r => r.correct).length / results.length * 100).toFixed(1) 
    : 0

  // 计算平均涨跌幅
  $: avgPriceChange = results.length > 0
    ? (results.reduce((sum, r) => sum + Math.abs(r.priceChange), 0) / results.length * 100).toFixed(1)
    : 0

  // 按方向统计
  $: directionStats = results.reduce((acc, r) => {
    acc[r.actualDirection] = (acc[r.actualDirection] || 0) + 1
    return acc
  }, {} as Record<string, number>)

  // 添加评语计算
  $: comment = (() => {
    const accuracyNum = parseFloat(accuracy)
    if (accuracyNum === 100) return '完美！您的预测非常准确！'
    if (accuracyNum >= 80) return '太棒了！您展现出了优秀的判断力'
    if (accuracyNum >= 60) return '不错！您的预测方向大体正确'
    if (accuracyNum >= 40) return '继续加油！市场总是充满变数'
    return '别灰心！这是提升的好机会'
  })()
</script>

<div class="result-panel" in:fly={{y: 20, duration: 400, easing: quintOut}}>
  <!-- 主要内容区域 -->
  <div class="main-section">
    <div class="records-header">
      <h3 class="records-title">详细记录</h3>
    </div>

    <div class="records-list">
      {#each results as result, index}
        <div class="record-card" class:correct={result.correct}>
          <div class="record-header">
            <div class="record-stock">
              <div class="stock-code">{result.stockCode}</div>
              <div class="stock-name">{result.stockName}</div>
            </div>
            <div class="record-status">
              {#if result.correct}
                <div class="status-icon">✓</div>
              {:else}
                <div class="status-icon">✕</div>
              {/if}
            </div>
          </div>

          <div class="record-predictions">
            <div class="prediction">
              <span class="label">预测</span>
              <span class="value {result.direction}">
                {#if result.direction === 'up'}
                  <TrendingUp size={14} />上涨
                {:else if result.direction === 'down'}
                  <TrendingDown size={14} />下跌
                {:else}
                  ↔️震荡
                {/if}
              </span>
            </div>
            <div class="prediction">
              <span class="label">实际</span>
              <span class="value {result.actualDirection}">
                {#if result.actualDirection === 'up'}
                  <TrendingUp size={14} />上涨
                {:else if result.actualDirection === 'down'}
                  <TrendingDown size={14} />下跌
                {:else}
                  ↔️震荡
                {/if}
                <span class="change">
                  ({(result.priceChange * 100).toFixed(1)}%)
                </span>
              </span>
            </div>
          </div>
        </div>
      {/each}
    </div>
  </div>

  <!-- 右侧概览区域 -->
  <div class="overview-section">
    <!-- 标题区域 -->
    <div class="overview-header">
      <div class="header-content">
        <h2 class="title">测试结果</h2>
        <p class="subtitle">共测试 {results.length} 只股票</p>
      </div>
      <button class="restart-btn" on:click>
        重新测试
        <ChevronRight size={16} />
      </button>
    </div>

    <!-- 主要统计区域 -->
    <div class="stats-container">
      <!-- 准确率 -->
      <div class="accuracy-section">
        <div class="accuracy-value">{accuracy}%</div>
        <div class="accuracy-label">预测准确率</div>
      </div>

      <!-- 统计数据 -->
      <div class="stats-grid">
        <div class="stat-item">
          <div class="stat-value up">{directionStats.up || 0}</div>
          <div class="stat-label">上涨</div>
        </div>
        <div class="stat-item">
          <div class="stat-value down">{directionStats.down || 0}</div>
          <div class="stat-label">下跌</div>
        </div>
        <div class="stat-item">
          <div class="stat-value">{avgPriceChange}%</div>
          <div class="stat-label">平均波动</div>
        </div>
      </div>
    </div>

    <!-- 评语区域 -->
    <div class="comment-section">
      <div class="comment-text">{comment}</div>
    </div>
  </div>
</div>

<style>
  .result-panel {
    height: 100%;
    display: grid;
    grid-template-columns: 1fr 400px;
    background: var(--neutral-50);
    overflow: hidden;
    border-radius: var(--radius-lg);
    box-shadow: 0 0 0 1px var(--border-color);
  }

  /* 主要内容区域 */
  .main-section {
    display: flex;
    flex-direction: column;
    min-height: 0;
    margin-right: 32px; /* 只保留与侧边栏的间距 */
  }

  .records-header {
    margin-bottom: 24px;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .records-title {
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--text-primary);
    display: flex;
    align-items: center;
    gap: 8px;
  }

  /* 移除下划线装饰 */
  .records-header::after {
    display: none;
  }

  .records-list {
    overflow-y: auto;
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 24px;
  }

  /* 美化滚动条 */
  .records-list::-webkit-scrollbar {
    width: 8px;
  }

  .records-list::-webkit-scrollbar-track {
    background: transparent;
  }

  .records-list::-webkit-scrollbar-thumb {
    background: var(--neutral-200);
    border-radius: 4px;
  }

  .records-list::-webkit-scrollbar-thumb:hover {
    background: var(--neutral-300);
  }

  .record-card {
    background: white;
    border-radius: var(--radius-lg);
    padding: 24px;
    display: flex;
    flex-direction: column;
    gap: 16px;
    border: 1px solid var(--border-color);
    position: relative;
    transition: all 0.2s;
  }

  .record-card:hover {
    border-color: var(--primary-200);
    box-shadow: 0 4px 12px rgba(var(--primary-500-rgb), 0.08);
  }

  .record-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    padding-bottom: 16px;
    border-bottom: 1px solid var(--border-color);
  }

  .record-stock {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .stock-code {
    font-family: var(--font-mono);
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--primary-600);
    letter-spacing: -0.02em;
  }

  .stock-name {
    font-size: 0.875rem;
    color: var(--text-secondary);
  }

  .record-predictions {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .prediction {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 8px 0;
  }

  .prediction .label {
    font-size: 0.875rem;
    color: var(--text-secondary);
    font-weight: 500;
    min-width: 40px; /* 固定标签宽度 */
  }

  .prediction .value {
    display: flex;
    align-items: center;
    gap: 6px;
    font-weight: 500;
  }

  .prediction .value.up { color: #ef4444; }
  .prediction .value.down { color: #10b981; }
  .prediction .value.shock { color: #6366f1; }

  .change {
    margin-left: 4px;
    color: var(--text-secondary);
    font-size: 0.75rem;
  }

  .status-icon {
    width: 28px;
    height: 28px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1rem;
    font-weight: 600;
    position: relative;
    overflow: hidden;
  }

  .record-card.correct .status-icon {
    color: #10b981;
  }

  .record-card:not(.correct) .status-icon {
    color: #ef4444;
  }

  /* 移除背景色,改用边框 */
  .status-icon::before {
    content: '';
    position: absolute;
    inset: 0;
    border: 2px solid currentColor;
    border-radius: 50%;
    opacity: 0.2;
  }

  /* 右侧概览区域 */
  .overview-section {
    background: white;
    padding: 32px;
    display: flex;
    flex-direction: column;
    gap: 32px;
    border-left: 1px solid var(--border-color);
  }

  .overview-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-bottom: 24px;
    border-bottom: 1px solid var(--border-color);
  }

  .header-content {
    flex: 1;
  }

  .title {
    font-size: 1.5rem;
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: 4px;
  }

  .subtitle {
    font-size: 0.875rem;
    color: var(--text-secondary);
  }

  .restart-btn {
    padding: 12px 24px;
    background: var(--primary-600);
    color: white;
    border: none;
    border-radius: var(--radius-lg);
    font-size: 1rem;
    font-weight: 500;
    display: flex;
    align-items: center;
    gap: 8px;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .restart-btn:hover {
    background: var(--primary-700);
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(var(--primary-500-rgb), 0.2);
  }

  .restart-btn:active {
    transform: translateY(0);
  }

  .accuracy-section {
    text-align: center;
    padding: 32px 0;
    position: relative;
    border-top: 1px solid var(--border-color);
    border-bottom: 1px solid var(--border-color);
  }

  .accuracy-value {
    font-size: 4rem;
    font-weight: 700;
    color: var(--primary-600);
    line-height: 1;
    margin-bottom: 12px;
  }

  .accuracy-label {
    font-size: 1rem;
    color: var(--text-secondary);
    font-weight: 500;
    position: relative;
    z-index: 1;
  }

  .stats-container {
    padding: 0;
    border: none;
  }

  .stats-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 24px;
  }

  .stat-item {
    text-align: center;
    padding: 16px;
    border-radius: var(--radius-lg);
    border: 1px solid var(--border-color);
    transition: all 0.2s;
  }

  .stat-item:hover {
    border-color: var(--primary-200);
  }

  .stat-value {
    font-size: 1.75rem;
    font-weight: 600;
    margin-bottom: 8px;
  }

  .stat-label {
    font-size: 0.875rem;
    font-weight: 500;
  }

  .comment-section {
    margin-top: auto;
    padding: 24px;
    border-top: 1px solid var(--border-color);
    text-align: center;
  }

  .comment-text {
    font-size: 1.125rem;
    line-height: 1.6;
    color: var(--primary-600);
    font-weight: 500;
  }

  /* 响应式调整 */
  @media (max-height: 800px) {
    .accuracy-section {
      padding: 24px 0;
    }

    .accuracy-value {
      font-size: 3.5rem;
    }

    .stat-value {
      font-size: 1.5rem;
    }
  }

  @media (max-height: 700px) {
    .overview-section {
      gap: 16px;
    }

    .accuracy-section {
      padding: 16px 0;
    }

    .accuracy-value {
      font-size: 3rem;
      margin-bottom: 8px;
    }

    .stats-grid {
      gap: 12px;
    }

    .stat-item {
      padding: 12px;
    }
  }

  @media (min-width: 1440px) {
    .records-list {
      grid-template-columns: repeat(3, 1fr);
    }
  }

  @media (max-width: 1280px) {
    .records-list {
      grid-template-columns: repeat(2, 1fr);
    }
  }

  @media (max-width: 768px) {
    .records-list {
      grid-template-columns: 1fr;
    }
  }
</style> 