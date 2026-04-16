# 任务清单: fix-stripe-webhook-empty-secret-bypass

> **@status:** completed | 2026-04-16 19:12

```yaml
@feature: fix-stripe-webhook-empty-secret-bypass
@created: 2026-04-16
@status: completed
@mode: R2
```

## 进度概览

| 完成 | 失败 | 跳过 | 总数 |
|------|------|------|------|
| 4 | 0 | 0 | 4 |

---

## 任务列表

### 1. 控制器修复

- [√] 1.1 在 `controller/topup_stripe.go` 中为 `StripeWebhook` 增加空/空白 `StripeWebhookSecret` 拒绝逻辑，并复用到普通 Stripe 支付入口的前置校验 | depends_on: []
- [√] 1.2 在 `controller/subscription_payment_stripe.go` 中补充相同的 webhook secret 配置校验，阻止未配置状态下创建订阅支付订单 | depends_on: [1.1]

### 2. 验证与回归

- [√] 2.1 新增或扩展控制器测试，覆盖空 webhook secret 下的 webhook 拒绝、普通充值拒绝、订阅支付拒绝 | depends_on: [1.1, 1.2]
- [√] 2.2 运行与本次变更相关的 Go 测试，确认修复行为与编译状态 | depends_on: [2.1]

---

## 执行日志

| 时间 | 任务 | 状态 | 备注 |
|------|------|------|------|
| 2026-04-16 19:07:00 | 1.1 | completed | 已在普通充值入口、webhook 入口和内部建链函数补充 webhook secret 判空防护 |
| 2026-04-16 19:08:00 | 1.2 | completed | 已在订阅支付入口和内部建链函数补充相同配置校验 |
| 2026-04-16 19:10:00 | 2.1 | completed | 新增 `controller/stripe_security_test.go` 覆盖空 secret 拒绝场景 |
| 2026-04-16 19:11:18 | 2.2 | completed | `go test ./controller -run 'TestStripeWebhookRejectsEmptyWebhookSecret|TestRequestStripePayRejectsMissingWebhookSecret|TestSubscriptionRequestStripePayRejectsMissingWebhookSecret'` 通过 |

---

## 执行备注

> 记录执行过程中的重要说明、决策变更、风险提示等

- 本次修复聚焦空 webhook secret 绕过链路，未扩展到完整支付对账或事件字段二次校验。
- 普通充值接口保持原有响应结构：失败时返回 `message=error` 与中文错误文本。
