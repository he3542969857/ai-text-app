---
name: ai-text-tool
description: 通过 ai-app CLI 调用 AI 文本处理后端,完成中译英、英译中、文本总结。当用户需要翻译文本或总结长文本时使用。
---

# AI 文本处理工具(ai-app)

本技能封装命令行工具 `ai-app`,它调用后端服务(默认 `http://localhost:8080`)完成三类文本处理:**中译英、英译中、文本总结**。结果以流式(打字机)方式返回。

## 何时使用

- 用户要把中文翻译成英文,或英文翻译成中文。
- 用户要把一段长文本总结为若干要点。

## 前置条件

后端服务已启动并可访问。可用环境变量或 `--server` 指定地址:

```bash
# 启动后端(项目根 backend/ 下)
DATABASE_URL=postgres://postgres:postgres@localhost:5432/aitext?sslmode=disable \
  go run ./cmd/server
```

构建 CLI:

```bash
cd cli && go build -o ai-app .   # Windows: ai-app.exe
```

## 命令用法

### 翻译

```bash
ai-app translate --text "<待翻译文本>" --from <zh|en> --to <zh|en> [--server URL] [--json]
```

示例:
```bash
ai-app translate --text "Hello world" --from en --to zh
ai-app translate --text "你好世界" --from zh --to en
```

### 总结

```bash
ai-app summarize --text "<长文本>" --max-points <N> [--server URL] [--json]
```

示例:
```bash
ai-app summarize --text "一段很长的文本……" --max-points 3
```

### 全局参数

| 参数 | 说明 | 默认 |
|------|------|------|
| `--server` | 后端服务地址 | `http://localhost:8080` |
| `--json` | 以 JSON 输出最终结果(taskId/status/elapsedMs) | 关闭 |

## 输出

- 默认:逐 token 流式打印译文/总结到标准输出。
- `--json`:打印结构化结果,便于 Agent 解析:

```json
{
  "taskId": "847de960-...",
  "status": "done",
  "elapsedMs": 443,
  "error": ""
}
```

非零退出码表示失败(如参数校验未通过),错误信息打印到 stderr。

## 被 Agent 发现与调用

- **Claude Code / OpenClaw**:将本仓库的 `skill.md` 放入技能目录(或本项目根),Agent 读取 frontmatter 的 `name`/`description` 即可发现本工具,并按上文命令格式执行 `ai-app`。
- Agent 典型调用流程:
  1. 识别用户意图(翻译 / 总结)。
  2. 组装命令:`ai-app translate --text "..." --from en --to zh --json`。
  3. 通过 shell 执行,读取 stdout(流式文本)或 `--json` 结构化结果。

> 📸 验证:在 Agent(如 Claude Code)中请求"把 Hello 翻译成中文",Agent 应发现并执行 `ai-app translate --text "Hello" --from en --to zh`。请将该调用过程截图保存为 `docs/agent-invocation.png` 作为交付凭证。
