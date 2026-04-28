# 方案归档索引

> 通过此文件快速查找历史方案
> 历史年份: [2024](_index-2024.md) | [2023](_index-2023.md) | ...

## 快速索引（当前年份）

| 时间戳 | 名称 | 类型 | 涉及模块 | 决策 | 结果 |
|--------|------|------|---------|------|------|
| 202604281458 | fix-usage-log-cache-billing-min-quota | implementation | service/text-quota | fix-usage-log-cache-billing-min-quota#D001 | ✅完成 |
| 202604241556 | add-ghcr-trigger-for-sync-origin-main-local7-20260318 | implementation | ci/ghcr-workflow | add-ghcr-trigger-for-sync-origin-main-local7-20260318#D001 | ✅完成 |
| 202604161906 | fix-stripe-webhook-empty-secret-bypass | - | - | - | ✅完成 |
| 202604140459 | sync-fork-with-upstream-preserving-local-changes | - | - | - | ✅完成 |
| 202604140454 | fix-vertex-glm5-claude-messages-openapi-compatibility | implementation | relay/channel/vertex, service/convert | fix-vertex-glm5-claude-messages-openapi-compatibility#D001 | ✅完成 |

## 按月归档

### 2026-04
- [202604281458_fix-usage-log-cache-billing-min-quota](./2026-04/202604281458_fix-usage-log-cache-billing-min-quota/) - 修复 Claude 风格缓存 usage 未显式标记语义时重复扣减缓存 token，导致使用日志落到最小扣费的问题
- [202604241556_add-ghcr-trigger-for-sync-origin-main-local7-20260318](./2026-04/202604241556_add-ghcr-trigger-for-sync-origin-main-local7-20260318/) - 为同步分支增加 GHCR 自动构建触发，并隔离分支镜像标签，避免覆盖主线 `main/latest`
- [202604161906_fix-stripe-webhook-empty-secret-bypass](./2026-04/202604161906_fix-stripe-webhook-empty-secret-bypass/) - 修复 Stripe webhook 在空密钥配置下可被伪造完成充值或订阅的问题
- [202604140454_fix-vertex-glm5-claude-messages-openapi-compatibility](./2026-04/202604140454_fix-vertex-glm5-claude-messages-openapi-compatibility/) - 修复 Vertex OpenSource `-maas` 模型对 Claude Messages 请求的 OpenAI 兼容转换

## 结果状态说明
- ✅ 完成
- ⚠️ 部分完成
- ❌ 失败/中止
- ⏸ 未执行
- 🔄 已回滚
- 📄 概述
