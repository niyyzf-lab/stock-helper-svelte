<script lang="ts">
  import Header from '../components/Header.svelte'
</script>

<div class="layout">
  <div class="sidebar">
    <slot name="sidebar">
      <!-- 默认侧边栏内容 -->
    </slot>
  </div>
  <div class="header">
    <Header />
  </div>
  <div class="main">
    <slot>
      <!-- 默认主要内容 -->
    </slot>
  </div>
</div>

<style>
  .layout {
    display: grid;
    grid-template-columns: auto 1fr;
    grid-template-rows: 60px 1fr;
    width: 100%;
    height: 100vh;
    overflow: hidden;
  }
  .sidebar {
    grid-column: 1;
    grid-row: 1 / 3;
  }
  .header{
    grid-column: 2;
    grid-row: 1;
  }
  .main{
    grid-column: 2;
    grid-row: 2;
    overflow-y: auto;
  }
</style> 