<script lang="ts">
  import Modal from '../Modal.svelte'
  import type { Stock } from '../../types/stock'

  export let show = false
  export let stocks: Stock[] = []
  export let onClose = () => {}
  export let onStart = () => {}

  function formatDate(dateStr: string) {
    if (!dateStr) return '未知时间'
    try {
      const date = new Date(dateStr)
      return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit'
      })
    } catch {
      return '无效日期'
    }
  }
</script>

<Modal 
  show={show}
  title="测试股票列表"
  on:close={onClose}
  class_="stock-list-modal"
>
  <div class="stock-list-container">
    <div class="header-section slide-up">
      <div class="info-text">
        系统已为您随机选择以下股票进行测试：
      </div>
      <div class="stock-count">
        共 <span>{stocks.length}</span> 只股票
      </div>
    </div>

    <div class="stock-list">
      {#each stocks as stock, i}
        <div class="stock-item slide-up" style="animation-delay: {i * 0.05}s">
          <div class="stock-main">
            <div class="stock-index">{i + 1}</div>
            <div class="stock-info">
              <div class="stock-name">{stock.mc}</div>
              <div class="stock-code">{stock.dm}</div>
            </div>
          </div>
          <div class="stock-date">
            <span class="date-label">测试日期</span>
            <span class="date-value">{formatDate(stock.testDate)}</span>
          </div>
        </div>
      {/each}
    </div>

    <button class="confirm-btn fade-in" on:click={onStart}>
      开始测试
    </button>
  </div>
</Modal>

<style>
  .stock-list-container {
    padding: 24px;
    background: #fff;
  }

  .header-section {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }

  .info-text {
    font-size: 0.875rem;
    color: var(--text-secondary);
  }

  .stock-count {
    font-size: 0.8125rem;
    color: var(--text-secondary);
  }

  .stock-count span {
    color: var(--primary-500);
    font-weight: 600;
    margin: 0 1px;
  }

  .stock-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
    max-height: 400px;
    overflow-y: auto;
    padding: 4px;
    margin: 0 -4px;
  }

  .stock-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 16px;
    background: #f8fafc;
    border-radius: 10px;
    transition: all 0.2s ease;
  }

  .stock-item:hover {
    background: #f1f5f9;
    transform: translateX(2px);
  }

  .stock-main {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .stock-index {
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: white;
    color: var(--primary-500);
    font-size: 0.8125rem;
    font-weight: 600;
    border-radius: 6px;
    border: 1px solid var(--primary-200);
  }

  .stock-info {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .stock-name {
    font-size: 0.9375rem;
    font-weight: 500;
    color: var(--text-primary);
  }

  .stock-code {
    font-family: var(--font-mono);
    font-size: 0.8125rem;
    color: var(--primary-500);
  }

  .stock-date {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 2px;
  }

  .date-label {
    font-size: 0.75rem;
    color: var(--text-secondary);
  }

  .date-value {
    font-size: 0.8125rem;
    color: var(--text-primary);
    font-weight: 500;
  }

  .confirm-btn {
    display: block;
    width: 100%;
    margin-top: 24px;
    padding: 12px;
    background: var(--primary-500);
    color: white;
    border: none;
    border-radius: 8px;
    font-size: 0.9375rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
  }

  .confirm-btn:hover {
    background: var(--primary-600);
    transform: translateY(-1px);
  }

  /* 动画 */
  .slide-up {
    animation: slideUp 0.4s ease-out forwards;
    opacity: 0;
  }

  .fade-in {
    animation: fadeIn 0.4s ease-out forwards;
    animation-delay: 0.2s;
    opacity: 0;
  }

  @keyframes slideUp {
    from {
      opacity: 0;
      transform: translateY(12px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }

  :global(.modal-content.stock-list-modal) {
    max-width: 480px !important;
  }

  /* 滚动条样式 */
  .stock-list::-webkit-scrollbar {
    width: 8px;
  }

  .stock-list::-webkit-scrollbar-track {
    background: transparent;
  }

  .stock-list::-webkit-scrollbar-thumb {
    background: #e2e8f0;
    border-radius: 4px;
  }

  .stock-list::-webkit-scrollbar-thumb:hover {
    background: #cbd5e1;
  }
</style> 