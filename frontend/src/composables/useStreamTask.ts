import { ref } from 'vue'
import { streamTask } from '../api/sse'
import { cancelTask } from '../api/client'

// useStreamTask 封装一次流式任务的状态机,供翻译/总结页复用。
export function useStreamTask() {
  const output = ref('')
  const status = ref<'idle' | 'running' | 'done' | 'failed' | 'cancelled'>('idle')
  const elapsedMs = ref(0)
  const errorMsg = ref('')
  const taskId = ref('')

  let controller: AbortController | null = null

  async function run(type: string, params: Record<string, unknown>) {
    output.value = ''
    errorMsg.value = ''
    elapsedMs.value = 0
    status.value = 'running'
    controller = new AbortController()

    await streamTask(
      type,
      params,
      {
        onMeta: (id) => (taskId.value = id),
        onToken: (t) => (output.value += t),
        onError: (m) => {
          errorMsg.value = m
          status.value = 'failed'
        },
        onDone: (s, ms) => {
          elapsedMs.value = ms
          if (status.value === 'running') {
            status.value = s as typeof status.value
          }
        },
      },
      controller.signal,
    )
    controller = null
  }

  // stop 中断流式连接并通知后端取消(停止生成)。
  async function stop() {
    if (controller) {
      controller.abort()
      controller = null
    }
    if (taskId.value) {
      await cancelTask(taskId.value)
    }
    status.value = 'cancelled'
  }

  return { output, status, elapsedMs, errorMsg, taskId, run, stop }
}
