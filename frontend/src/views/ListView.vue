<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { NCard, NGrid, NGi, NText, NSpin, NAlert } from 'naive-ui'
import { fetchFunctions, type FunctionDef } from '../api/client'

const router = useRouter()
const functions = ref<FunctionDef[]>([])
const loading = ref(true)
const error = ref('')

const routeFor = (type: string) => (type === 'summarize' ? '/summarize' : '/translate')

onMounted(async () => {
  try {
    functions.value = await fetchFunctions()
  } catch (e) {
    error.value = (e as Error).message
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <n-text strong style="font-size: 22px">功能列表</n-text>
    <p><n-text depth="3">选择一个功能开始使用</n-text></p>

    <n-spin :show="loading">
      <n-alert v-if="error" type="error" :title="error" style="margin-bottom: 16px" />
      <n-grid :cols="3" :x-gap="16" :y-gap="16" responsive="screen" item-responsive>
        <n-gi v-for="(fn, i) in functions" :key="i" span="3 m:1">
          <n-card
            hoverable
            style="cursor: pointer; height: 100%"
            @click="router.push(routeFor(fn.type))"
          >
            <template #header>{{ fn.name }}</template>
            <n-text depth="3">{{ fn.description }}</n-text>
          </n-card>
        </n-gi>
      </n-grid>
    </n-spin>
  </div>
</template>
