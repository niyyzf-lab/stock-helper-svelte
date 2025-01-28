<script lang="ts">
  import { link } from 'svelte-spa-router'
</script>

<div class="overview-page">
  <div class="page-header">
    <h1>概览</h1>
    <p class="description">查看市场整体情况和关键指标</p>
  </div>

  <div class="content-area">
    <div class="test-content">
      <p>这是概览页面的内容</p>
      <nav class="test-nav">
        <a href="#/filter" use:link>去筛选页面</a>
        <a href="#/analysis" use:link>去分析页面</a>
      </nav>
    </div>
  </div>
</div>

<style>
  .overview-page {
    min-height: 100%;
  }

  .page-header {
    margin-bottom: 32px;
  }

  h1 {
    font-size: var(--text-xl);
    font-weight: var(--font-semibold);
    color: var(--text-primary);
    margin-bottom: 8px;
  }

  .description {
    color: var(--text-secondary);
    font-size: var(--text-sm);
  }

  .content-area {
    background: white;
    border-radius: var(--radius-lg);
    padding: 24px;
    box-shadow: var(--shadow-sm);
    border: 1px solid var(--neutral-100);
  }

  .test-nav {
    margin-top: 16px;
    display: flex;
    gap: 16px;
  }

  a {
    color: var(--primary-500);
    text-decoration: none;
  }

  a:hover {
    text-decoration: underline;
  }
</style> 