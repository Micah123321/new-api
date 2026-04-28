# 项目上下文

## 基本信息

- 项目类型: Go 编写的 AI API 网关 / 代理服务
- 核心技术栈: Gin, GORM, React, Vite, Redis
- 当前修复焦点: Claude 风格缓存 usage 未显式标记语义时的文本计费最小扣费修复

## 当前约束

- 需要同时兼容 OpenAI、Claude、Gemini 三类请求格式
- Vertex AI 渠道同时存在原生 Claude、Gemini 和 OpenAPI OpenSource 三条路径
- 文本计费链路需要区分 OpenAI 总输入语义与 Claude 独立缓存读/写语义，避免重复扣减缓存 token
- 修复时不得影响现有 tiered billing expression、日志表结构和受保护项目标识
