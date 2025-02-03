<script lang="ts">
  import { fly } from 'svelte/transition'
  import { quintOut } from 'svelte/easing'
  import type { TestResult } from 'src/types/stock'
  import { 
    TrendingUpIcon, 
    TrendingDownIcon, 
    ChevronRightIcon 
  } from 'lucide-svelte'
  import { createEventDispatcher } from 'svelte'

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

  // 评语计算
  $: comment = (() => {
    const accuracyNum = parseFloat(accuracy.toString())
    if (accuracyNum === 100) return '完美！您的预测非常准确！'
    if (accuracyNum >= 80) return '太棒了！您展现出了优秀的判断力'
    if (accuracyNum >= 60) return '不错！您的预测方向大体正确'
    if (accuracyNum >= 40) return '继续加油！市场总是充满变数'
    return '别灰心！这是提升的好机会'
  })()

  // 添加事件分发器
  const dispatch = createEventDispatcher()

  function handleReset() {
    dispatch('reset')
  }
</script>

<div class="result-panel" in:fly={{y: 20, duration: 400, easing: quintOut}}>
  <div class="content">
    <!-- 巨大的分数显示 -->
    <div class="score-section">
      <div class="score">{accuracy}<span class="percent">%</span></div>
      <div class="score-label">预测准确率</div>
      <div class="comment">{comment}</div>
    </div>

    <!-- 指标统计 -->
    <div class="stats-grid">
      <div class="stat-item">
        <div class="stat-value">{results.length}</div>
        <div class="stat-label">测试数量</div>
      </div>
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

    <!-- 返回按钮 -->
    <button class="restart-btn" on:click={handleReset}>
      重新测试
      <ChevronRightIcon size="16" />
    </button>
  </div>
</div>

<style>
  .result-panel {
    height: 100%;
    background: var(--neutral-50);
    border-radius: var(--radius-lg);
    box-shadow: 0 0 0 1px var(--border-color);
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .content {
    max-width: 600px;
    width: 100%;
    padding: 48px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 48px;
  }

  .score-section {
    text-align: center;
  }

  .score {
    font-size: 8rem;
    font-weight: 700;
    color: var(--primary-600);
    line-height: 1;
    margin-bottom: 16px;
  }

  .percent {
    font-size: 4rem;
    margin-left: 8px;
  }

  .score-label {
    font-size: 1.5rem;
    color: var(--text-secondary);
    margin-bottom: 24px;
  }

  .comment {
    font-size: 1.25rem;
    color: var(--primary-600);
    font-weight: 500;
  }

  .stats-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 24px;
    width: 100%;
  }

  .stat-item {
    text-align: center;
    padding: 24px 16px;
    border-radius: var(--radius-lg);
    background: white;
    border: 1px solid var(--border-color);
    transition: all 0.2s;
  }

  .stat-item:hover {
    border-color: var(--primary-200);
    transform: translateY(-2px);
  }

  .stat-value {
    font-size: 2rem;
    font-weight: 600;
    margin-bottom: 8px;
  }

  .stat-value.up { color: #ef4444; }
  .stat-value.down { color: #10b981; }

  .stat-label {
    font-size: 0.875rem;
    color: var(--text-secondary);
    font-weight: 500;
  }

  .restart-btn {
    padding: 16px 32px;
    background: var(--primary-600);
    color: white;
    border: none;
    border-radius: var(--radius-lg);
    font-size: 1.125rem;
    font-weight: 500;
    display: flex;
    align-items: center;
    gap: 8px;
    transition: all 0.2s;
  }

  .restart-btn:hover {
    background: var(--primary-700);
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(var(--primary-500-rgb), 0.2);
  }

  .restart-btn:active {
    transform: translateY(0);
  }

  @media (max-width: 768px) {
    .content {
      padding: 32px;
      gap: 32px;
    }

    .score {
      font-size: 6rem;
    }

    .percent {
      font-size: 3rem;
    }

    .stats-grid {
      grid-template-columns: repeat(2, 1fr);
    }
  }
</style> 