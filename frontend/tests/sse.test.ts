import { describe, it, expect, vi } from 'vitest'
import { createSSEParser, dispatchEvent, type SSEEvent } from '../src/api/sse'

describe('createSSEParser', () => {
  it('解析完整的多个事件', () => {
    const events: SSEEvent[] = []
    const p = createSSEParser((e) => events.push(e))
    p.feed('event: meta\ndata: {"taskId":"t1"}\n\n')
    p.feed('event: token\ndata: {"text":"你"}\n\n')
    expect(events).toHaveLength(2)
    expect(events[0]).toEqual({ event: 'meta', data: '{"taskId":"t1"}' })
    expect(events[1].event).toBe('token')
  })

  it('跨 chunk 切分的事件能被正确拼接', () => {
    const events: SSEEvent[] = []
    const p = createSSEParser((e) => events.push(e))
    p.feed('event: tok')
    p.feed('en\ndata: {"text":"好"}')
    expect(events).toHaveLength(0) // 尚未遇到空行
    p.feed('\n\n')
    expect(events).toHaveLength(1)
    expect(events[0].event).toBe('token')
    expect(events[0].data).toBe('{"text":"好"}')
  })
})

describe('dispatchEvent', () => {
  it('按事件类型分发到对应回调', () => {
    const onMeta = vi.fn()
    const onToken = vi.fn()
    const onDone = vi.fn()
    const onError = vi.fn()
    const handlers = { onMeta, onToken, onDone, onError }

    dispatchEvent({ event: 'meta', data: '{"taskId":"abc"}' }, handlers)
    dispatchEvent({ event: 'token', data: '{"text":"Hi"}' }, handlers)
    dispatchEvent({ event: 'done', data: '{"status":"done","elapsedMs":42}' }, handlers)
    dispatchEvent({ event: 'error', data: '{"message":"boom"}' }, handlers)

    expect(onMeta).toHaveBeenCalledWith('abc')
    expect(onToken).toHaveBeenCalledWith('Hi')
    expect(onDone).toHaveBeenCalledWith('done', 42)
    expect(onError).toHaveBeenCalledWith('boom')
  })

  it('data 非法 JSON 时不抛错', () => {
    const onToken = vi.fn()
    expect(() => dispatchEvent({ event: 'token', data: 'not-json' }, { onToken })).not.toThrow()
  })
})
