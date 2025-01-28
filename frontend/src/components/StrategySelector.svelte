<script lang="ts">
    import { Check, ListFilter, BarChart3, TrendingUp, RefreshCw, AlertTriangle } from 'lucide-svelte';
  import { createEventDispatcher } from 'svelte';
    import { fade, fly, scale } from 'svelte/transition';
    import { elasticOut } from 'svelte/easing';


  interface Strategy {
    id: number;
    name: string;
    description: string;
    filePath: string;
  }

  export let strategies: Strategy[] = [];
  export let selectedStrategyId: number | null = null;
  export let loading = false;
  export let error: string | null = null;

  const dispatch = createEventDispatcher();

  // 根据策略ID获取对应图标
  function getStrategyIcon(id: number) {
    switch(id) {
      case 1: return ListFilter;
      case 2: return BarChart3;
      default: return TrendingUp;
    }
  }

  function selectStrategy(id: number) {
    selectedStrategyId = id;
  }

  function confirmStrategy() {
    const strategy = strategies.find(s => s.id === selectedStrategyId);
    if (strategy) {
      dispatch('confirm', strategy);
    }
  }
</script>

<div class="strategy-selector" in:fade={{duration: 200}}>
  {#if loading}
    <div class="loading-state" in:fade={{duration: 200}}>
      <div class="loading-icon spin">
        <RefreshCw size={24} />
      </div>
      <span>加载中...</span>
    </div>
  {:else if error}
    <div class="error-state" in:fly={{y: 20, duration: 300}}>
      <div class="error-icon">
        <AlertTriangle size={24} />
      </div>
      <div class="error-content">
        <h3>加载失败</h3>
        <p>{error}</p>
        <button class="btn outline" on:click={() => dispatch('retry')}>
          重试
        </button>
      </div>
    </div>
  {:else}
    <div class="strategies-grid">
      {#each strategies as strategy, i}
        <button 
          class="strategy-card" 
          class:selected={selectedStrategyId === strategy.id}
          on:click={() => selectStrategy(strategy.id)}
          in:scale={{
            start: 0.9,
            duration: 300,
            delay: i * 50,
            easing: elasticOut
          }}
        >
          <div class="strategy-icon" in:scale={{start: 0.8, duration: 300, delay: i * 50 + 50}}>
            <svelte:component this={getStrategyIcon(strategy.id)} size={20} strokeWidth={1.5} />
          </div>
          <div class="strategy-info">
            <span class="strategy-name" in:fly={{y: 10, duration: 300, delay: i * 50 + 100}}>{strategy.name}</span>
            <span class="strategy-desc" in:fly={{y: 10, duration: 300, delay: i * 50 + 150}}>{strategy.description}</span>
          </div>
          <div class="check-icon" in:scale={{start: 0.5, duration: 300, delay: i * 50 + 200}}>
            <Check size={16} />
          </div>
        </button>
      {/each}
    </div>

    <div class="actions" in:fly={{y: 20, duration: 300, delay: 200}}>
      <button class="btn outline" on:click={() => dispatch('cancel')}>
        取消
      </button>
      <button 
        class="btn solid" 
        disabled={!selectedStrategyId} 
        on:click={confirmStrategy}
      >
        确认
      </button>
    </div>
  {/if}
</div>

<style>
  /* 容器样式 */
  .strategy-selector {
    display: flex;
    flex-direction: column;
    gap: 24px;
    height: 100%;
    max-height: 80vh;
    transform-origin: center;
  }

  /* 加载状态 */
  .loading-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 16px;
    min-height: 300px;
    color: #6b7280;
  }

  .loading-icon {
    color: #2563eb;
  }

  .spin {
    animation: spin 1s linear infinite;
  }

  /* 错误状态 */
  .error-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 20px;
    min-height: 300px;
    text-align: center;
    transform-origin: center;
  }

  .error-icon {
    color: #ef4444;
    animation: shake 0.5s cubic-bezier(0.36, 0.07, 0.19, 0.97) both;
  }

  /* 策略网格 */
  .strategies-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
    gap: 20px;
    width: 100%;
    padding: 24px;
    padding-bottom: 0;
    overflow-y: auto;
    min-height: 0;
    /* 自定义滚动条样式 */
    scrollbar-width: thin;
    scrollbar-color: #e5e7eb transparent;
  }

  .strategies-grid::-webkit-scrollbar {
    width: 6px;
  }

  .strategies-grid::-webkit-scrollbar-track {
    background: transparent;
    margin: 4px;
  }

  .strategies-grid::-webkit-scrollbar-thumb {
    background-color: #e5e7eb;
    border-radius: 3px;
    border: 1px solid transparent;
    background-clip: padding-box;
  }

  .strategies-grid::-webkit-scrollbar-thumb:hover {
    background-color: #d1d5db;
  }

  /* 策略卡片 */
  .strategy-card {
    position: relative;
    display: flex;
    align-items: flex-start;
    gap: 20px;
    padding: 24px;
    background: white;
    border: 1px solid #e5e7eb;
    border-radius: 12px;
    cursor: pointer;
    text-align: left;
    transform-origin: center;
    transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
    width: 100%;
    min-height: 140px;
    overflow: hidden;
  }

  .strategy-card:hover {
    transform: translateY(-2px) scale(1.02);
    border-color: #93c5fd;
    box-shadow: 0 8px 24px rgba(37, 99, 235, 0.12);
    background: linear-gradient(to right bottom, #ffffff, #f8fafc);
  }

  .strategy-card.selected {
    background: linear-gradient(to right bottom, #f0f7ff, #e0e7ff);
    border-color: #2563eb;
    box-shadow: 0 8px 24px rgba(37, 99, 235, 0.12);
  }

  .strategy-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 52px;
    height: 52px;
    background: #f3f4f6;
    border-radius: 14px;
    color: #2563eb;
    flex-shrink: 0;
    transform-origin: center;
    transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  }

  .strategy-card:hover .strategy-icon {
    transform: scale(1.1) rotate(-5deg);
    background: #e0e7ff;
  }

  .strategy-card.selected .strategy-icon {
    background: white;
    transform: scale(1.1) rotate(-5deg);
    box-shadow: 0 4px 12px rgba(37, 99, 235, 0.2);
  }

  .strategy-info {
    flex: 1;
    min-width: 0;
    transform-origin: left;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .strategy-name {
    display: block;
    font-size: 16px;
    font-weight: 600;
    color: #111827;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .strategy-desc {
    display: -webkit-box;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;
    overflow: hidden;
    font-size: 14px;
    color: #6b7280;
    line-height: 1.6;
  }

  .check-icon {
    position: absolute;
    top: 24px;
    right: 24px;
    width: 20px;
    height: 20px;
    color: #2563eb;
    opacity: 0;
    transform: scale(0.5);
    transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
    background: white;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 2px 8px rgba(37, 99, 235, 0.2);
  }

  .strategy-card.selected .check-icon {
    opacity: 1;
    transform: scale(1);
  }

  /* 按钮组 */
  .actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    padding: 24px;
    padding-top: 20px;
    border-top: 1px solid #e5e7eb;
    background: white;
    flex-shrink: 0;
  }

  .btn {
    height: 38px;
    padding: 0 20px;
    font-size: 14px;
    font-weight: 500;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .btn.outline {
    color: #4b5563;
    background: white;
    border: 1px solid #e5e7eb;
  }

  .btn.outline:hover {
    color: #111827;
    border-color: #d1d5db;
    background: #f9fafb;
    transform: translateY(-1px);
  }

  .btn.solid {
    color: white;
    background: #2563eb;
    border: none;
  }

  .btn.solid:hover:not(:disabled) {
    background: #1d4ed8;
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(37, 99, 235, 0.2);
  }

  .btn.solid:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .btn:active:not(:disabled) {
    transform: translateY(0) scale(0.98);
  }

  /* 动画 */
  @keyframes spin {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
  }

  @keyframes shake {
    10%, 90% { transform: translate3d(-1px, 0, 0); }
    20%, 80% { transform: translate3d(2px, 0, 0); }
    30%, 50%, 70% { transform: translate3d(-4px, 0, 0); }
    40%, 60% { transform: translate3d(4px, 0, 0); }
  }
</style> 