package helper

import (
	"testing"

	relaycommon "github.com/QuantumNous/new-api/relay/common"
	relayconstant "github.com/QuantumNous/new-api/relay/constant"
	"github.com/gin-gonic/gin"
)

func TestModelMappedHelperResponsesCompactPrefersFullCompactKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)
	c.Set("model_mapping", `{"gpt-5.5-openai-compact":"gpt-5.4","gpt-5.5":"gpt-5"}`)

	info := &relaycommon.RelayInfo{
		RelayMode:       relayconstant.RelayModeResponsesCompact,
		OriginModelName: "gpt-5.5-openai-compact",
		ChannelMeta: &relaycommon.ChannelMeta{
			UpstreamModelName: "gpt-5.5-openai-compact",
		},
	}

	if err := ModelMappedHelper(c, info, nil); err != nil {
		t.Fatalf("ModelMappedHelper() error = %v", err)
	}

	if !info.IsModelMapped {
		t.Fatal("ModelMappedHelper() did not mark compact full-key mapping as mapped")
	}
	if info.UpstreamModelName != "gpt-5.4" {
		t.Fatalf("UpstreamModelName = %q, want %q", info.UpstreamModelName, "gpt-5.4")
	}
	if info.OriginModelName != "gpt-5.4-openai-compact" {
		t.Fatalf("OriginModelName = %q, want %q", info.OriginModelName, "gpt-5.4-openai-compact")
	}
}

func TestModelMappedHelperResponsesCompactFallsBackToBaseKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)
	c.Set("model_mapping", `{"gpt-5.5":"gpt-5.4"}`)

	info := &relaycommon.RelayInfo{
		RelayMode:       relayconstant.RelayModeResponsesCompact,
		OriginModelName: "gpt-5.5-openai-compact",
		ChannelMeta: &relaycommon.ChannelMeta{
			UpstreamModelName: "gpt-5.5-openai-compact",
		},
	}

	if err := ModelMappedHelper(c, info, nil); err != nil {
		t.Fatalf("ModelMappedHelper() error = %v", err)
	}

	if !info.IsModelMapped {
		t.Fatal("ModelMappedHelper() did not mark compact base-key mapping as mapped")
	}
	if info.UpstreamModelName != "gpt-5.4" {
		t.Fatalf("UpstreamModelName = %q, want %q", info.UpstreamModelName, "gpt-5.4")
	}
	if info.OriginModelName != "gpt-5.4-openai-compact" {
		t.Fatalf("OriginModelName = %q, want %q", info.OriginModelName, "gpt-5.4-openai-compact")
	}
}

func TestModelMappedHelperCompactOnlyKeyDoesNotAffectNormalResponses(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)
	c.Set("model_mapping", `{"gpt-5.5-openai-compact":"gpt-5.4"}`)

	info := &relaycommon.RelayInfo{
		RelayMode:       relayconstant.RelayModeResponses,
		OriginModelName: "gpt-5.5",
		ChannelMeta: &relaycommon.ChannelMeta{
			UpstreamModelName: "gpt-5.5",
		},
	}

	if err := ModelMappedHelper(c, info, nil); err != nil {
		t.Fatalf("ModelMappedHelper() error = %v", err)
	}

	if info.IsModelMapped {
		t.Fatal("ModelMappedHelper() marked normal responses as mapped by compact-only key")
	}
	if info.UpstreamModelName != "gpt-5.5" {
		t.Fatalf("UpstreamModelName = %q, want %q", info.UpstreamModelName, "gpt-5.5")
	}
	if info.OriginModelName != "gpt-5.5" {
		t.Fatalf("OriginModelName = %q, want %q", info.OriginModelName, "gpt-5.5")
	}
}
