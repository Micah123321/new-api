# CHANGELOG

## [v0.12.10] - 2026-04-16

### 修复
- **[controller/stripe]**: 修复 Stripe webhook 在 `StripeWebhookSecret` 为空时仍可进入验签与订单完成逻辑的问题，并阻止未配置状态下继续创建 Stripe 支付订单 — by yinjianm
  - 方案: [202604161906_fix-stripe-webhook-empty-secret-bypass](archive/2026-04/202604161906_fix-stripe-webhook-empty-secret-bypass/)

## [0.0.1] - 2026-04-14

### 修复
- **[relay/channel/vertex]**: 修复 Vertex OpenSource `-maas` 模型处理 `/v1/messages` 时错误发送 Anthropic payload 的问题，并补齐 Claude `tool_choice` 到 OpenAI 请求的兼容映射 — by yinjianm
  - 方案: [202604140454_fix-vertex-glm5-claude-messages-openapi-compatibility](archive/2026-04/202604140454_fix-vertex-glm5-claude-messages-openapi-compatibility/)
