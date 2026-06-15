<script setup lang="ts">
import { ref } from 'vue'
import {
  NCard,
  NInput,
  NInputNumber,
  NButton,
  NSpace,
  NText,
  NTag,
  NAlert,
} from 'naive-ui'
import StreamOutput from '../components/StreamOutput.vue'
import { useStreamTask } from '../composables/useStreamTask'

const text = ref('')
const maxPoints = ref(3)

const { output, status, elapsedMs, errorMsg, run, stop } = useStreamTask()

async function submit() {
  if (!text.value.trim()) return
  await run('summarize', { text: text.value, maxPoints: maxPoints.value })
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
  <n-card title="文本总结">
    <n-space vertical size="large">
      <n-space align="center">
        <n-text>要点数</n-text>
        <n-input-number v-model:value="maxPoints" :min="1" :max="10" style="width: 120px" />
      </n-space>
      <n-input
        v-model:value="text"
        type="textarea"
        placeholder="粘贴需要总结的长文本…"
        :autosize="{ minRows: 6, maxRows: 16 }"
      />
      <n-space>
        <n-button type="primary" :disabled="status === 'running' || !text.trim()" @click="submit">
          总结
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
