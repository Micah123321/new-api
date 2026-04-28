# service/text_quota.go 模块说明

## 职责

- 汇总文本请求的输入、输出、缓存读、缓存写、图片、音频和工具调用附加费用。
- 根据 `types.PriceData` 中的模型倍率、分组倍率、缓存倍率和固定价格计算最终 `quota`。
- 将计费结果写入用户、渠道和使用日志，并为日志详情生成价格、缓存和 billing metadata。

## 行为规范

- `UsageSemantic=anthropic` 表示 Claude 语义：`PromptTokens` 是非缓存输入 token，缓存读/缓存写 token 独立计费，不应从 prompt 中再次扣减。
- `UsageSemantic=openai` 表示 OpenAI 语义：`PromptTokens` 是总输入 token，缓存读/缓存写 token 需要从普通输入中拆分出来按缓存倍率计费。
- 当上游未显式标记语义，但 `cached_tokens + cache_creation_tokens > prompt_tokens` 时，计费层按 Claude 风格缓存语义处理；该 token 形态在 OpenAI 总输入语义下不合理，会导致重复扣减缓存并落到最小扣费。
- 显式 `usage.UsageSemantic` 优先级最高，`FinalRequestRelayFormat == RelayFormatClaude` 次之；兜底识别不得覆盖显式语义。
- OpenRouter OpenAI-format 缓存用例中缓存 token 小于等于 prompt token，仍按 OpenAI 总输入语义拆分计费。

## 依赖关系

- `dto.Usage`: 上游 usage 结构和缓存 token 字段来源。
- `relay/common.RelayInfo`: 请求格式、渠道信息、模型价格和分组倍率来源。
- `types.PriceData`: 模型倍率、缓存倍率、固定价格和附加倍率来源。
- `pkg/billingexpr`: tiered billing expression 的预扣费与结算变量体系。
- `model.RecordConsumeLog`: 使用日志落库。

## 验证覆盖

- `go test ./service -run "TestCalculateTextQuotaSummary"`
- `go test ./service -run "TestComposeTieredTextQuota"`
