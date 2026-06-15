// 后端 REST 接口封装(非流式部分)。

export interface FunctionDef {
  type: string
  name: string
  description: string
  params: Record<string, unknown>
}

export interface TaskRecord {
  id: string
  type: string
  params: Record<string, unknown>
  status: string
  result: string
  error?: string
  createdAt: string
  elapsedMs: number
}

export async function fetchFunctions(): Promise<FunctionDef[]> {
  const r = await fetch('/api/functions')
  if (!r.ok) throw new Error('获取功能列表失败')
  return r.json()
}

export async function fetchHistory(limit = 50): Promise<TaskRecord[]> {
  const r = await fetch(`/api/tasks?limit=${limit}`)
  if (!r.ok) throw new Error('获取历史失败')
  return r.json()
}

export async function fetchTask(id: string): Promise<TaskRecord> {
  const r = await fetch(`/api/task/${id}`)
  if (!r.ok) throw new Error('任务不存在')
  return r.json()
}

export async function cancelTask(id: string): Promise<void> {
  await fetch(`/api/task/${id}`, { method: 'DELETE' })
}
