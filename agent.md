# agent.md — AI Agent 在本项目中的角色与协作方式

本项目采用 **Vibe Coding + 规范驱动开发(SDD)+ 测试驱动开发(TDD)**,由 AI 编程 Agent(Claude Code,配合 [superpowers](https://github.com/obra/superpowers) 技能库)与人类开发者协作完成。本文件说明 Agent 承担的角色、协作流程与产物。

## Agent 承担的角色

| 角色 | 职责 |
|------|------|
| **需求分析师** | 通过 brainstorming 技能,把模糊需求逼问为明确规范,产出 `docs/superpowers/specs/` 设计文档 |
| **架构师** | 选型(Vue+Go+Postgres)、设计"统一任务模型+双消费通道"架构,权衡 SSE 流式与异步队列的统一 |
| **计划制定者** | 通过 writing-plans 技能,把设计拆成可执行的 TDD 分步计划(`docs/superpowers/plans/`) |
| **实现工程师** | 红绿 TDD:先写失败测试 → 最小实现 → 测试通过 → 频繁语义化提交 |
| **质量把关** | 收尾时验证全量测试、合并分支(finishing-a-development-branch 技能) |

## 协作流程(superpowers 工作流)

```
brainstorming(头脑风暴) → 设计文档 → 人类确认
      ↓
writing-plans(写计划) → TDD 分步计划
      ↓
executing-plans(内联执行) → 逐 Task 红绿循环 + 提交
      ↓
finishing-a-development-branch(收尾) → 验证测试 → 合并
```

人类开发者在以下节点介入决策:大模型选型、加分项范围、技术栈(Vue+Go)、数据库(Postgres)、执行方式。Agent 负责其余的设计、实现与验证。

## 人机分工边界

- **Agent 自主**:写测试与实现、运行测试、修编译/逻辑错误、提交、装配冒烟、写文档。
- **人类决策**:产品取舍、技术选型、是否合并/发布、提供 API Key、最终验收。
- **Agent 不擅自**:在未确认设计前写实现代码(brainstorming 硬门禁);删除人类未授权的工作。

## 关键协作产物

- `docs/superpowers/specs/` — 设计规范(SDD)
- `docs/superpowers/plans/` — TDD 实现计划
- `spec/` — 需求拆分、接口设计、页面原型
- `skill.md` — CLI 工具的 Agent 可发现说明
- 语义化 git 历史 — 每个 TDD 循环一个 commit,可追溯

## 可复用的 Agent 工具链

- **后端调用**:`ai-app` CLI(见 `skill.md`)让 Agent 直接调用本服务完成翻译/总结。
- **大模型**:DeepSeek(OpenAI 兼容流式);无 Key 时自动降级为 Mock,保留真实调用链路。
