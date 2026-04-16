# 模块: controller-stripe-payment

## 职责

- 负责 Stripe 普通充值与订阅支付的控制器入口校验。
- 负责处理 `/api/stripe/webhook` 回调的前置验签与事件分发。
- 在配置不完整时阻断支付创建与回调处理，避免误配置导致订单被伪造完成。

## 行为规范

- `StripeWebhookSecret` 必须在业务层显式校验，空值或仅空白字符串均视为未配置。
- 普通 Stripe 充值与订阅支付在 webhook secret 未配置时不得创建新的待支付订单。
- Stripe webhook 在 secret 未配置时直接拒绝请求，不进入 Stripe SDK 验签或本地订单完成逻辑。
- 充值和订阅支付入口保持现有响应风格，不改变正常已配置实例的业务路径。

## 依赖关系

- 依赖 `setting` 读取 Stripe 配置项。
- 依赖 `model` 执行充值订单、订阅计划与订阅订单相关查询。
- 依赖 `common` 提供 JSON、响应封装和通用校验能力。
