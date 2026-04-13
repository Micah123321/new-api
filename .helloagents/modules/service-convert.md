# service/convert.go

## 职责

- 负责 Claude、OpenAI、Gemini 之间的请求与响应转换
- 为多渠道适配层提供统一的协议互转能力

## 本次更新

- 新增 Claude `tool_choice` 到 OpenAI `tool_choice` 的兼容映射
- 当 Claude 请求显式禁止并行工具调用时，映射到 OpenAI `parallel_tool_calls=false`

## 依赖关系

- 被 `relay/channel/openai`、`relay/channel/vertex` 等多条渠道转换链复用

