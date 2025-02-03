import { writable } from 'svelte/store';

// 获取系统主题偏好
const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;

// 从本地存储获取主题设置，如果没有则使用系统偏好
const storedTheme = localStorage.getItem('theme');
const initialTheme = storedTheme || (prefersDark ? 'dark' : 'light');

// 创建主题 store
export const themeStore = writable(initialTheme);

// 订阅主题变化并保存到本地存储
themeStore.subscribe(value => {
  if (typeof window !== 'undefined') {
    localStorage.setItem('theme', value);
    
    // 添加过渡类
    document.documentElement.classList.add('no-transition');
    
    // 切换主题类
    if (value === 'dark') {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
    
    // 强制重排以应用新样式
    void document.documentElement.offsetHeight;
    
    // 移除过渡类，允许后续变化有动画
    requestAnimationFrame(() => {
      document.documentElement.classList.remove('no-transition');
    });
  }
});

// 切换主题
export function toggleTheme() {
  themeStore.update(current => current === 'light' ? 'dark' : 'light');
} 