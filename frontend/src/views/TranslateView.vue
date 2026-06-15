<script setup lang="ts">
import { ref } from 'vue'
import {
  NCard,
  NInput,
  NSelect,
  NButton,
  NSpace,
  NText,
  NTag,
  NAlert,
} from 'naive-ui'
import StreamOutput from '../components/StreamOutput.vue'
import { useStreamTask } from '../composables/useStreamTask'

const text = ref('')
const direction = ref<'en2zh' | 'zh2en'>('en2zh')
const directionOptions = [
  { label: '英译中', value: 'en2zh' },
  { label: '中译英', value: 'zh2en' },
]

const { output, status, elapsedMs, errorMsg, run, stop } = useStreamTask()

async function submit() {
  if (!text.value.trim()) return
  const [from, to] = direction.value === 'en2zh' ? ['en', 'zh'] : ['zh', 'en']
  await run('translate', { text: text.value, from, to })
}

const statusType: Record<string, 'default' | 'info' | 'success' | 'error' | 'warning'> = {
  idle: 'default',
  running: 'info',
  done: 'success',
  failed: 'error',
  cancelled: 'warning',
}
</script>

<template>
  <n-card title="文本翻译">
    <n-space vertical size="large">
      <n-select v-model:value="direction" :options="directionOptions" style="width: 160px" />
      <n-input
        v-model:value="text"
        type="textarea"
        placeholder="输入要翻译的文本…"
        :autosize="{ minRows: 4, maxRows: 10 }"
      />
      <n-space>
        <n-button type="primary" :disabled="status === 'running' || !text.trim()" @click="submit">
          翻译
        </n-button>
        <n-button v-if="status === 'running'" type="error" @click="stop">停止生成</n-button>
        <n-tag :type="statusType[status]">{{ status }}</n-tag>
        <n-text v-if="elapsedMs > 0" depth="3">耗时 {{ elapsedMs }} ms</n-text>
      </n-space>
      <n-alert v-if="errorMsg" type="error" :title="errorMsg" />
      <div>
        <n-text depth="3">结果</n-text>
        <StreamOutput :text="output" />
      </div>
    </n-space>
  </n-card>
</template>
