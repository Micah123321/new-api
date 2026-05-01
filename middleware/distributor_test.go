package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/model"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupDistributorTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	gin.SetMode(gin.TestMode)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	model.DB = db
	require.NoError(t, db.AutoMigrate(&model.Channel{}))

	t.Cleanup(func() {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	})

	return db
}

func TestSetupContextForSelectedChannelRefreshesCompactModelMapping(t *testing.T) {
	db := setupDistributorTestDB(t)

	latestMapping := `{"gpt-5.5-openai-compact":"gpt-5.4"}`
	require.NoError(t, db.Create(&model.Channel{
		Id:           38,
		Type:         constant.ChannelTypeCodex,
		Key:          `{"access_token":"token","account_id":"account"}`,
		Status:       common.ChannelStatusEnabled,
		Name:         "codex",
		ModelMapping: &latestMapping,
	}).Error)

	req := httptest.NewRequest(http.MethodPost, "/v1/responses/compact", nil)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = req

	staleMapping := ""
	channel := &model.Channel{
		Id:           38,
		Type:         constant.ChannelTypeCodex,
		Key:          `{"access_token":"cached","account_id":"cached"}`,
		Status:       common.ChannelStatusEnabled,
		Name:         "cached codex",
		ModelMapping: &staleMapping,
	}

	err := SetupContextForSelectedChannel(ctx, channel, "gpt-5.5-openai-compact")
	require.Nil(t, err)
	require.Equal(t, latestMapping, common.GetContextKeyString(ctx, constant.ContextKeyChannelModelMapping))
	require.Equal(t, `{"access_token":"cached","account_id":"cached"}`, common.GetContextKeyString(ctx, constant.ContextKeyChannelKey))
}
