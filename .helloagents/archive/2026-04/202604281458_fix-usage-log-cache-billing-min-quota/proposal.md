# 变更提案: fix-usage-log-cache-billing-min-quota

## 元信息
```yaml
类型: 修复
方案类型: implementation
优先级: P1
状态: 已规划
创建: 2026-04-28
```

---

## 1. 需求

### 背景
用户反馈已在模型价格中设置输入、缓存读等价格，但使用日志的花费长期显示 `$0.000002`。该金额等价于内部 `quota=1`，说明实际计费链路可能把应计 token 算成了非正数后触发最小扣费兜底。

截图中的请求包含很小的 `prompt_tokens` 和很大的缓存读 token。该形态符合 Claude 语义：`input_tokens` 不包含 `cache_read_input_tokens`。当前部分兼容路径未标记为 `anthropic` 语义时，会按 OpenAI 语义从 `prompt_tokens` 中再次扣减缓存 token，导致基数为负并最终落到 `quota=1`。

### 目标
修复文本计费汇总逻辑，使 Claude 风格缓存 usage 即使缺少显式 `UsageSemantic` 标记，也不会因二次扣减缓存读/缓存写 token 而产生最小扣费。

### 约束条件
```yaml
时间约束: 当前修复内完成
性能约束: 只增加常量时间判断，不引入额外 I/O 或外部依赖
兼容性约束: 保持现有 OpenAI/OpenRouter OpenAI-format 缓存计费行为不变
业务约束: 不修改价格展示层、不改数据库结构、不改受保护项目标识
```

### 验收标准
- [ ] 截图同类数据形态（`cached_tokens > prompt_tokens`）不再得到 `quota=1`。
- [ ] OpenRouter OpenAI-format 缓存读场景仍按总输入拆分缓存计费。
- [ ] 现有 Claude 语义、legacy Claude-derived OpenAI usage、工具附加费用相关单测通过。

---

## 2. 方案

### 技术方案
在 `service/text_quota.go` 中增强 usage 语义识别：
- 保留显式 `usage.UsageSemantic` 和 `FinalRequestRelayFormat == RelayFormatClaude` 的优先级。
- 增加保守判定：当 usage 中存在 Claude 风格缓存字段，且 `cached_tokens + cache_creation_tokens` 大于 `prompt_tokens` 时，判为 `anthropic` 语义。
- 不基于模型名猜测，不改变缓存 token 小于等于 prompt token 的 OpenAI-format 计费路径。

同步在 `service/text_quota_test.go` 增加回归测试，覆盖用户截图中的最小扣费问题，并保留 OpenRouter OpenAI-format 用例。

### 影响范围
```yaml
涉及模块:
  - service/text_quota.go: 文本计费汇总和 usage 语义识别
  - service/text_quota_test.go: 计费语义回归测试
预计变更文件: 2 个代码文件 + 方案包/知识库记录
```

### 风险评估
| 风险 | 等级 | 应对 |
|------|------|------|
| 将 OpenAI-format 缓存 usage 误判为 Claude 语义 | 中 | 只在缓存 token 总量大于 prompt token 这种 OpenAI 总输入语义下不合理的形态触发 |
| 影响 OpenRouter 既有缓存折扣 | 中 | 保留缓存小于 prompt 的 OpenRouter 测试并运行相关单测 |
| 上游返回异常 usage 导致语义判断不准 | 低 | 仅修复明确可识别的异常形态，未知形态保持原逻辑 |

### 方案取舍
```yaml
唯一方案理由: 该方案直接修复根因，同时保持现有显式语义标记和 OpenAI-format 逻辑不变，影响面最小。
放弃的替代路径:
  - 前端花费显示兜底: 只能掩盖 quota=1 的结果，不能修复实际扣费。
  - 按模型名包含 claude 判定: 易误伤中转模型名和非 Claude 兼容模型。
  - 修改所有 Claude 兼容渠道 usage 标记: 覆盖面大且容易遗漏，仍需计费层兜底。
回滚边界: 回退 `service/text_quota.go` 的语义识别 helper 和新增测试即可恢复原行为。
```

---

## 3. 技术设计

### 计费语义判定
```text
usage.UsageSemantic 非空
  -> 使用显式语义
FinalRequestRelayFormat == RelayFormatClaude
  -> anthropic
hasClaudeStyleCacheUsage(usage)
  -> anthropic
其他
  -> openai
```

`hasClaudeStyleCacheUsage` 只在 `cached_tokens + cache_creation_tokens > prompt_tokens` 时返回 true。OpenAI 语义下 `prompt_tokens` 是总输入，缓存读/写应包含在总输入中，因此缓存 token 不应大于 prompt token。

---

## 4. 核心场景

### 场景: Claude 风格缓存读 usage 正确计费
**模块**: service/text_quota.go
**条件**: `PromptTokens=2132`，`CachedTokens=133158`，未显式标记 `UsageSemantic`
**行为**: 计费汇总将 usage 识别为 `anthropic` 语义，不从 prompt token 中扣减缓存读 token
**结果**: quota 按输入、缓存读、输出分别计费，不再落到最小值 1

---

## 5. 技术决策

### fix-usage-log-cache-billing-min-quota#D001: 使用 token 形态兜底识别 Claude 缓存语义
**日期**: 2026-04-28
**状态**: ✅采纳
**背景**: 部分 Claude 兼容链路返回的 `prompt_tokens` 不包含缓存读 token，但没有显式设置 `UsageSemantic=anthropic`，导致计费层按 OpenAI 语义重复扣减缓存。
**选项分析**:
| 选项 | 优点 | 缺点 |
|------|------|------|
| A: token 形态兜底识别 | 精准覆盖截图问题，改动小，兼容既有标记 | 只覆盖可识别的异常形态 |
| B: 按模型名识别 Claude | 覆盖更多未标记路径 | 误判风险高，模型别名不可控 |
| C: 修改所有上游适配器标记 | 语义来源更前置 | 改动范围大，无法保证第三方兼容路径全覆盖 |
**决策**: 选择方案 A
**理由**: `cached_tokens > prompt_tokens` 与 OpenAI 总输入语义冲突，但符合 Claude 独立缓存字段语义，是当前最小且可验证的修复。
**影响**: 仅影响文本计费汇总中的 usage 语义判定。

---

## 6. 验证策略

```yaml
verifyMode: test-first
reviewerFocus:
  - service/text_quota.go 中 usage 语义判定是否过宽
  - OpenRouter OpenAI-format 缓存计费是否保持不变
testerFocus:
  - go test ./service -run "TestCalculateTextQuotaSummary"
  - go test ./service -run "TestComposeTieredTextQuota"
uiValidation: none
riskBoundary:
  - 不修改价格配置、日志表结构或前端显示换算逻辑
  - 不变更 tiered billing expression 变量语义
```

---

## 7. 成果设计

N/A，非视觉任务。
