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
    gap: 40px;
    height: 100%;
    max-height: 85vh;
    transform-origin: center;
    padding: 40px;
    min-width: 800px;
    max-width: 1200px;
    margin: 0 auto;
    position: relative;
  }

  .strategy-selector::before {
    content: '';
    position: absolute;
    inset: 0;
    background: radial-gradient(circle at top right, 
      color-mix(in srgb, var(--primary-500) 2%, transparent),
      transparent 60%
    );
    pointer-events: none;
  }

  /* 加载状态 */
  .loading-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 20px;
    min-height: 300px;
    color: var(--text-secondary);
    padding: 48px;
    text-align: center;
    position: relative;
  }

  .loading-icon {
    color: var(--primary-500);
    width: 48px;
    height: 48px;
    filter: drop-shadow(0 4px 8px color-mix(in srgb, var(--primary-500) 30%, transparent));
  }

  .spin {
    animation: spin 1.2s cubic-bezier(0.4, 0, 0.2, 1) infinite;
  }

  /* 错误状态 */
  .error-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 24px;
    min-height: 300px;
    text-align: center;
    transform-origin: center;
    padding: 48px;
    color: var(--text-secondary);
    position: relative;
  }

  .error-icon {
    color: var(--error-500);
    animation: shake 0.5s cubic-bezier(0.36, 0.07, 0.19, 0.97) both;
    width: 56px;
    height: 56px;
    margin: 0 auto 16px;
    filter: drop-shadow(0 4px 8px color-mix(in srgb, var(--error-500) 30%, transparent));
  }

  /* 策略网格 */
  .strategies-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 32px;
    margin-bottom: 32px;
    position: relative;
  }

  /* 策略卡片 */
  .strategy-card {
    position: relative;
    padding: 32px;
    background: color-mix(in srgb, var(--surface) 97%, transparent);
    border: 1px solid color-mix(in srgb, var(--border-color) 80%, transparent);
    border-radius: var(--radius-xl);
    cursor: pointer;
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
    display: flex;
    align-items: flex-start;
    gap: 24px;
    min-height: 180px;
    height: 180px;
    backdrop-filter: blur(8px);
    -webkit-backdrop-filter: blur(8px);
    touch-action: manipulation;
    user-select: none;
    -webkit-tap-highlight-color: transparent;
    overflow: hidden;
    box-sizing: border-box;
  }

  .strategy-card::before {
    content: '';
    position: absolute;
    inset: 0;
    background: radial-gradient(circle at top left,
      color-mix(in srgb, var(--primary-500) 3%, transparent),
      transparent 70%
    );
    opacity: 0;
    transition: opacity 0.4s ease;
  }

  .strategy-card:hover::before {
    opacity: 1;
  }

  .strategy-card:hover {
    background: color-mix(in srgb, var(--surface-variant) 97%, transparent);
    transform: translateY(-3px);
    box-shadow: var(--shadow-lg),
                0 0 0 1px color-mix(in srgb, var(--primary-500) 8%, transparent);
  }

  .strategy-card:active {
    transform: translateY(0) scale(0.98);
  }

  .strategy-card.selected {
    background: color-mix(in srgb, var(--primary-500) 8%, var(--surface));
    border: 2px solid var(--primary-500);
    padding: 31px;
    box-shadow: 0 0 0 1px var(--primary-500),
                0 0 0 4px color-mix(in srgb, var(--primary-500) 20%, transparent),
                0 8px 24px -8px color-mix(in srgb, var(--primary-500) 30%, transparent);
    transform: translateY(-3px);
  }

  :global(.dark) .strategy-card.selected {
    background: color-mix(in srgb, var(--primary-500) 15%, var(--surface));
    border-color: var(--primary-400);
    box-shadow: 0 0 0 1px var(--primary-400),
                0 0 0 4px color-mix(in srgb, var(--primary-400) 20%, transparent),
                0 8px 24px -8px color-mix(in srgb, var(--primary-400) 30%, transparent);
  }

  .strategy-icon {
    width: 56px;
    height: 56px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: color-mix(in srgb, var(--primary-500) 12%, var(--surface));
    border-radius: var(--radius-xl);
    color: var(--primary-500);
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
    flex-shrink: 0;
    box-shadow: inset 0 0 0 1px color-mix(in srgb, var(--primary-500) 15%, transparent),
                0 4px 12px -2px color-mix(in srgb, var(--primary-500) 20%, transparent);
    position: relative;
    z-index: 1;
  }

  .strategy-icon::before {
    content: '';
    position: absolute;
    inset: -2px;
    background: radial-gradient(circle at 30% 30%,
      color-mix(in srgb, var(--primary-500) 10%, transparent),
      transparent 70%
    );
    border-radius: inherit;
    opacity: 0;
    transition: opacity 0.4s ease;
  }

  .strategy-card:hover .strategy-icon::before {
    opacity: 1;
  }

  .strategy-card:hover .strategy-icon {
    transform: scale(1.1) rotate(5deg);
    background: color-mix(in srgb, var(--primary-500) 18%, var(--surface));
    color: var(--primary-600);
    box-shadow: inset 0 0 0 1px color-mix(in srgb, var(--primary-500) 20%, transparent),
                0 8px 16px -4px color-mix(in srgb, var(--primary-500) 25%, transparent);
  }

  .strategy-card.selected .strategy-icon {
    background: color-mix(in srgb, var(--primary-500) 25%, var(--surface));
    color: var(--primary-500);
    box-shadow: 0 0 0 2px var(--surface),
                0 0 0 4px color-mix(in srgb, var(--primary-500) 30%, transparent),
                0 8px 16px -4px color-mix(in srgb, var(--primary-500) 30%, transparent);
  }

  .strategy-info {
    flex: 1;
    min-width: 0;
    transform-origin: left;
    display: flex;
    flex-direction: column;
    gap: 16px;
    position: relative;
    z-index: 1;
    height: 120px;
  }

  .strategy-name {
    font-size: var(--text-xl);
    font-weight: var(--font-semibold);
    color: var(--text-primary);
    letter-spacing: -0.01em;
    text-shadow: 0 1px 2px color-mix(in srgb, var(--shadow-color) 10%, transparent);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .strategy-desc {
    font-size: var(--text-base);
    color: var(--text-secondary);
    line-height: 1.7;
    display: -webkit-box;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
    height: calc(1.7em * 2);
    margin: 0;
    padding: 0;
  }

  .check-icon {
    position: absolute;
    top: 28px;
    right: 28px;
    width: 28px;
    height: 28px;
    color: var(--primary-500);
    opacity: 0;
    transform: scale(0.5) rotate(-90deg);
    transition: all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
    background: var(--surface);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: var(--shadow-md),
                0 0 0 2px color-mix(in srgb, var(--primary-500) 15%, transparent),
                0 4px 12px -2px color-mix(in srgb, var(--primary-500) 20%, transparent);
    z-index: 2;
  }

  .strategy-card.selected .check-icon {
    opacity: 1;
    transform: scale(1) rotate(0);
  }

  /* 按钮组 */
  .actions {
    display: flex;
    justify-content: flex-end;
    gap: 20px;
    margin-top: auto;
    padding-top: 32px;
    border-top: 1px solid var(--border-color);
    position: relative;
  }

  .btn {
    height: 48px;
    padding: 0 32px;
    font-size: var(--text-base);
    font-weight: var(--font-semibold);
    border-radius: var(--radius-xl);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    min-width: 120px;
    letter-spacing: -0.01em;
    touch-action: manipulation;
    position: relative;
    overflow: hidden;
  }

  .btn::before {
    content: '';
    position: absolute;
    inset: 0;
    background: radial-gradient(circle at var(--x, 50%) var(--y, 50%),
      color-mix(in srgb, var(--primary-500) 20%, transparent),
      transparent 70%
    );
    opacity: 0;
    transition: opacity 0.3s ease;
  }

  .btn:hover::before {
    opacity: 0.8;
    --x: var(--hover-x, 50%);
    --y: var(--hover-y, 50%);
  }

  .btn.outline {
    color: var(--text-secondary);
    background: var(--surface);
    border: 1px solid var(--border-color);
  }

  .btn.outline:hover {
    color: var(--text-primary);
    border-color: var(--border-hover);
    background: var(--hover-bg);
    transform: translateY(-1px);
    box-shadow: var(--shadow-sm);
  }

  .btn.solid {
    color: white;
    background: var(--primary-500);
    border: none;
    box-shadow: 0 1px 3px color-mix(in srgb, var(--primary-500) 30%, transparent),
                0 0 0 1px color-mix(in srgb, var(--primary-500) 20%, transparent);
  }

  .btn.solid:hover:not(:disabled) {
    background: var(--primary-600);
    transform: translateY(-2px);
    box-shadow: 0 4px 12px color-mix(in srgb, var(--primary-500) 30%, transparent),
                0 0 0 1px color-mix(in srgb, var(--primary-500) 30%, transparent);
  }

  .btn.solid:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .btn:active:not(:disabled) {
    transform: translateY(0) scale(0.97);
  }

  /* 动画 */
  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  @keyframes shake {
    10%, 90% { transform: translate3d(-1px, 0, 0); }
    20%, 80% { transform: translate3d(2px, 0, 0); }
    30%, 50%, 70% { transform: translate3d(-4px, 0, 0); }
    40%, 60% { transform: translate3d(4px, 0, 0); }
  }

  /* 确保选中和悬停状态的过渡更加平滑 */
  .strategy-card,
  .strategy-card:hover,
  .strategy-card.selected,
  .strategy-card:active {
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1),
                padding 0s linear,
                border 0s linear;
  }
</style> 