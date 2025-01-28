<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  function handleSelectStrategy() {
    dispatch('select');
  }
</script>

<div class="empty-state">
  <div class="empty-content">
    <svg class="empty-icon" viewBox="0 0 24 24" width="32" height="32" stroke="currentColor" fill="none">
      <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z" stroke-width="1.5"/>
      <path d="M3.29 7L12 12l8.71-5M12 22V12" stroke-width="1.5"/>
    </svg>
    <h3>开始使用股票筛选</h3>
    <p>选择一个策略开始筛选股票</p>
    <button class="btn solid large" on:click={handleSelectStrategy}>
      选择策略
    </button>
  </div>
</div>

<style>
  .empty-state {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 48px;
    min-height: 400px;
  }

  .empty-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    text-align: center;
  }

  .empty-icon {
    color: #2563eb;
    margin-bottom: 24px;
  }

  h3 {
    font-size: 16px;
    font-weight: 600;
    color: #111827;
    margin: 0 0 8px;
  }

  p {
    font-size: 14px;
    color: #6b7280;
    margin: 0 0 24px;
  }

  .btn.solid.large {
    height: 40px;
    padding: 0 24px;
    font-size: 14px;
    color: white;
    background: #2563eb;
    border: none;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.15s;
  }

  .btn.solid.large:hover {
    background: #1d4ed8;
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(37, 99, 235, 0.2);
  }

  .btn.solid.large:active {
    transform: translateY(0);
  }
</style> 