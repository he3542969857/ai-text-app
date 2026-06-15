import { ref } from 'vue'
import { streamTask } from '../api/sse'
import { cancelTask, fetchTask } from '../api/client'

// useStreamTask 封装一次流式任务的状态机,供翻译/总结页复用。
// 同时演示「前端轮询任务状态」:拿到 taskId 后定时轮询 GET /api/task/{id}
// 展示后端侧进度反馈(与 SSE 实时流并存)。
export function useStreamTask() {
  const output = ref('')
  const status = ref<'idle' | 'running' | 'done' | 'failed' | 'cancelled'>('idle')
  const polledStatus = ref('') // 轮询得到的后端任务状态
  const elapsedMs = ref(0)
  const errorMsg = ref('')
  const taskId = ref('')

  let controller: AbortController | null = null
  let pollTimer: number | null = null

  function startPolling(id: string) {
    stopPolling()
    pollTimer = window.setInterval(async () => {
      try {
        const tk = await fetchTask(id)
        polledStatus.value = tk.status
        if (['done', 'failed', 'cancelled'].includes(tk.status)) stopPolling()
      } catch {
        /* 任务可能尚未注册,忽略 */
      }
    }, 800)
  }
  function stopPolling() {
    if (pollTimer !== null) {
      clearInterval(pollTimer)
      pollTimer = null
    }
  }

  async function run(type: string, params: Record<string, unknown>) {
    output.value = ''
    errorMsg.value = ''
    elapsedMs.value = 0
    polledStatus.value = ''
    status.value = 'running'
    controller = new AbortController()

    await streamTask(
      type,
      params,
      {
        onMeta: (id) => {
          taskId.value = id
          startPolling(id)
        },
        onToken: (t) => (output.value += t),
        onError: (m) => {
          errorMsg.value = m
          status.value = 'failed'
        },
        onDone: (s, ms) => {
          elapsedMs.value = ms
          if (status.value === 'running') status.value = s as typeof status.value
        },
      },
      controller.signal,
    )
    controller = null
    stopPolling()
  }

  // stop 中断流式连接并通知后端取消(停止生成)。
  async function stop() {
    if (controller) {
      controller.abort()
      controller = null
    }
    stopPolling()
    if (taskId.value) await cancelTask(taskId.value)
    status.value = 'cancelled'
  }

  return { output, status, polledStatus, elapsedMs, errorMsg, taskId, run, stop }
}
