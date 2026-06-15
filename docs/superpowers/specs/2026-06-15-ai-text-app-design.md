# AI 文本处理应用 — 设计文档

> 状态:已批准 | 日期:2026-06-15 | 作者:Vibe Coding (Claude Code + superpowers)

## 1. 目标与范围

实现一个 AI 文本处理应用,提供三个功能:**中译英、英译中、文本总结**。
通过 DeepSeek API 真实调用大模型,SSE 流式逐 token 输出(打字机效果),支持取消/停止生成。

### 交付范围(含全部加分项)

- **基础**:前端三页 + SSE 流式 + 取消;后端三接口 + SSE;Go CLI 工具 + `skill.md`(Agent 可发现调用)。
- **加分①**:SDD/TDD 流程(本设计文档 + `spec/` + `agent.md`,经 superpowers 工作流产出)。
- **加分②**:全栈深度 — 前端虚拟滚动/流式渲染/深浅主题/响应式;后端参数校验/统一错误/日志追踪/超时控制。
- **加分③**:异步任务与轮询 — 内存队列 + worker + 状态机 + 前端轮询。
- **加分④**:数据闭环 — PostgreSQL 记录每次调用的输入/输出/耗时/状态 + 历史查询页。
- **加分⑤**:Docker — `Dockerfile` + `docker-compose.yml` 一键启动。

## 2. 技术栈

| 层 | 选型 |
|----|------|
| 前端 | Vue 3 + TypeScript + Vite + Naive UI + Pinia + Vue Router |
| 后端 | Go + Gin |
| 大模型 | DeepSeek(`deepseek-chat`,OpenAI 兼容,`stream:true`);无 Key 时 Mock fallback |
| 存储 | PostgreSQL(数据闭环历史,pgx/pgxpool 驱动) |
| CLI | Go + cobra |
| 部署 | Docker + docker-compose |

## 3. 架构:统一任务模型 + 双消费通道(方案 A)

核心:一个 **Task** 抽象,既能 SSE 实时消费,也能轮询查状态,还落库供历史查询 —— 把"SSE 流式"与"异步队列+轮询"两个需求统一。

```
POST /api/task
  → 创建 Task(pending)、入内存队列、返回 SSE 连接
  → Worker 取任务 → running → 调 DeepSeek 流式
  → 每个 token publish 到该任务的 pub/sub broker
  → SSE handler subscribe → 逐 token 转发前端(打字机)
  → 完成 → done/failed,结果+耗时写入 Postgres

GET    /api/task/{id}   → 轮询状态/进度(pending/running/done/failed/cancelled)
GET    /api/tasks       → 历史列表(查询页)
DELETE /api/task/{id}   → cancel context → 中断 LLM 调用 → cancelled
```

## 4. 后端设计

### 4.1 目录结构

```
backend/
  cmd/server/main.go
  internal/
    config/      # 读 .env:DEEPSEEK_API_KEY、端口、超时秒数
    model/       # Task、TaskStatus、TaskType、FunctionDef
    llm/         # Client 接口 + DeepSeekClient(流式) + MockClient
    task/        # Queue、Worker、Broker(pub/sub)、Manager(registry)
    store/       # PostgreSQL Repository(pgxpool)
    handler/     # functions / task(SSE) / cancel / history
    middleware/  # traceID 日志 / recover 统一错误 / 参数校验
```

### 4.2 数据模型

```go
type TaskType string   // "translate" | "summarize"
type TaskStatus string // pending | running | done | failed | cancelled

type Task struct {
    ID        string
    Type      TaskType
    Params    map[string]any   // translate:{text,from,to}; summarize:{text,maxPoints}
    Status    TaskStatus
    Result    string
    Err       string
    CreatedAt time.Time
    ElapsedMs int64
}
```

### 4.3 接口契约

| 方法 | 路径 | 请求 | 响应 |
|------|------|------|------|
| GET | `/api/functions` | — | `[{type,name,description,params}]` 共 3 项 |
| POST | `/api/task` | `{type, params}` | **SSE 流**:`meta`→`token*`→`done`/`error` |
| GET | `/api/task/{id}` | — | `{id,type,status,result,elapsedMs,...}` |
| GET | `/api/tasks?limit=50` | — | 历史列表 |
| DELETE | `/api/task/{id}` | — | `{status:"cancelled"}` |

**功能列表**(`/api/functions` 返回):
1. 中译英 — type=translate, from=zh, to=en
2. 英译中 — type=translate, from=en, to=zh
3. 文本总结 — type=summarize, 支持 maxPoints

**SSE 事件格式**:
```
event: meta\ndata: {"taskId":"..."}\n\n
event: token\ndata: {"text":"H"}\n\n
event: done\ndata: {"status":"done","elapsedMs":1234}\n\n
event: error\ndata: {"message":"...","traceId":"..."}\n\n
```

