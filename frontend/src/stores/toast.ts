import toast from 'svelte-french-toast'
import type { ToastPosition } from 'svelte-french-toast'

// 定义通知类型
type ToastType = 'success' | 'error' | 'info' | 'warning'

// 通知选项接口
interface ToastOptions {
  duration?: number
  position?: ToastPosition
  icon?: string
}

// 默认选项
const defaultOptions: ToastOptions = {
  duration: 3000,
  position: 'top-right'
}

// 创建通知函数
function createToast() {
  return {
    show(message: string, type: ToastType = 'info', options: ToastOptions = {}) {
      const mergedOptions = { ...defaultOptions, ...options }
      
      switch (type) {
        case 'success':
          toast.success(message, mergedOptions)
          break
        case 'error':
          toast.error(message, mergedOptions)
          break
        case 'info':
          toast(message, mergedOptions)
          break
        case 'warning':
          toast(message, {
            ...mergedOptions,
            icon: '⚠️'
          })
          break
      }
    },
    success(message: string, options: ToastOptions = {}) {
      this.show(message, 'success', options)
    },
    error(message: string, options: ToastOptions = {}) {
      this.show(message, 'error', options)
    },
    info(message: string, options: ToastOptions = {}) {
      this.show(message, 'info', options)
    },
    warning(message: string, options: ToastOptions = {}) {
      this.show(message, 'warning', options)
    },
    dismiss() {
      toast.dismiss()
    },
    promise<T>(
      promise: Promise<T>,
      messages: {
        loading: string
        success: string
        error: string
      },
      options: ToastOptions = {}
    ) {
      return toast.promise(
        promise,
        {
          loading: messages.loading,
          success: messages.success,
          error: messages.error
        },
        options
      )
    }
  }
}

export const toastStore = createToast() 