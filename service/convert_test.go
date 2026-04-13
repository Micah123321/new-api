package service

import (
	"testing"

	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/dto"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/stretchr/testify/require"
)

func TestClaudeToOpenAIRequestConvertsNamedToolChoice(t *testing.T) {
	t.Parallel()

	info := &relaycommon.RelayInfo{
		OriginModelName: "zai-org/glm-5-maas",
		ChannelMeta: &relaycommon.ChannelMeta{
			ChannelType:       constant.ChannelTypeVertexAi,
			UpstreamModelName: "zai-org/glm-5-maas",
		},
	}

	req := dto.ClaudeRequest{
		Model: "zai-org/glm-5-maas",
		Messages: []dto.ClaudeMessage{
			{Role: "user", Content: "hello"},
		},
		ToolChoice: dto.ClaudeToolChoice{
			Type:                   "tool",
			Name:                   "lookup",
			DisableParallelToolUse: true,
		},
	}

	converted, err := ClaudeToOpenAIRequest(req, info)
	require.NoError(t, err)

	toolChoice, ok := converted.ToolChoice.(map[string]any)
	require.True(t, ok)
	require.Equal(t, "function", toolChoice["type"])

	functionChoice, ok := toolChoice["function"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "lookup", functionChoice["name"])

	require.NotNil(t, converted.ParallelTooCalls)
	require.False(t, *converted.ParallelTooCalls)
}

func TestClaudeToOpenAIRequestConvertsAnyToolChoiceToRequired(t *testing.T) {
	t.Parallel()

	info := &relaycommon.RelayInfo{
		OriginModelName: "zai-org/glm-5-maas",
		ChannelMeta: &relaycommon.ChannelMeta{
			ChannelType:       constant.ChannelTypeVertexAi,
			UpstreamModelName: "zai-org/glm-5-maas",
		},
	}

	req := dto.ClaudeRequest{
		Model: "zai-org/glm-5-maas",
		Messages: []dto.ClaudeMessage{
			{Role: "user", Content: "hello"},
		},
		ToolChoice: dto.ClaudeToolChoice{
			Type: "any",
		},
	}

	converted, err := ClaudeToOpenAIRequest(req, info)
	require.NoError(t, err)
	require.Equal(t, "required", converted.ToolChoice)
	require.Nil(t, converted.ParallelTooCalls)
}
