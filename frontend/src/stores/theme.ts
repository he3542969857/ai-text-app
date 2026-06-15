import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

// 主题 store:持久化深/浅色偏好,并同步到 <html data-theme> 供全局 CSS 变量切换。
export const useThemeStore = defineStore('theme', () => {
  const dark = ref(localStorage.getItem('theme') === 'dark')

  function apply() {
    document.documentElement.dataset.theme = dark.value ? 'dark' : 'light'
  }
  apply()
  watch(dark, apply)

  function toggle() {
    dark.value = !dark.value
    localStorage.setItem('theme', dark.value ? 'dark' : 'light')
  }

  return { dark, toggle }
})
