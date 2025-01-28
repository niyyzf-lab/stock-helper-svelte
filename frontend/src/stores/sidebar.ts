import { writable } from 'svelte/store'

interface SidebarStore {
  isCollapsed: boolean
}

const STORAGE_KEY = 'sidebar-state'

// 从 localStorage 读取初始状态
function getInitialState(): SidebarStore {
  try {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored) {
      return JSON.parse(stored)
    }
  } catch (err) {
    console.error('Error reading sidebar state from localStorage:', err)
  }
  return { isCollapsed: false }
}

function createSidebarStore() {
  const { subscribe, set, update } = writable<SidebarStore>(getInitialState())

  // 保存状态到 localStorage
  const saveState = (state: SidebarStore) => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(state))
    } catch (err) {
      console.error('Error saving sidebar state to localStorage:', err)
    }
  }

  return {
    subscribe,
    toggle: () => update(state => {
      const newState = { ...state, isCollapsed: !state.isCollapsed }
      saveState(newState)
      return newState
    }),
    setCollapsed: (value: boolean) => update(state => {
      const newState = { ...state, isCollapsed: value }
      saveState(newState)
      return newState
    }),
    reset: () => {
      const newState = { isCollapsed: false }
      saveState(newState)
      set(newState)
    }
  }
}

export const sidebarStore = createSidebarStore() 