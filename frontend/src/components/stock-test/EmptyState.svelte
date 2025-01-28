<script lang="ts">
  import { Brain, ChartLine, Target } from 'lucide-svelte'

  export let loading = false
  export let progress = {
    current: 0,
    total: 0,
    currentStock: '',
    stage: ''
  }
  export let onStart: () => void
</script>

<section class="empty-state">
  <div class="empty-content animate-in">
    <div class="hero-section">
      <div class="icon-container fade-in">
        <Brain size={72} strokeWidth={1} />
      </div>
      <h2 class="slide-up">交易能力测试</h2>
      <p class="slide-up">通过真实历史数据，科学评估您的市场判断能力</p>
    </div>

    <div class="features">
      <div class="feature slide-in-left">
        <Target size={20} strokeWidth={1.5} />
        <span>随机选择优质标的</span>
      </div>
      <div class="feature slide-in-right">
        <ChartLine size={20} strokeWidth={1.5} />
        <span>真实历史行情数据</span>
      </div>
    </div>

    <button class="start-button fade-in-up" on:click={onStart} disabled={loading}>
      {#if loading}
        <div class="loading">
          <div class="spinner" />
          <div class="status">
            <span>{progress.stage || '准备中...'}</span>
            {#if progress.currentStock}
              <span class="detail">{progress.current}/{progress.total} - {progress.currentStock}</span>
            {/if}
          </div>
        </div>
      {:else}
        开始测试
      {/if}
    </button>
  </div>
</section>

<style>
  .empty-state {
    min-height: 100%;
    width: 100%;
    display: grid;
    place-items: center;
    padding: 2rem;
    position: relative;
    overflow: hidden;
    border-left: 2px dashed rgba(255, 255, 255, 0.3);
  }

  /* 移除内部虚线边框 */
  .empty-state::before {
    display: none;
  }

  /* 移除网格背景 */
  .empty-state::after {
    display: none;
  }

  .empty-content {
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 4rem;
    max-width: 480px;
    width: 100%;
    padding: 2rem;
    z-index: 1;
  }

  .hero-section {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    gap: 1.5rem;
  }

  .icon-container {
    color: var(--primary-600);
    animation: 
      iconEntrance 0.8s ease-out forwards,
      float 6s ease-in-out infinite 0.8s;
  }

  @keyframes iconEntrance {
    0% {
      opacity: 0;
      transform: translateY(-20px) scale(0.9);
    }
    100% {
      opacity: 1;
      transform: translateY(0) scale(1);
    }
  }

  @keyframes float {
    0%, 100% {
      transform: translateY(0);
    }
    50% {
      transform: translateY(-10px);
    }
  }

  h2 {
    font-size: 2rem;
    font-weight: 600;
    color: var(--text-primary);
    letter-spacing: -0.02em;
    margin: 0;
  }

  .hero-section p {
    font-size: 1rem;
    color: var(--text-secondary);
    line-height: 1.6;
    margin: 0;
  }

  .features {
    display: flex;
    gap: 2rem;
    color: var(--text-secondary);
  }

  .feature {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    font-size: 0.9375rem;
  }

  .feature :global(svg) {
    color: var(--primary-500);
  }

  .start-button {
    min-width: 180px;
    height: 56px;
    padding: 0 2rem;
    border-radius: 16px;
    border: none;
    background: var(--primary-600);
    color: white;
    font-size: 1rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .start-button:hover {
    background: var(--primary-700);
    transform: translateY(-1px);
  }

  .start-button:disabled {
    opacity: 0.7;
    cursor: wait;
  }

  .loading {
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  .spinner {
    width: 20px;
    height: 20px;
    border: 2px solid rgba(255, 255, 255, 0.3);
    border-top-color: white;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  .status {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 0.25rem;
    font-size: 0.875rem;
  }

  .detail {
    opacity: 0.8;
    font-size: 0.8125rem;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  /* 基础动画关键帧 */
  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }

  @keyframes slideUp {
    from { 
      opacity: 0;
      transform: translateY(20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  @keyframes slideInLeft {
    from {
      opacity: 0;
      transform: translateX(-20px);
    }
    to {
      opacity: 1;
      transform: translateX(0);
    }
  }

  @keyframes slideInRight {
    from {
      opacity: 0;
      transform: translateX(20px);
    }
    to {
      opacity: 1;
      transform: translateX(0);
    }
  }

  /* 动画类 */
  .fade-in {
    opacity: 0;
    animation: fadeIn 0.6s ease-out forwards;
    animation-delay: 0.2s;
  }

  .slide-up {
    opacity: 0;
    animation: slideUp 0.6s ease-out forwards;
  }

  .slide-up:nth-child(2) {
    animation-delay: 0.3s;
  }

  .slide-up:nth-child(3) {
    animation-delay: 0.4s;
  }

  .slide-in-left {
    opacity: 0;
    animation: slideInLeft 0.6s ease-out forwards;
    animation-delay: 0.5s;
  }

  .slide-in-right {
    opacity: 0;
    animation: slideInRight 0.6s ease-out forwards;
    animation-delay: 0.5s;
  }

  .fade-in-up {
    opacity: 0;
    animation: slideUp 0.6s ease-out forwards;
    animation-delay: 0.7s;
  }
</style> 