### 4.4 pub/sub Broker

每个 Task 一个 Broker:`Publish(token)` 写入缓冲 + 广播给订阅者;`Subscribe()` 返回 channel;缓冲区保存已产出 token,支持 SSE 断线重连补发。Worker 是 publisher,SSE handler 是 subscriber。

### 4.5 LLM 层

```go
type Client interface {
    Stream(ctx context.Context, messages []Message, out chan<- string) error
}
```
- `DeepSeekClient`:POST `https://api.deepseek.com/chat/completions`,`stream:true`,解析 SSE 增量 `delta.content` 推入 out。
- `MockClient`:无 `DEEPSEEK_API_KEY` 时启用,逐字吐预设结果模拟流式;**保留真实调用链路,注释标明 mock**。
- Prompt:翻译用"将以下文本从 X 翻译为 Y,只输出译文";总结用"用不超过 N 个要点总结"。

### 4.6 横切关注点(加分②)

- **超时**:`context.WithTimeout`(默认 60s),超时取消 LLM 调用,状态 failed。
- **统一错误**:中间件 `recover` + 错误响应 `{code,message,traceId}`。
- **日志追踪**:每请求注入 `traceId`,结构化日志贯穿 handler→worker→llm。
- **参数校验**:空文本、文本超长(如 >20000 字)、非法语向/类型 → 400。

## 5. 前端设计

### 5.1 目录结构

```
frontend/src/
  api/        # sse.ts(fetch+ReadableStream 消费 SSE)、client.ts
  views/      # ListView、TranslateView、SummarizeView、HistoryView
  components/ # StreamOutput(打字机+虚拟滚动)、ThemeToggle、TaskStatusBadge
  stores/     # task.ts(Pinia:当前任务状态/进度/轮询)
  router/     # 路由
  App.vue main.ts
```

### 5.2 关键点

- **SSE 消费**:用 `fetch` + `ReadableStream`(因需 POST body,不用 EventSource);手动解析 `event:`/`data:` 帧逐 token 追加。
- **取消/停止**:`AbortController.abort()` 断流 + 调 `DELETE /api/task/{id}` 通知后端终止。
- **StreamOutput**:打字机渲染;大文本**虚拟滚动**只渲染可视区(加分②)。
- **主题**:Naive UI `darkTheme` 深浅切换(加分②);响应式(窄屏单列)。
- **轮询模式**:Pinia store 暴露任务状态,可对 `GET /api/task/{id}` 轮询展示 pending/running 进度(加分③)。
- **历史页**:拉 `GET /api/tasks` 展示输入/输出/耗时/状态(加分④)。

## 6. CLI 设计(cobra)

```
ai-app translate --text "Hello" --from en --to zh
ai-app summarize --text "长文本..." --max-points 3
```
- 调后端 `POST /api/task`,消费 SSE,流式打印 stdout。
- flags:`--server`(后端地址,默认 `http://localhost:8080`)、`--json`(结构化输出)。
- `skill.md` 描述工具能力、参数、调用示例,使其可被 Claude Code/OpenClaw 等 Agent 发现并执行(附调用截图)。

## 7. 数据流(一次翻译)

```
前端输入 → POST /api/task → 入队 → worker 取 → running
  → DeepSeek 流式 → token 经 Broker → SSE → 前端打字机
  → 完成 → 写 Postgres → 历史页可查
取消:DELETE → cancel context → HTTP 中断 → cancelled
```

## 8. 测试策略(TDD 红绿)

- **后端** `go test` + `httptest`:覆盖 Broker pub/sub、任务生命周期状态机、SSE 事件解析、取消、Postgres 仓储、Mock LLM 流式。表驱动。
- **前端** Vitest + Vue Test Utils:测 StreamOutput 渲染、Pinia store(mock SSE 流)、SSE 解析器。
- **CLI**:对 httptest server 测命令解析与流式输出。

## 9. 交付物清单

- [x] 项目源码仓库(本仓库)
- [x] `skill.md`(Agent 可发现) + ⏳ Agent 调用截图(需人工补)
- [x] `agent.md`、`spec/`(加分①)
- [x] 前后端全栈深度(加分②)
- [x] 异步任务+轮询(加分③)
- [x] 数据闭环 + 历史页(加分④)
- [x] `Dockerfile` + `docker-compose.yml`(加分⑤)
- [x] `README.md`:项目介绍/演示、技术栈、本地运行、API 文档

## 10. 非目标(YAGNI)

- 不做用户认证/多租户。
- 不做 Redis/Celery(内存队列足够;设计预留接口便于替换)。Postgres 仅用于历史记录,不承担队列。
- 不做生产级分布式;聚焦单机一键启动与功能完整。
