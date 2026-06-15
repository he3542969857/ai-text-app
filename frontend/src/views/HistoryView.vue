<script setup lang="ts">
import { onMounted, ref, h } from 'vue'
import { NDataTable, NTag } from 'naive-ui'
import { fetchHistory, type TaskRecord } from '../api/client'

const rows = ref<TaskRecord[]>([])
const loading = ref(true)

const typeLabel: Record<string, string> = { translate: '🌐 翻译', summarize: '📝 总结' }
const statusType: Record<string, 'default' | 'info' | 'success' | 'error' | 'warning'> = {
  pending: 'default',
  running: 'info',
  done: 'success',
  failed: 'error',
  cancelled: 'warning',
}

function truncate(s: string, n = 36) {
  return s && s.length > n ? s.slice(0, n) + '…' : s
}

const columns = [
  { title: '功能', key: 'type', width: 90, render: (r: TaskRecord) => typeLabel[r.type] ?? r.type },
  {
    title: '输入',
    key: 'input',
    render: (r: TaskRecord) => truncate(String((r.params as any)?.text ?? '')),
  },
  { title: '输出', key: 'result', render: (r: TaskRecord) => truncate(r.result) },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render: (r: TaskRecord) =>
      h(NTag, { type: statusType[r.status], size: 'small', round: true }, () => r.status),
  },
  { title: '耗时', key: 'elapsedMs', width: 90, render: (r: TaskRecord) => `${r.elapsedMs} ms` },
  {
    title: '时间',
    key: 'createdAt',
    width: 170,
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
  <div class="wrap ap-fade-up">
    <header class="head">
      <div>
        <h1 class="ap-title">🕘 调用历史</h1>
        <p class="ap-subtitle">每次调用的输入、输出、耗时与状态</p>
      </div>
      <button class="ap-btn ap-btn-ghost" @click="load">🔄 刷新</button>
    </header>

    <section class="ap-card table-card">
      <n-data-table
        :columns="columns"
        :data="rows"
        :loading="loading"
        :pagination="{ pageSize: 10 }"
        :bordered="false"
      />
    </section>
  </div>
</template>

<style scoped>
.wrap {
  display: flex;
  flex-direction: column;
  gap: 22px;
}
.head {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 16px;
}
.head h1 {
  font-size: 32px;
}
.head p {
  margin-top: 8px;
}
.table-card {
  padding: 14px 16px;
}
</style>
