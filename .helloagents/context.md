# 项目上下文

## 基本信息

- 项目类型: Go 编写的 AI API 网关 / 代理服务
- 核心技术栈: Gin, GORM, React, Vite, Redis
- 当前修复焦点: Stripe webhook 空密钥绕过支付验证修复

## 当前约束

- 需要同时兼容 OpenAI、Claude、Gemini 三类请求格式
- Vertex AI 渠道同时存在原生 Claude、Gemini 和 OpenAPI OpenSource 三条路径
- Stripe 支付链路中，配置项校验必须先于 webhook 验签和订单完成逻辑
- 修复时不得影响现有计费、日志和受保护项目标识
