package codex

import (
	"testing"

	"github.com/QuantumNous/new-api/dto"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	relayconstant "github.com/QuantumNous/new-api/relay/constant"
)

func TestConvertOpenAIResponsesRequest_CompactKeepsInstructionsOptional(t *testing.T) {
	adaptor := &Adaptor{}
	info := &relaycommon.RelayInfo{
		RelayMode: relayconstant.RelayModeResponsesCompact,
		ChannelMeta: &relaycommon.ChannelMeta{
			ChannelSetting: dto.ChannelSettings{},
		},
	}

	convertedAny, err := adaptor.ConvertOpenAIResponsesRequest(nil, info, dto.OpenAIResponsesRequest{
		Model: "gpt-5",
	})
	if err != nil {
		t.Fatalf("ConvertOpenAIResponsesRequest() error = %v", err)
	}

	converted, ok := convertedAny.(dto.OpenAIResponsesRequest)
	if !ok {
		t.Fatalf("ConvertOpenAIResponsesRequest() type = %T, want dto.OpenAIResponsesRequest", convertedAny)
	}
	if len(converted.Instructions) != 0 {
		t.Fatalf("compact request should not auto inject empty instructions, got %q", string(converted.Instructions))
	}
	if len(converted.Store) != 0 {
		t.Fatalf("compact request should not force store=false, got %q", string(converted.Store))
	}
}

func TestConvertOpenAIResponsesRequest_ResponsesInjectsDefaultInstructions(t *testing.T) {
	adaptor := &Adaptor{}
	info := &relaycommon.RelayInfo{
		RelayMode: relayconstant.RelayModeResponses,
		ChannelMeta: &relaycommon.ChannelMeta{
			ChannelSetting: dto.ChannelSettings{},
		},
	}

	convertedAny, err := adaptor.ConvertOpenAIResponsesRequest(nil, info, dto.OpenAIResponsesRequest{
		Model: "gpt-5",
	})
	if err != nil {
		t.Fatalf("ConvertOpenAIResponsesRequest() error = %v", err)
	}

	converted, ok := convertedAny.(dto.OpenAIResponsesRequest)
	if !ok {
		t.Fatalf("ConvertOpenAIResponsesRequest() type = %T, want dto.OpenAIResponsesRequest", convertedAny)
	}
	if string(converted.Instructions) != `""` {
		t.Fatalf("responses request should auto inject empty instructions, got %q", string(converted.Instructions))
	}
	if string(converted.Store) != "false" {
		t.Fatalf("responses request should force store=false, got %q", string(converted.Store))
	}
}
