package vertex

import (
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/dto"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func textPtr(v string) *string {
	return &v
}

func TestGetRequestURL_OpenSourceUsesV1Endpoint(t *testing.T) {
	t.Parallel()

	info := &relaycommon.RelayInfo{
		OriginModelName: "zai-org/glm-5-maas",
		ChannelMeta: &relaycommon.ChannelMeta{
			ApiVersion:        `{"default":"global"}`,
			ApiKey:            `{"project_id":"demo-project"}`,
			UpstreamModelName: "zai-org/glm-5-maas",
			ChannelOtherSettings: dto.ChannelOtherSettings{
				VertexKeyType: dto.VertexKeyTypeJSON,
			},
		},
	}

	adaptor := &Adaptor{}
	adaptor.Init(info)

	url, err := adaptor.GetRequestURL(info)
	require.NoError(t, err)
	require.Equal(t, "https://aiplatform.googleapis.com/v1/projects/demo-project/locations/global/endpoints/openapi/chat/completions", url)
}

func TestGetRequestURL_OpenSourceWithAPIKeyReturnsError(t *testing.T) {
	t.Parallel()

	info := &relaycommon.RelayInfo{
		OriginModelName: "zai-org/glm-5-maas",
		ChannelMeta: &relaycommon.ChannelMeta{
			ApiVersion:        `{"default":"global"}`,
			ApiKey:            "AIza-example",
			UpstreamModelName: "zai-org/glm-5-maas",
			ChannelOtherSettings: dto.ChannelOtherSettings{
				VertexKeyType: dto.VertexKeyTypeAPIKey,
			},
		},
	}

	adaptor := &Adaptor{}
	adaptor.Init(info)

	_, err := adaptor.GetRequestURL(info)
	require.Error(t, err)
	require.Contains(t, err.Error(), "service account json credentials")
}

func TestConvertClaudeRequest_OpenSourceUsesOpenAICompatPayload(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	info := &relaycommon.RelayInfo{
		OriginModelName: "zai-org/glm-5-maas",
		ChannelMeta: &relaycommon.ChannelMeta{
			UpstreamModelName: "zai-org/glm-5-maas",
			ChannelOtherSettings: dto.ChannelOtherSettings{
				VertexKeyType: dto.VertexKeyTypeJSON,
			},
		},
	}

	adaptor := &Adaptor{}
	adaptor.Init(info)
	require.Equal(t, RequestModeOpenSource, adaptor.RequestMode)

	req := &dto.ClaudeRequest{
		Model: "zai-org/glm-5-maas",
		System: []dto.ClaudeMediaMessage{
			{Type: "text", Text: textPtr("system instruction")},
		},
		Messages: []dto.ClaudeMessage{
			{
				Role: "user",
				Content: []dto.ClaudeMediaMessage{
					{Type: "text", Text: textPtr("first")},
					{Type: "text", Text: textPtr(" second")},
				},
			},
		},
		MaxTokens: lo.ToPtr(uint(256)),
		Stream:    lo.ToPtr(true),
		Tools: []any{
			dto.Tool{
				Name:        "lookup",
				Description: "Find data",
				InputSchema: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"q": map[string]any{"type": "string"},
					},
				},
			},
		},
		ToolChoice: dto.ClaudeToolChoice{
			Type:                   "tool",
			Name:                   "lookup",
			DisableParallelToolUse: true,
		},
	}

	converted, err := adaptor.ConvertClaudeRequest(ctx, info, req)
	require.NoError(t, err)

	openAIReq, ok := converted.(*dto.GeneralOpenAIRequest)
	require.True(t, ok, "expected OpenAI request payload for Vertex OpenSource mode")
	require.Equal(t, "zai-org/glm-5-maas", openAIReq.Model)
	require.Equal(t, "zai-org/glm-5-maas", ctx.GetString("request_model"))
	require.Equal(t, lo.ToPtr(uint(256)), openAIReq.MaxTokens)
	require.Equal(t, lo.ToPtr(true), openAIReq.Stream)

	require.Len(t, openAIReq.Messages, 2)
	require.Equal(t, "system", openAIReq.Messages[0].Role)
	require.Equal(t, "system instruction", openAIReq.Messages[0].StringContent())
	require.Equal(t, "user", openAIReq.Messages[1].Role)
	require.Len(t, openAIReq.Messages[1].ParseContent(), 2)
	require.Equal(t, "first", openAIReq.Messages[1].ParseContent()[0].Text)
	require.Equal(t, " second", openAIReq.Messages[1].ParseContent()[1].Text)

	require.Len(t, openAIReq.Tools, 1)
	require.Equal(t, "lookup", openAIReq.Tools[0].Function.Name)
	toolChoice, ok := openAIReq.ToolChoice.(map[string]any)
	require.True(t, ok, "expected named tool_choice to be converted into OpenAI function selection")
	require.Equal(t, "function", toolChoice["type"])
	functionChoice, ok := toolChoice["function"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "lookup", functionChoice["name"])
	require.NotNil(t, openAIReq.ParallelTooCalls)
	require.False(t, *openAIReq.ParallelTooCalls)
}

func TestConvertClaudeRequest_ClaudeModeKeepsVertexAnthropicPayload(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	info := &relaycommon.RelayInfo{
		OriginModelName: "claude-3-5-sonnet-20240620",
		ChannelMeta: &relaycommon.ChannelMeta{
			UpstreamModelName: "claude-3-5-sonnet-20240620",
			ChannelOtherSettings: dto.ChannelOtherSettings{
				VertexKeyType: dto.VertexKeyTypeJSON,
			},
		},
	}

	adaptor := &Adaptor{}
	adaptor.Init(info)
	require.Equal(t, RequestModeClaude, adaptor.RequestMode)

	req := &dto.ClaudeRequest{
		Model: "claude-3-5-sonnet-20240620",
		Messages: []dto.ClaudeMessage{
			{Role: "user", Content: "hello"},
		},
	}

	converted, err := adaptor.ConvertClaudeRequest(ctx, info, req)
	require.NoError(t, err)

	vertexReq, ok := converted.(*VertexAIClaudeRequest)
	require.True(t, ok, "expected Vertex Anthropic payload for Claude mode")
	require.Equal(t, anthropicVersion, vertexReq.AnthropicVersion)
	require.Len(t, vertexReq.Messages, 1)
	require.Equal(t, "claude-3-5-sonnet@20240620", ctx.GetString("request_model"))
}
