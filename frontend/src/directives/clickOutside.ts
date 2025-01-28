/// <reference types="svelte" />

type ClickOutsideEvent = CustomEvent<void>;

declare namespace svelteHTML {
  interface HTMLAttributes<T> {
    'on:clickoutside'?: (event: ClickOutsideEvent) => void;
  }
}

export function clickOutside(node: HTMLElement) {
  const handleClick = (event: MouseEvent) => {
    if (node && !node.contains(event.target as Node) && !event.defaultPrevented) {
      node.dispatchEvent(new CustomEvent('clickoutside'))
    }
  }

  document.addEventListener('click', handleClick, true)

  return {
    destroy() {
      document.removeEventListener('click', handleClick, true)
    }
  }
} 