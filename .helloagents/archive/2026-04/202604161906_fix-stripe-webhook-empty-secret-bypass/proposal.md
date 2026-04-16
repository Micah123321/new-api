# 变更提案: fix-stripe-webhook-empty-secret-bypass

## 元信息
```yaml
类型: 修复
方案类型: implementation
优先级: P0
状态: 已确认
创建: 2026-04-16
```

---

## 1. 需求

### 背景
当前 `StripeWebhook` 入口会直接使用 `setting.StripeWebhookSecret` 参与 Stripe SDK 验签，未在业务层显式拒绝空密钥。若实例误配置为空字符串，攻击者可伪造 `checkout.session.completed` 回调，将本地待支付订单直接标记为成功，导致未真实支付却完成充值或订阅开通。

### 目标
- 在 Stripe webhook 入口处阻断空密钥配置下的所有回调处理。
- 在发起 Stripe 支付或订阅支付前补充 webhook 配置校验，避免继续创建可被伪造完成的待支付订单。
- 增加最小回归测试，覆盖充值与订阅两条入口的关键防护行为。

### 约束条件
```yaml
时间约束: 本次仅做服务端安全修复，不扩展为完整支付对账重构
性能约束: 不新增额外外部依赖，不引入明显请求开销
兼容性约束: 保持现有 Stripe 集成方式、数据库模型和三数据库兼容性不变
业务约束: 不修改受保护的项目标识，不改变正常已配置实例的支付流程
```

### 验收标准
- [ ] `StripeWebhookSecret` 为空或仅空白时，`/api/stripe/webhook` 返回拒绝结果，不再进入事件处理逻辑
- [ ] `RequestStripePay` 与 `SubscriptionRequestStripePay` 在 webhook secret 未配置时拒绝创建支付会话/订单
- [ ] 新增测试可验证上述保护逻辑，且不破坏现有代码编译

---

## 2. 方案

### 技术方案
在 `controller/topup_stripe.go` 的 `StripeWebhook` 入口增加配置校验，先对 `StripeWebhookSecret` 做 `TrimSpace` 并显式判空；为空时直接返回错误状态并记录日志，不再调用 Stripe SDK 验签。

在充值和订阅两条 Stripe 支付创建入口补充同样的 webhook secret 校验，使系统在配置不完整时无法继续创建待支付订单。

新增控制器层最小回归测试，使用 `httptest` 验证：
- webhook 在空 secret 下直接拒绝；
- 充值支付入口在空 secret 下直接返回业务错误；
- 订阅支付入口在空 secret 下直接返回业务错误。

### 影响范围
```yaml
涉及模块:
  - controller/topup_stripe.go: Stripe webhook 与普通充值支付入口加固
  - controller/subscription_payment_stripe.go: 订阅支付入口加固
  - controller/*_test.go: 新增或扩展回归测试
预计变更文件: 3
```

### 风险评估
| 风险 | 等级 | 应对 |
|------|------|------|
| 误伤已配置但含前后空白的 secret | 低 | 校验前先 `TrimSpace`，仅对实际空值拒绝 |
| 测试依赖现有数据库/上下文导致不稳定 | 中 | 优先选取无需真实 Stripe/数据库写入的短路径校验，缩小测试面 |
| 仅修复空密钥绕过，未覆盖更强的支付对账问题 | 中 | 本次限定为止血修复，后续可单独规划对账与事件字段校验增强 |

---

## 3. 技术设计（可选）

> 本次不涉及架构、API 或数据模型变更，N/A。

### 架构设计
N/A

### API设计
N/A

### 数据模型
| 字段 | 类型 | 说明 |
|------|------|------|
| N/A | N/A | 本次无数据模型变更 |

---

## 4. 核心场景

> 执行完成后同步到对应模块文档

### 场景: 空 webhook secret 下的 Stripe 回调保护
**模块**: controller/topup_stripe.go
**条件**: `StripeWebhookSecret` 未配置或为空白字符串
**行为**: 收到 `/api/stripe/webhook` 请求
**结果**: 服务端直接拒绝请求，不调用 Stripe SDK，不执行充值或订阅完成逻辑

### 场景: 配置不完整时阻止创建 Stripe 待支付订单
**模块**: controller/topup_stripe.go, controller/subscription_payment_stripe.go
**条件**: `StripeApiSecret` 可用，但 `StripeWebhookSecret` 未配置
**行为**: 用户请求普通充值或订阅支付链接
**结果**: 服务端返回业务错误，不创建新的待支付订单

---

## 5. 技术决策

> 本方案涉及的技术决策，归档后成为决策的唯一完整记录

### fix-stripe-webhook-empty-secret-bypass#D001: 在业务入口显式拒绝空 webhook secret
**日期**: 2026-04-16
**状态**: ✅采纳
**背景**: Stripe SDK 的 HMAC 验签对空字符串 secret 不会自动失败，若业务层不先拒绝空配置，则“未配置”会退化成“使用空密钥验签”，导致伪造回调可被接受。
**选项分析**:
| 选项 | 优点 | 缺点 |
|------|------|------|
| A: 仅在 webhook 入口判空 | 能快速止血，改动最小 | 仍允许前台继续创建待支付 Stripe 订单，运维误配置下风险持续存在 |
| B: 在 webhook 与支付创建入口都判空 | 同时阻断伪造回调与新订单生成，行为更一致 | 改动面略大，需要补两类入口测试 |
**决策**: 选择方案 B
**理由**: 这是对现有业务影响最小、但能完整封住当前利用链的方案；无需改数据模型，也不引入额外依赖。
**影响**: 影响 Stripe 充值控制器、Stripe 订阅支付控制器及相关测试

---

## 6. 成果设计

> 含视觉产出的任务由 DESIGN Phase2 填充。非视觉任务整节标注"N/A"。

N/A
