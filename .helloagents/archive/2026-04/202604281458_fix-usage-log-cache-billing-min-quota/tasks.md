# 任务清单: fix-usage-log-cache-billing-min-quota

> **@status:** completed | 2026-04-28 15:04

```yaml
@feature: fix-usage-log-cache-billing-min-quota
@created: 2026-04-28
@status: completed
@mode: R2
```

## LIVE_STATUS
```json
{"status":"completed","completed":3,"failed":0,"pending":0,"total":3,"percent":100,"current":"全部任务完成，准备归档方案包","updated_at":"2026-04-28 15:04:08"}
```

## 进度概览

| 完成 | 失败 | 跳过 | 总数 |
|------|------|------|------|
| 3 | 0 | 0 | 3 |

---

## 任务列表

### 1. 计费语义修复

- [√] 1.1 修改 `service/text_quota.go` 的 usage 语义识别
  - 预期变更: 增加 Claude 风格缓存 usage 的保守兜底识别，避免缓存 token 大于 prompt token 时仍按 OpenAI 语义二次扣减。
  - 完成标准: 截图同类数据形态被识别为 `anthropic` 语义，最终 quota 不再为 1。
  - 验证方式: 新增单测覆盖 `cached_tokens > prompt_tokens` 的场景。
  - depends_on: []

### 2. 回归测试

- [√] 2.1 更新 `service/text_quota_test.go`
  - 预期变更: 增加用户截图同类用例，并确保既有 OpenRouter OpenAI-format 缓存拆分用例仍覆盖。
  - 完成标准: 新用例断言 quota 按输入、缓存读、输出计费，`IsClaudeUsageSemantic=true`。
  - 验证方式: `go test ./service -run "TestCalculateTextQuotaSummary"`
  - depends_on: [1.1]

### 3. 验证与知识库同步

- [√] 3.1 运行相关测试并更新知识库记录
  - 预期变更: 执行相关 Go 单测；更新 `.helloagents/CHANGELOG.md`，并将方案包状态同步为完成。
  - 完成标准: 测试通过或明确记录失败原因；知识库记录包含本次修复和方案包链接。
  - 验证方式: 检查测试输出、方案包 `LIVE_STATUS`、CHANGELOG 记录。
  - depends_on: [2.1]

---

## 执行日志

| 时间 | 任务 | 状态 | 备注 |
|------|------|------|------|
| 2026-04-28 14:58:00 | 方案设计 | in_progress | 已创建方案包并确定唯一修复方案 |
| 2026-04-28 15:03:09 | 1.1 | completed | 已增加 Claude 风格缓存 usage 兜底识别 |
| 2026-04-28 15:03:09 | 2.1 | completed | 已增加缓存 token 大于 prompt token 的回归测试 |
| 2026-04-28 15:04:08 | 3.1 | completed | 相关 service 单测通过，知识库记录已更新；完整 `go test ./service` 仍受 channel affinity usage cache 既有测试隔离问题影响 |

---

## 执行备注

- 本次修复不修改价格配置、日志表结构和前端金额换算。
- `pkg/billingexpr/expr.md` 已按项目规则阅读；本次不改变 tiered billing expression 变量语义。
- 验证说明: `go test ./service -run "TestCalculateTextQuotaSummary|TestComposeTieredTextQuota"` 通过；`go test ./service` 失败于 `TestObserveChannelAffinityUsageCacheByRelayFormat_UnsupportedModeKeepsEmpty`，单独运行同组用例也失败，属于本次计费修复范围外的全局指标测试隔离问题。
