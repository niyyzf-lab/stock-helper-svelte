<script lang="ts">
  export let stock: {
    code: string;
    name: string;
    price: number;
    turnover: number;
    change: number;
    volume: number;
    amount: number;
    reason: string;
  };

  // 格式化数字
  function formatNumber(num: number, precision: number = 2): string {
    return num.toFixed(precision);
  }

  // 格式化成交量/金额（以万为单位）
  function formatLargeNumber(num: number): string {
    if (num >= 100000000) {
      return (num / 100000000).toFixed(2) + '亿';
    }
    return (num / 10000).toFixed(2) + '万';
  }
</script>

<div class="stock-card">
  <div class="card-header">
    <div class="stock-info">
      <h3 class="stock-name">{stock.name}</h3>
      <span class="stock-code">{stock.code}</span>
    </div>
    <div class="price-info">
      <span class="price">{formatNumber(stock.price)}</span>
      <span class="change" class:up={stock.change > 0} class:down={stock.change < 0}>
        {stock.change > 0 ? '+' : ''}{formatNumber(stock.change)}%
      </span>
    </div>
  </div>
  
  <div class="card-content">
    <div class="metrics">
      <div class="metric">
        <span class="label">换手率</span>
        <span class="value">{formatNumber(stock.turnover)}%</span>
      </div>
      <div class="metric">
        <span class="label">成交量</span>
        <span class="value">{formatLargeNumber(stock.volume)}</span>
      </div>
      <div class="metric">
        <span class="label">成交额</span>
        <span class="value">{formatLargeNumber(stock.amount)}</span>
      </div>
    </div>
    
    <div class="reason">
      <span class="reason-label">选股理由</span>
      <p class="reason-text">{stock.reason}</p>
    </div>
  </div>
</div>

<style>
  .stock-card {
    background: var(--background-color);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-lg);
    padding: 16px;
    transition: all 0.2s ease;
  }

  .stock-card:hover {
    border-color: var(--neutral-200);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 16px;
  }

  .stock-info {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .stock-name {
    font-size: var(--text-lg);
    font-weight: var(--font-medium);
    color: var(--text-primary);
    margin: 0;
  }

  .stock-code {
    font-size: var(--text-sm);
    color: var(--text-secondary);
  }

  .price-info {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 4px;
  }

  .price {
    font-size: var(--text-xl);
    font-weight: var(--font-medium);
    color: var(--text-primary);
  }

  .change {
    font-size: var(--text-sm);
    font-weight: var(--font-medium);
    padding: 2px 8px;
    border-radius: var(--radius-full);
  }

  .change.up {
    color: var(--success-600);
    background: var(--success-50);
  }

  .change.down {
    color: var(--error-600);
    background: var(--error-50);
  }

  .metrics {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 12px;
    margin-bottom: 16px;
  }

  .metric {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .label {
    font-size: var(--text-xs);
    color: var(--text-secondary);
  }

  .value {
    font-size: var(--text-sm);
    color: var(--text-primary);
    font-weight: var(--font-medium);
  }

  .reason {
    border-top: 1px solid var(--border-color);
    padding-top: 12px;
  }

  .reason-label {
    font-size: var(--text-xs);
    color: var(--text-secondary);
    margin-bottom: 4px;
    display: block;
  }

  .reason-text {
    font-size: var(--text-sm);
    color: var(--text-primary);
    margin: 0;
    line-height: 1.5;
  }
</style> 