<script setup lang="ts">
import { computed } from 'vue'
import { RouterView, useRoute, useRouter } from 'vue-router'
import {
  NConfigProvider,
  NGlobalStyle,
  NLayout,
  NLayoutHeader,
  NLayoutContent,
  NSpace,
  NButton,
  NText,
  darkTheme,
  zhCN,
} from 'naive-ui'
import { useThemeStore } from './stores/theme'
import ThemeToggle from './components/ThemeToggle.vue'

const themeStore = useThemeStore()
const theme = computed(() => (themeStore.dark ? darkTheme : null))

const route = useRoute()
const router = useRouter()

const navs = [
  { name: 'list', label: '首页' },
  { name: 'translate', label: '翻译' },
  { name: 'summarize', label: '总结' },
  { name: 'history', label: '历史' },
]
</script>

<template>
  <n-config-provider :theme="theme" :locale="zhCN">
    <n-global-style />
    <n-layout style="min-height: 100vh">
      <n-layout-header bordered style="padding: 12px 20px">
        <n-space justify="space-between" align="center">
          <n-space align="center">
            <n-text strong style="font-size: 18px; cursor: pointer" @click="router.push('/')">
              🤖 AI 文本处理
            </n-text>
            <n-button
              v-for="n in navs"
              :key="n.name"
              text
              :type="route.name === n.name ? 'primary' : 'default'"
              @click="router.push({ name: n.name })"
            >
              {{ n.label }}
            </n-button>
          </n-space>
          <ThemeToggle />
        </n-space>
      </n-layout-header>
      <n-layout-content style="padding: 24px; max-width: 960px; margin: 0 auto; width: 100%">
        <RouterView />
      </n-layout-content>
    </n-layout>
  </n-config-provider>
</template>
