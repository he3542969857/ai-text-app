<script setup lang="ts">
import { NButton } from 'naive-ui'
import { useThemeStore } from '../stores/theme'

const theme = useThemeStore()
</script>

<template>
  <n-button quaternary circle @click="theme.toggle()" :title="theme.dark ? '切换浅色' : '切换深色'">
    {{ theme.dark ? '🌙' : '☀️' }}
  </n-button>
</template>
