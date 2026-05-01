# CHANGELOG

## [v0.12.13] - 2026-05-01

### 快速修改
- **[relay/helper]**: 修复 Responses compact 模式模型重定向先剥离 `-openai-compact` 导致完整 compact key 无法命中的问题，并保持基础模型映射兼容 — by yinjianm
  - 类型: 快速修改（无方案包）
  - 文件: relay/helper/model_mapped.go:16-90, relay/helper/model_mapped_test.go:10-93
- **[middleware/distributor]**: 修复 Responses compact 普通请求复用渠道缓存时可能拿不到最新 `model_mapping`，导致实际请求仍按原模型转发和计费的问题 — by yinjianm
  - 类型: 快速修改（无方案包）
  - 文件: middleware/distributor.go:345-376, middleware/distributor_test.go:36-62, model/channel.go:358-366

## [v0.12.12] - 2026-04-28

### 修复
- **[service/text-quota]**: 修复 Claude 风格缓存 usage 未显式标记语义时被按 OpenAI 语义二次扣减缓存 token，导致使用日志花费落到 `$0.000002` 最小扣费的问题 — by yinjianm
  - 方案: [202604281458_fix-usage-log-cache-billing-min-quota](archive/2026-04/202604281458_fix-usage-log-cache-billing-min-quota/)
  - 决策: fix-usage-log-cache-billing-min-quota#D001(使用 token 形态兜底识别 Claude 缓存语义)

## [工作区修复] - 2026-04-24

### 快速修改
- **[relay/tool-billing]**: 修复同步分支在合并上游后遗漏工具计费接口迁移与 channel test 计费辅助函数，消除 Docker `go build` 中的 `operation_setting` 未定义符号和本地包级编译阻断 [快速修改] [文件: relay/compatible_handler.go:357-392, controller/channel-test.go:559-603]

## [v0.12.11] - 2026-04-24

### 修复
- **[ci/ghcr-workflow]**: 为 `sync/origin-main-with-local7-20260318` 增加 GHCR 自动构建触发，并将同步分支镜像标签隔离为 `origin-main-with-local7-20260318` / `latest-origin-main-with-local7-20260318`，避免覆盖主线 `main/latest` — by yinjianm
  - 方案: [202604241556_add-ghcr-trigger-for-sync-origin-main-local7-20260318](archive/2026-04/202604241556_add-ghcr-trigger-for-sync-origin-main-local7-20260318/)

## [工作区修复] - 2026-04-22

### 快速修改
- **[controller/stripe]**: 修复 Stripe webhook 在 `stripeWebhookSecretConfigured()` 为 false 时误用未导入的标准库 `log.Printf`，导致 Docker `go build` 编译失败 [快速修改] [文件: controller/topup_stripe.go:170-173]

## [v0.12.10] - 2026-04-16

### 修复
- **[controller/stripe]**: 修复 Stripe webhook 在 `StripeWebhookSecret` 为空时仍可进入验签与订单完成逻辑的问题，并阻止未配置状态下继续创建 Stripe 支付订单 — by yinjianm
  - 方案: [202604161906_fix-stripe-webhook-empty-secret-bypass](archive/2026-04/202604161906_fix-stripe-webhook-empty-secret-bypass/)

## [0.0.1] - 2026-04-14

### 修复
- **[relay/channel/vertex]**: 修复 Vertex OpenSource `-maas` 模型处理 `/v1/messages` 时错误发送 Anthropic payload 的问题，并补齐 Claude `tool_choice` 到 OpenAI 请求的兼容映射 — by yinjianm
  - 方案: [202604140454_fix-vertex-glm5-claude-messages-openapi-compatibility](archive/2026-04/202604140454_fix-vertex-glm5-claude-messages-openapi-compatibility/)
