<script setup lang="ts">
import { ref, computed } from 'vue'
import { NInput, NInputNumber } from 'naive-ui'
import StreamOutput from '../components/StreamOutput.vue'
import { useStreamTask } from '../composables/useStreamTask'

const text = ref('')
const maxPoints = ref(3)

const { output, status, elapsedMs, errorMsg, run, stop } = useStreamTask()

async function submit() {
  if (!text.value.trim()) return
  await run('summarize', { text: text.value, maxPoints: maxPoints.value })
}

const statusMeta = computed(() => {
  const map: Record<string, { emoji: string; label: string; cls: string }> = {
    idle: { emoji: '💤', label: '待命', cls: 'idle' },
    running: { emoji: '⚡️', label: '生成中', cls: 'running' },
    done: { emoji: '✅', label: '完成', cls: 'done' },
    failed: { emoji: '❌', label: '失败', cls: 'failed' },
    cancelled: { emoji: '🛑', label: '已停止', cls: 'cancelled' },
  }
  return map[status.value] ?? map.idle
})
</script>

<template>
  <div class="wrap ap-fade-up">
    <header class="head">
      <h1 class="ap-title">📝 文本总结</h1>
      <p class="ap-subtitle">长文本一键提炼为要点</p>
    </header>

    <section class="ap-card panel">
      <div class="points">
        <span class="points-label">🎯 要点数</span>
        <n-input-number v-model:value="maxPoints" :min="1" :max="10" class="points-input" />
      </div>
      <n-input
        v-model:value="text"
        type="textarea"
        placeholder="粘贴需要总结的长文本…"
        :autosize="{ minRows: 6, maxRows: 18 }"
      />
      <div class="actions">
        <button class="ap-btn" :disabled="status === 'running' || !text.trim()" @click="submit">
          总结 ✨
        </button>
        <button v-if="status === 'running'" class="ap-btn ap-btn-danger" @click="stop">
          停止生成
        </button>
        <span class="status" :class="statusMeta.cls">{{ statusMeta.emoji }} {{ statusMeta.label }}</span>
        <span v-if="elapsedMs > 0" class="elapsed">⏱ {{ elapsedMs }} ms</span>
      </div>
      <p v-if="errorMsg" class="err">⚠️ {{ errorMsg }}</p>
    </section>

    <section class="ap-card result">
      <div class="result-label">🗒 要点</div>
      <StreamOutput :text="output" />
    </section>
  </div>
</template>

<style scoped>
.wrap {
  display: flex;
  flex-direction: column;
  gap: 22px;
}
.head h1 {
  font-size: 32px;
}
.head p {
  margin-top: 8px;
}
.panel {
  padding: 22px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.points {
  display: flex;
  align-items: center;
  gap: 12px;
}
.points-label {
  font-size: 15px;
  color: var(--text-secondary);
}
.points-input {
  width: 130px;
}
.actions {
  display: flex;
  align-items: center;
  gap: 14px;
  flex-wrap: wrap;
}
.status {
  font-size: 14px;
  font-weight: 500;
  padding: 5px 12px;
  border-radius: 980px;
  background: color-mix(in srgb, var(--text) 6%, transparent);
  color: var(--text-secondary);
}
.status.running {
  color: var(--accent);
  background: color-mix(in srgb, var(--accent) 12%, transparent);
}
.status.done {
  color: #34c759;
  background: color-mix(in srgb, #34c759 14%, transparent);
}
.status.failed {
  color: #ff453a;
  background: color-mix(in srgb, #ff453a 14%, transparent);
}
.elapsed {
  font-size: 13px;
  color: var(--text-tertiary);
}
.err {
  margin: 0;
  color: #ff453a;
  font-size: 14px;
}
.result {
  padding: 18px 20px 20px;
}
.result-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-secondary);
  margin-bottom: 12px;
}
</style>
