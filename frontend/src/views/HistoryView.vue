<script setup lang="ts">
import { onMounted, ref, h } from 'vue'
import { NCard, NDataTable, NButton, NSpace, NTag, NText } from 'naive-ui'
import { fetchHistory, type TaskRecord } from '../api/client'

const rows = ref<TaskRecord[]>([])
const loading = ref(true)

const typeLabel: Record<string, string> = { translate: '翻译', summarize: '总结' }
const statusType: Record<string, 'default' | 'info' | 'success' | 'error' | 'warning'> = {
  pending: 'default',
  running: 'info',
  done: 'success',
  failed: 'error',
  cancelled: 'warning',
}

function truncate(s: string, n = 40) {
  return s && s.length > n ? s.slice(0, n) + '…' : s
}

const columns = [
  { title: '功能', key: 'type', width: 70, render: (r: TaskRecord) => typeLabel[r.type] ?? r.type },
  {
    title: '输入',
    key: 'input',
    render: (r: TaskRecord) => truncate(String((r.params as any)?.text ?? '')),
  },
  { title: '输出', key: 'result', render: (r: TaskRecord) => truncate(r.result) },
  {
    title: '状态',
    key: 'status',
    width: 90,
    render: (r: TaskRecord) => h(NTag, { type: statusType[r.status], size: 'small' }, () => r.status),
  },
  { title: '耗时(ms)', key: 'elapsedMs', width: 90 },
  {
    title: '时间',
    key: 'createdAt',
    width: 160,
    render: (r: TaskRecord) => new Date(r.createdAt).toLocaleString(),
  },
]

async function load() {
  loading.value = true
  try {
    rows.value = await fetchHistory(100)
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <n-card title="调用历史">
    <template #header-extra>
      <n-button size="small" @click="load">刷新</n-button>
    </template>
    <n-space vertical>
      <n-text depth="3">记录每次调用的输入、输出、耗时与状态(数据闭环)</n-text>
      <n-data-table
        :columns="columns"
        :data="rows"
        :loading="loading"
        :pagination="{ pageSize: 10 }"
        :bordered="false"
      />
    </n-space>
  </n-card>
</template>
