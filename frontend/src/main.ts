import './styles.css'
import App from './App.svelte'

const app = new App({
  target: document.getElementById('app')!,
})

// 确保初始路由正确
if (!location.hash) {
  location.hash = '#/'
}

export default app
