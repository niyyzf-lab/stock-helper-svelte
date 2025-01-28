<script lang="ts">
  import { Minimize2, X, Minus } from 'lucide-svelte'

  // 这些函数需要通过 Wails 绑定实现
  function minimizeWindow() {
    (window as any).runtime.WindowMinimise()
  }
  
  function maximizeWindow() {
    (window as any).runtime.WindowToggleMaximise()
  }
  
  function closeWindow() {
    (window as any).runtime.Quit()
  }
</script>

<div class="window-controls" style="--wails-draggable:no-drag">
  <button on:click={minimizeWindow}>
    <Minus size={16} />
  </button>
  <button on:click={maximizeWindow}>
    <Minimize2 size={16} />
  </button>
  <button class="close" on:click={closeWindow}>
    <X size={16} />
  </button>
</div>

<style>
  .window-controls {
    display: flex;
    gap: 8px;
    margin-left: auto;
    padding-right: 8px;
  }

  button {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    border: none;
    border-radius: 6px;
    background: transparent;
    color: var(--neutral-400);
    cursor: pointer;
    transition: all 0.2s ease;
  }

  button:hover {
    background: var(--neutral-100);
    color: var(--neutral-700);
  }

  .close:hover {
    background: #ef4444;
    color: white;
  }
</style> 