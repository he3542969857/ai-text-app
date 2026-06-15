# AI 文本处理应用

一个全栈 AI 文本处理应用,提供 **中译英、英译中、文本总结** 三个功能。通过 DeepSeek API 真实调用大模型,SSE 流式逐 token 输出(打字机效果),支持取消/停止生成。

> 本项目以 **Vibe Coding + SDD/TDD** 方式开发(Claude Code + [superpowers](https://github.com/obra/superpowers) 工作流):brainstorming → 设计规范 → TDD 实现计划 → 红绿实现 → 收尾。详见 [`agent.md`](./agent.md)、[`docs/superpowers/`](./docs/superpowers/) 与 [`spec/`](./spec/)。

## ✨ 功能与特性

- **三大功能**:中译英 / 英译中 / 文本总结(可控要点数)
- **SSE 流式**:逐 token 打字机渲染,支持**停止生成**(前端断流 + 后端 `context` 取消中断 LLM)
- **统一任务模型**:一个 Task 抽象同时支持「SSE 实时消费」「轮询查状态」「历史落库」
- **异步队列**:内存队列 + worker 池 + 状态机(pending/running/done/failed/cancelled)
- **数据闭环**:PostgreSQL 记录每次调用的输入/输出/耗时/状态,提供历史查询页
- **全栈深度**:前端虚拟滚动 / 流式渲染 / 深浅主题 / 响应式;后端参数校验 / 统一错误 / traceID 日志 / 超时控制
- **CLI + Agent**:`ai-app` 命令行工具,可被 Claude Code / OpenClaw 等 Agent 发现并调用(见 [`skill.md`](./skill.md))
- **无 Key 也能跑**:未配置 `DEEPSEEK_API_KEY` 时自动降级为 Mock(逐字模拟,**保留真实调用链路**)

## 🧱 技术栈

| 层 | 选型 |
|----|------|
| 前端 | Vue 3 · TypeScript · Vite · Naive UI · Pinia · Vue Router |
| 后端 | Go · Gin |
| 大模型 | DeepSeek(`deepseek-chat`,OpenAI 兼容流式) |
| 存储 | PostgreSQL(pgx / pgxpool) |
| CLI | Go · cobra |
| 部署 | Docker · docker-compose |
| 测试 | Go `testing` + httptest;Vitest + Vue Test Utils |

## 📂 项目结构

```
ai-text-app/
├── backend/         # Go + Gin 后端
│   ├── cmd/server/  # 入口
│   └── internal/    # config/model/llm/task/store/middleware/handler
├── cli/             # Go + cobra CLI(ai-app)
├── frontend/        # Vue3 + TS 前端
│   └── src/         # api/views/components/stores/composables/router
├── spec/            # 需求拆分 / 接口设计 / 页面原型(SDD)
├── docs/superpowers/# 设计文档 + 实现计划
├── agent.md         # AI Agent 角色与协作
├── skill.md         # CLI 工具的 Agent 可发现说明
├── docker-compose.yml
└── README.md
```

## 🚀 本地运行

### 方式一:Docker 一键启动(推荐)

```bash
# 可选:配置真实 Key(不配则用 Mock)
export DEEPSEEK_API_KEY=sk-xxx

docker compose up --build
```
- 前端:http://localhost:3000
- 后端:http://localhost:8080
- Postgres:localhost:5432(compose 内置)

### 方式二:本地分别启动

**前置**:Go 1.26+、Node 20+、PostgreSQL 17,并建库:
```bash
createdb -U postgres aitext
```

**后端**:
```bash
cd backend
cp .env.example .env   # 按需填入 DEEPSEEK_API_KEY
# 设置数据库连接(或写入 .env / 环境变量)
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/aitext?sslmode=disable"
go run ./cmd/server     # 监听 :8080
```

**前端**:
```bash
cd frontend
npm install
npm run dev              # http://localhost:5173(已配置 /api 代理到 :8080)
```

**CLI**:
```bash
cd cli
go build -o ai-app .     # Windows: ai-app.exe
./ai-app translate --text "Hello world" --from en --to zh
./ai-app summarize --text "一段长文本…" --max-points 3
./ai-app translate --text "你好" --from zh --to en --json
```

## 🧪 测试

```bash
# 后端(需可连 Postgres;store/handler 为集成测试,无库则自动跳过)
cd backend
TEST_DATABASE_URL="postgres://postgres:postgres@localhost:5432/aitext?sslmode=disable" go test ./...

# CLI
cd cli && go test ./...

# 前端
cd frontend && npm run test
```

## 📡 API 文档

完整接口见 [`spec/02-api.md`](./spec/02-api.md)。摘要:

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/functions` | 功能列表 |
| POST | `/api/task` | 提交任务,**SSE 流式**返回 `meta`→`token*`→`done`/`error` |
| GET | `/api/task/{id}` | 轮询任务状态 |
| GET | `/api/tasks?limit=50` | 历史记录 |
| DELETE | `/api/task/{id}` | 取消任务 |

**SSE 示例**:
```
event: meta
data: {"taskId":"<uuid>"}

event: token
data: {"text":"你"}

event: done
data: {"status":"done","elapsedMs":1234}
```

## 🤖 Agent 集成

`ai-app` CLI 可被 Agent 发现并调用,详见 [`skill.md`](./skill.md)。典型流程:Agent 识别意图 → 组装 `ai-app translate --text "..." --from en --to zh --json` → 执行并读取结果。

下图为 Claude Code 读取 `skill.md` 后真实执行 `ai-app` 的调用凭证:

![Agent 调用截图](docs/agent-invocation.png)

## 📐 架构说明

**统一任务模型 + 双消费通道**:`POST /api/task` 创建任务并入内存队列,同时返回 SSE 连接;worker 调 DeepSeek 流式,token 经 per-task pub/sub Broker 广播给 SSE 订阅者;任务完成写入 PostgreSQL。这样同一个任务既能实时流式消费,又能轮询查状态、落库查历史 —— 把"SSE 流式"与"异步队列+轮询"两个需求统一在一套抽象里。
