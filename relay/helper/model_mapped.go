package helper

import (
	"errors"
	"fmt"
	"strings"

	appcommon "github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/relay/common"
	relayconstant "github.com/QuantumNous/new-api/relay/constant"
	"github.com/QuantumNous/new-api/setting/ratio_setting"
	"github.com/gin-gonic/gin"
)

func ModelMappedHelper(c *gin.Context, info *common.RelayInfo, request dto.Request) error {
	if info.ChannelMeta == nil {
		info.ChannelMeta = &common.ChannelMeta{}
	}

	isResponsesCompact := info.RelayMode == relayconstant.RelayModeResponsesCompact
	originModelName := info.OriginModelName
	baseModelName := originModelName
	if isResponsesCompact && strings.HasSuffix(originModelName, ratio_setting.CompactModelSuffix) {
		baseModelName = strings.TrimSuffix(originModelName, ratio_setting.CompactModelSuffix)
	}

	// map model name
	modelMapping := c.GetString("model_mapping")
	if modelMapping != "" && modelMapping != "{}" {
		modelMap := make(map[string]string)
		err := appcommon.Unmarshal([]byte(modelMapping), &modelMap)
		if err != nil {
			return fmt.Errorf("unmarshal_model_mapping_failed")
		}

		for _, candidate := range modelMappingCandidates(originModelName, baseModelName, isResponsesCompact) {
			mappedModel, isMapped, err := resolveModelMapping(modelMap, candidate)
			if err != nil {
				return err
			}
			if isMapped {
				info.IsModelMapped = true
				info.UpstreamModelName = mappedModel
				break
			}
		}
	}

	if isResponsesCompact {
		finalUpstreamModelName := baseModelName
		if info.IsModelMapped && info.UpstreamModelName != "" {
			finalUpstreamModelName = info.UpstreamModelName
		}
		finalUpstreamModelName = strings.TrimSuffix(finalUpstreamModelName, ratio_setting.CompactModelSuffix)
		info.UpstreamModelName = finalUpstreamModelName
		info.OriginModelName = ratio_setting.WithCompactModelSuffix(finalUpstreamModelName)
	}
	if request != nil {
		request.SetModelName(info.UpstreamModelName)
	}
	return nil
}

func modelMappingCandidates(originModelName, baseModelName string, isResponsesCompact bool) []string {
	if isResponsesCompact && originModelName != baseModelName {
		return []string{originModelName, baseModelName}
	}
	return []string{baseModelName}
}

func resolveModelMapping(modelMap map[string]string, originModel string) (string, bool, error) {
	currentModel := originModel
	visitedModels := map[string]bool{
		currentModel: true,
	}

	for {
		mappedModel, exists := modelMap[currentModel]
		if !exists || mappedModel == "" {
			return currentModel, currentModel != originModel, nil
		}
		if visitedModels[mappedModel] {
			if mappedModel == currentModel {
				return currentModel, currentModel != originModel, nil
			}
			return "", false, errors.New("model_mapping_contains_cycle")
		}
		visitedModels[mappedModel] = true
		currentModel = mappedModel
	}
}
