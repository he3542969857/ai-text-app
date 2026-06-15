# 接口设计

基地址:`http://localhost:8080`

## GET /api/functions

返回所有可用功能。

**响应 200**
```json
[
  { "type": "translate", "name": "中译英", "description": "将中文翻译为英文", "params": { "from": "zh", "to": "en" } },
  { "type": "translate", "name": "英译中", "description": "将英文翻译为中文", "params": { "from": "en", "to": "zh" } },
  { "type": "summarize", "name": "文本总结", "description": "将长文本总结为要点", "params": { "maxPoints": 3 } }
]
```

## POST /api/task

提交任务,**SSE 流式**返回。

**请求**
```json
{ "type": "translate", "params": { "text": "Hello", "from": "en", "to": "zh" } }
```
```json
{ "type": "summarize", "params": { "text": "长文本…", "maxPoints": 3 } }
```

**响应 200**(`Content-Type: text/event-stream`)
```
event: meta
data: {"taskId":"<uuid>"}

event: token
data: {"text":"你"}

event: token
data: {"text":"好"}

event: done
data: {"status":"done","elapsedMs":1234}
```

失败时在 `done` 前发 `error` 事件:
```
event: error
data: {"message":"<错误信息>"}
```

**校验失败 400**
```json
{ "code": "validation", "message": "text 不能为空" }
```

## GET /api/task/{id}

轮询单个任务状态。

**响应 200**
```json
{ "id":"<uuid>","type":"translate","params":{...},"status":"done","result":"你好","elapsedMs":1234,"createdAt":"2026-06-15T..." }
```
**404** 任务不存在。

状态枚举:`pending | running | done | failed | cancelled`

## GET /api/tasks?limit=50

历史记录列表(按创建时间倒序)。

**响应 200**:`TaskRecord[]`(结构同上)。

## DELETE /api/task/{id}

取消运行中的任务。

**响应 200**
```json
{ "status": "cancelled" }
```
**404** 任务不存在或已结束。

## 错误响应通用结构

```json
{ "code": "<错误码>", "message": "<说明>", "traceId": "<追踪ID>" }
```
服务器内部错误(panic)返回 500,带 `traceId`。
