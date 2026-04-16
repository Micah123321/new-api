package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/setting"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type stripeAPIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func setupStripeControllerTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	gin.SetMode(gin.TestMode)
	common.UsingSQLite = true
	common.UsingMySQL = false
	common.UsingPostgreSQL = false
	common.RedisEnabled = false

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}
	model.DB = db
	model.LOG_DB = db

	if err := db.AutoMigrate(&model.SubscriptionPlan{}); err != nil {
		t.Fatalf("failed to migrate subscription plan table: %v", err)
	}

	t.Cleanup(func() {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	})

	return db
}

func newStripeTestContext(t *testing.T, method string, target string, body any) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()

	var requestBody *bytes.Reader
	if body != nil {
		payload, err := common.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}
		requestBody = bytes.NewReader(payload)
	} else {
		requestBody = bytes.NewReader(nil)
	}

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(method, target, requestBody)
	if body != nil {
		ctx.Request.Header.Set("Content-Type", "application/json")
	}
	ctx.Set("id", 1)
	return ctx, recorder
}

func decodeStripeAPIResponse(t *testing.T, recorder *httptest.ResponseRecorder) stripeAPIResponse {
	t.Helper()

	var response stripeAPIResponse
	if err := common.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode api response: %v", err)
	}
	return response
}

func setStripePaymentConfigForTest(t *testing.T, apiSecret string, webhookSecret string) {
	t.Helper()

	originalAPISecret := setting.StripeApiSecret
	originalWebhookSecret := setting.StripeWebhookSecret
	t.Cleanup(func() {
		setting.StripeApiSecret = originalAPISecret
		setting.StripeWebhookSecret = originalWebhookSecret
	})

	setting.StripeApiSecret = apiSecret
	setting.StripeWebhookSecret = webhookSecret
}

func TestStripeWebhookRejectsEmptyWebhookSecret(t *testing.T) {
	setStripePaymentConfigForTest(t, "sk_test_valid", "")

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/stripe/webhook", strings.NewReader(`{"type":"checkout.session.completed"}`))
	ctx.Request.Header.Set("Stripe-Signature", "t=1,v1=fake")

	StripeWebhook(ctx)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestRequestStripePayRejectsMissingWebhookSecret(t *testing.T) {
	setStripePaymentConfigForTest(t, "sk_test_valid", " ")

	ctx, recorder := newStripeTestContext(t, http.MethodPost, "/api/user/stripe/pay", StripePayRequest{
		Amount:        getStripeMinTopup(),
		PaymentMethod: PaymentMethodStripe,
	})

	RequestStripePay(ctx)

	response := decodeStripeAPIResponse(t, recorder)
	if response.Message != "error" {
		t.Fatalf("expected error response, got %+v", response)
	}
	if response.Data != "Stripe Webhook 未配置" {
		t.Fatalf("expected webhook config error, got %+v", response)
	}
}

func TestSubscriptionRequestStripePayRejectsMissingWebhookSecret(t *testing.T) {
	db := setupStripeControllerTestDB(t)
	setStripePaymentConfigForTest(t, "sk_test_valid", "")

	plan := &model.SubscriptionPlan{
		Title:         "Starter",
		PriceAmount:   9.9,
		Currency:      "USD",
		DurationUnit:  model.SubscriptionDurationMonth,
		DurationValue: 1,
		Enabled:       true,
		StripePriceId: "price_test_123",
	}
	if err := db.Create(plan).Error; err != nil {
		t.Fatalf("failed to create subscription plan: %v", err)
	}

	ctx, recorder := newStripeTestContext(t, http.MethodPost, "/api/subscription/stripe/pay", SubscriptionStripePayRequest{
		PlanId: plan.Id,
	})

	SubscriptionRequestStripePay(ctx)

	response := decodeStripeAPIResponse(t, recorder)
	if response.Success {
		t.Fatalf("expected error response, got %+v", response)
	}
	if response.Message != "Stripe Webhook 未配置" {
		t.Fatalf("expected webhook config error, got %+v", response)
	}
}
