// SSE 解析与任务流式调用。
// 拆成「纯解析器」+「fetch 流式消费」两层,前者可独立单元测试。

export interface SSEEvent {
  event: string
  data: string
}

/**
 * createSSEParser 返回一个增量解析器:可多次 feed 任意切分的 chunk,
 * 每解析出一个完整 SSE 事件(以空行 \n\n 分隔)就回调 onEvent。
 */
export function createSSEParser(onEvent: (e: SSEEvent) => void) {
  let buf = ''
  return {
    feed(chunk: string) {
      buf += chunk
      let idx: number
      while ((idx = buf.indexOf('\n\n')) >= 0) {
        const raw = buf.slice(0, idx)
        buf = buf.slice(idx + 2)
        const e = parseEvent(raw)
        if (e) onEvent(e)
      }
    },
  }
}

function parseEvent(raw: string): SSEEvent | null {
  let event = 'message'
  const dataLines: string[] = []
  for (const line of raw.split('\n')) {
    if (line.startsWith('event:')) {
      event = line.slice(6).trim()
    } else if (line.startsWith('data:')) {
      dataLines.push(line.slice(5).trim())
    }
  }
  if (dataLines.length === 0) return null
  return { event, data: dataLines.join('\n') }
}

export interface TaskHandlers {
  onMeta?: (taskId: string) => void
  onToken?: (text: string) => void
  onDone?: (status: string, elapsedMs: number) => void
  onError?: (message: string) => void
}

/** dispatchEvent 将一个 SSE 事件分发到对应回调。导出以便测试。 */
export function dispatchEvent(e: SSEEvent, h: TaskHandlers) {
  let payload: any = {}
  try {
    payload = JSON.parse(e.data)
  } catch {
    payload = {}
  }
  switch (e.event) {
    case 'meta':
      h.onMeta?.(payload.taskId)
      break
    case 'token':
      h.onToken?.(payload.text ?? '')
      break
    case 'error':
      h.onError?.(payload.message ?? '未知错误')
      break
    case 'done':
      h.onDone?.(payload.status ?? 'done', payload.elapsedMs ?? 0)
      break
  }
}

/**
 * streamTask 提交任务并消费 SSE 流。使用 fetch + ReadableStream(因需 POST body)。
 * 传入 signal 可通过 AbortController 中断(停止生成)。
 */
export async function streamTask(
  type: string,
  params: Record<string, unknown>,
  handlers: TaskHandlers,
  signal?: AbortSignal,
): Promise<void> {
  const resp = await fetch('/api/task', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ type, params }),
    signal,
  })

  if (!resp.ok) {
    let msg = `请求失败(${resp.status})`
    try {
      const j = await resp.json()
      msg = j.message || msg
    } catch {
      /* ignore */
    }
    handlers.onError?.(msg)
    return
  }

  const reader = resp.body?.getReader()
  if (!reader) {
    handlers.onError?.('响应不支持流式读取')
    return
  }
  const decoder = new TextDecoder()
  const parser = createSSEParser((e) => dispatchEvent(e, handlers))

  try {
    while (true) {
      const { done, value } = await reader.read()
      if (done) break
      parser.feed(decoder.decode(value, { stream: true }))
    }
  } catch (err) {
    if ((err as Error).name !== 'AbortError') {
      handlers.onError?.((err as Error).message)
    }
  }
}
