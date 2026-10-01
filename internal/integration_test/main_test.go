package integration_test

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/tranvantuan81/bookmark-management/internal/api"
	"github.com/tranvantuan81/bookmark-management/internal/config"
	redisPkg "github.com/tranvantuan81/bookmark-management/pkg/redis"
)

const (
	testServiceName = "bookmark-management"
	testInstanceID  = "0f1092e7-38ed-4701-b500-c4697c9dc122"

	// expectedPasswordLength mirrors handler.passwordLength
	expectedPasswordLength = 12
	// expectedCodeLength mirrors service.DefaultCodeLength
	expectedCodeLength = 7
)

// TestMain silences gin's debug output so failures stay readable.
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

// setupTestAPI wires the real engine to an in-memory Redis, so routes, handler,
// service and repository run together exactly like in production.
func setupTestAPI(t *testing.T) (api.Engine, *redis.Client) {
	t.Helper()

	redisClient := redisPkg.InitMockRedis(t)

	engine := api.NewEngine(&config.Config{
		AppPort:     "8080",
		ServiceName: testServiceName,
		InstanceID:  testInstanceID,
	}, redisClient)

	return engine, redisClient
}

// decodeBody unmarshals the JSON response body into dst. Use it when a field value
// is not deterministic (generated code, generated password, ttl).
func decodeBody(t *testing.T, rec *httptest.ResponseRecorder, dst any) {
	t.Helper()

	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), dst))
}
