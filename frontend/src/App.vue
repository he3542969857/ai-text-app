<script setup lang="ts">
import { computed } from 'vue'
import { RouterView, useRoute, useRouter } from 'vue-router'
import { NConfigProvider, NGlobalStyle, darkTheme, zhCN, type GlobalThemeOverrides } from 'naive-ui'
import { useThemeStore } from './stores/theme'

const themeStore = useThemeStore()
const theme = computed(() => (themeStore.dark ? darkTheme : null))

const themeOverrides = computed<GlobalThemeOverrides>(() => {
  const accent = themeStore.dark ? '#0a84ff' : '#0071e3'
  return {
    common: {
      primaryColor: accent,
      primaryColorHover: themeStore.dark ? '#2a96ff' : '#0077ed',
      primaryColorPressed: themeStore.dark ? '#0a84ff' : '#006edb',
      borderRadius: '12px',
      borderRadiusSmall: '10px',
      fontFamily:
        '-apple-system, BlinkMacSystemFont, "SF Pro Display", "SF Pro Text", "Helvetica Neue", "PingFang SC", "Microsoft YaHei", system-ui, sans-serif',
      fontSize: '15px',
    },
  }
})

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
  <n-config-provider :theme="theme" :theme-overrides="themeOverrides" :locale="zhCN">
    <n-global-style />
    <header class="glass-nav">
      <div class="ap-container nav-inner">
        <div class="brand" @click="router.push('/')">
          <span class="brand-emoji">✨</span>
          <span class="brand-name">AI 文本工坊</span>
        </div>
        <nav class="nav-links">
          <button
            v-for="n in navs"
            :key="n.name"
            class="nav-link"
            :class="{ active: route.name === n.name }"
            @click="router.push({ name: n.name })"
          >
            {{ n.label }}
          </button>
        </nav>
        <button class="theme-btn" @click="themeStore.toggle()">
          {{ themeStore.dark ? '☀️' : '🌙' }}
        </button>
      </div>
    </header>

    <main class="ap-container page">
      <RouterView v-slot="{ Component }">
        <component :is="Component" :key="route.path" />
      </RouterView>
    </main>
  </n-config-provider>
</template>

<style scoped>
.nav-inner {
  height: 52px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}
.brand {
  display: flex;
  align-items: center;
  gap: 7px;
  cursor: pointer;
  font-weight: 600;
  letter-spacing: -0.02em;
  font-size: 17px;
  color: var(--text);
}
.brand-emoji {
  font-size: 18px;
}
.nav-links {
  display: flex;
  gap: 4px;
  background: color-mix(in srgb, var(--text) 6%, transparent);
  padding: 3px;
  border-radius: 980px;
}
.nav-link {
  appearance: none;
  border: none;
  background: transparent;
  cursor: pointer;
  font-family: inherit;
  font-size: 14px;
  font-weight: 500;
  letter-spacing: -0.01em;
  color: var(--text-secondary);
  padding: 6px 16px;
  border-radius: 980px;
  transition: color 0.2s ease, background 0.3s ease;
}
.nav-link:hover {
  color: var(--text);
}
.nav-link.active {
  background: var(--bg-elevated);
  color: var(--text);
  box-shadow: var(--shadow-sm);
}
.theme-btn {
  appearance: none;
  border: none;
  background: transparent;
  cursor: pointer;
  font-size: 18px;
  width: 38px;
  height: 38px;
  border-radius: 50%;
  transition: background 0.2s ease;
}
.theme-btn:hover {
  background: color-mix(in srgb, var(--text) 8%, transparent);
}
.page {
  padding-top: 48px;
  padding-bottom: 80px;
}

@media (max-width: 560px) {
  .brand-name {
    display: none;
  }
  .nav-link {
    padding: 6px 12px;
  }
}
</style>
