# relay/channel/vertex

## 职责

- 根据 `UpstreamModelName` 将 Vertex 请求分发到 Claude、Gemini 或 OpenSource(OpenAPI) 模式
- 组装对应上游 URL、请求头和请求体
- 处理 Vertex OpenAPI 返回的 OpenAI 风格响应，以及 Vertex 原生 Claude/Gemini 响应

## 本次更新

- `RequestModeOpenSource` 下的 `ConvertClaudeRequest` 不再返回 Vertex Anthropic payload
- `-maas` 模型收到 `/v1/messages` 时，改为复用 `service.ClaudeToOpenAIRequest` 生成 OpenAI Chat Completions 请求体
- 原生 Claude 与 Gemini 模式保持原有分支和 URL 逻辑不变

## 依赖关系

- 依赖 `service/convert.go` 完成 Claude -> OpenAI 的协议转换
- 依赖 `relay/channel/openai` 处理 OpenSource 上游响应

