import { defineStore } from 'pinia'
import { ref } from 'vue'

// 主题 store:持久化深/浅色偏好(加分项②主题切换)。
export const useThemeStore = defineStore('theme', () => {
  const dark = ref(localStorage.getItem('theme') === 'dark')

  function toggle() {
    dark.value = !dark.value
    localStorage.setItem('theme', dark.value ? 'dark' : 'light')
  }

  return { dark, toggle }
})
