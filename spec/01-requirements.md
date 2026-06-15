# 需求拆分

## 1. 核心功能(必做)

| 编号 | 功能 | 说明 |
|------|------|------|
| F1 | 中译英 | 中文 → 英文翻译 |
| F2 | 英译中 | 英文 → 中文翻译 |
| F3 | 文本总结 | 长文本 → 要点(可控要点数) |

## 2. 前端需求

| 编号 | 需求 | 实现 |
|------|------|------|
| FE1 | 列表页:展示所有功能 | `ListView` 拉 `/api/functions` 渲染卡片 |
| FE2 | 翻译页:输入/语向/结果 | `TranslateView` + 语向下拉 |
| FE3 | 总结页:长文本/要点数/结果 | `SummarizeView` + 要点数控件 |
| FE4 | SSE 流式打字机渲染 | `useStreamTask` + `StreamOutput` |
| FE5 | 取消/停止生成 | `AbortController.abort()` + `DELETE /api/task/{id}` |

## 3. 后端需求

| 编号 | 需求 | 实现 |
|------|------|------|
| BE1 | `GET /api/functions` | 返回功能列表 |
| BE2 | `POST /api/task`(SSE) | 提交任务,流式返回 taskId + 结果 |
| BE3 | `DELETE /api/task/{id}` | 取消任务(context 取消) |
| BE4 | 逐 token 推送 | per-task pub/sub Broker |

## 4. CLI 需求

| 编号 | 需求 | 实现 |
|------|------|------|
| C1 | `ai-app translate` | 调后端流式翻译 |
| C2 | `ai-app summarize` | 调后端流式总结 |
| C3 | `skill.md` | Agent 可发现调用 |

## 5. 加分项

| 编号 | 需求 | 实现 |
|------|------|------|
| B1 | SDD/TDD + agent.md + spec/ | superpowers 工作流 |
| B2 | 全栈深度 | 虚拟滚动/流式优化/主题/响应式;校验/统一错误/日志/超时 |
| B3 | 异步任务+轮询 | 内存队列+worker+状态机;`GET /api/task/{id}` 轮询 |
| B4 | 数据闭环 | PostgreSQL 记录;历史查询页 |
| B5 | Docker | Dockerfile + docker-compose 一键启动 |

## 6. 非目标(YAGNI)

- 不做用户认证/多租户。
- 不做 Redis/Celery(内存队列足够,设计预留接口)。
- 不做分布式;聚焦单机功能完整与一键启动。
