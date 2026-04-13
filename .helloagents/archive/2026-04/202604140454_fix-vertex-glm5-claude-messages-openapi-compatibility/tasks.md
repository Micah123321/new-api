# 任务清单: fix vertex glm5 claude messages openapi compatibility

> **@status:** completed | 2026-04-14 05:03

```yaml
@feature: fix vertex glm5 claude messages openapi compatibility
@created: 2026-04-14
@status: completed
@mode: R2
```

## 进度概览

| 完成 | 失败 | 跳过 | 总数 |
|------|------|------|------|
| 4 | 0 | 0 | 4 |

---

## 任务列表

### 1. Vertex 请求转换

- [√] 1.1 在 `relay/channel/vertex/adaptor.go` 中让 `RequestModeOpenSource` 的 `ConvertClaudeRequest` 改走 `service.ClaudeToOpenAIRequest`，输出 OpenAI Chat Completions payload | depends_on: []
- [√] 1.2 保持 Vertex 原生 Claude/Gemini 分支不变，并补齐必要的 `request_model` / 请求模型处理 | depends_on: [1.1]

### 2. 回归验证

- [√] 2.1 在 `relay/channel/vertex/adaptor_test.go` 增加 OpenSource 模式下 Claude Messages 转 OpenAI payload 的回归测试，覆盖 `system`、`tools/tool_choice`、多段 content、thinking | depends_on: [1.2]
- [√] 2.2 运行与本次修改直接相关的 `go test`，确认 Vertex 适配和 OpenAI 兼容逻辑通过 | depends_on: [2.1]

---

## 执行日志

| 时间 | 任务 | 状态 | 备注 |
|------|------|------|------|
| 2026-04-14 04:58:12 | 1.1 | completed | Vertex OpenSource 分支改为复用 `service.ClaudeToOpenAIRequest` |
| 2026-04-14 04:59:04 | 1.2 | completed | 保持 Vertex Claude/Gemini 原有分支不变 |
| 2026-04-14 05:00:08 | 2.1 | completed | 新增 Vertex 和 service 定向回归测试 |
| 2026-04-14 05:01:32 | 2.2 | completed | `go test ./relay/channel/vertex ./relay/channel/openai ./service -run 'Test(ConvertClaudeRequest_|ClaudeToOpenAIRequest|SendStreamData_)'` 通过 |

---

## 执行备注

- D001: OpenSource Vertex 复用通用 `ClaudeToOpenAIRequest` 转换链，而不是新增 Vertex 专用重复映射。
- 说明: `go test ./service` 全量执行时仍有与本次修改无关的既有失败：`TestObserveChannelAffinityUsageCacheByRelayFormat_MixedMode` 与 `TestObserveChannelAffinityUsageCacheByRelayFormat_UnsupportedModeKeepsEmpty`。
