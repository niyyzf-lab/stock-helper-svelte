<script lang="ts">
  import { fade, scale } from 'svelte/transition'
  
  export let show = false
  export let title = ''
  export let class_ = '' // 添加 class 属性用于自定义样式

  function close() {
    show = false
  }

  // 处理 ESC 键关闭
  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape' && show) {
      close()
    }
  }
</script>

<svelte:window on:keydown={handleKeydown} />

{#if show}
  <div 
    class="modal-wrapper"
    transition:fade={{ duration: 200 }}
  >
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div 
      class="modal-overlay"
      on:click={close}
      transition:fade={{ duration: 150 }}
    />
    
    <div 
      class="modal-container"
      transition:scale={{
        duration: 200,
        start: 0.95,
        opacity: 0
      }}
    >
      <!-- svelte-ignore a11y-click-events-have-key-events -->
      <div 
        class="modal-content {class_}"
        on:click|stopPropagation
      >
        <div class="modal-header">
          <h2>{title}</h2>
          <button 
            class="close-btn" 
            on:click={close}
            aria-label="关闭"
          >
            &times;
          </button>
        </div>
        <div class="modal-body">
          <slot />
        </div>
      </div>
    </div>
  </div>
{/if}

<style>
  .modal-wrapper {
    position: fixed;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 16px;
    z-index: 9999;
  }

  .modal-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.4);
    backdrop-filter: blur(4px);
    -webkit-backdrop-filter: blur(4px);
  }

  .modal-container {
    position: relative;
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1;
  }

  .modal-content {
    position: relative;
    width: 90%;
    max-width: 90vw;
    min-width: min-content;
    max-height: 90vh;
    background: white;
    border-radius: var(--radius-lg);
    box-shadow: 
      0 10px 25px -5px rgba(0, 0, 0, 0.1),
      0 8px 10px -6px rgba(0, 0, 0, 0.1);
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  .modal-header {
    padding: 16px 24px;
    border-bottom: 1px solid var(--border-color);
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: rgba(255, 255, 255, 0.9);
    backdrop-filter: blur(8px);
    -webkit-backdrop-filter: blur(8px);
    flex-shrink: 0;
  }

  h2 {
    font-size: var(--text-lg);
    font-weight: var(--font-semibold);
    color: var(--text-primary);
    margin: 0;
  }

  .close-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    padding: 0;
    background: none;
    border: none;
    border-radius: var(--radius-md);
    font-size: 24px;
    color: var(--text-secondary);
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .close-btn:hover {
    background: var(--hover-bg);
    color: var(--text-primary);
  }

  .modal-body {
    flex: 1;
    min-height: 0;
    overflow: auto;
    display: flex;
    flex-direction: column;
  }

  /* 确保模态窗口在所有内容之上 */
  :global(body:has(.modal-wrapper)) {
    overflow: hidden;
  }

  /* 移除 width 限制，让子组件自己控制宽度 */
  :global(.modal-content.result-modal) {
    max-width: none !important;
    width: auto !important;
    height: 90vh !important;
    display: flex;
    flex-direction: column;
  }
</style> 