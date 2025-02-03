<script lang="ts">
  import Header from '../components/Header.svelte'
</script>

<div class="layout">
  <div class="sidebar">
    <slot name="sidebar">
      <!-- 默认侧边栏内容 -->
    </slot>
  </div>
  <div class="content">
    <div class="header">
      <Header />
    </div>
    <div class="main">
      <slot>
        <!-- 默认主要内容 -->
      </slot>
    </div>
  </div>
</div>

<style>
  .layout {
    display: flex;
    width: 100%;
    height: 100vh;
    overflow: hidden;
  }

  .sidebar {
    position: relative;
    z-index: 2; /* 提高侧边栏的层级 */
  }

  .content {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-width: 0;
    z-index: 1;
  }

  .header {
    height: 60px;
    flex-shrink: 0;
    
  }

  .main {
    flex: 1;
    overflow-y: auto;
    min-height: 0;
  }
</style> 