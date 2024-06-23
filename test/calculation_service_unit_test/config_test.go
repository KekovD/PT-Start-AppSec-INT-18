package calculation_service_unit_test

import (
	"calculation_service/controller"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseEnvVariables(t *testing.T) {
	_ = os.Setenv("REDIS_HOST", "localhost")
	_ = os.Setenv("REDIS_PORT", "6379")
	_ = os.Setenv("REDIS_DATABASE", "0")
	_ = os.Setenv("REDIS_TTL", "10s")
	_ = os.Setenv("REDIS_SYNC_INTERVAL", "1m")
	_ = os.Setenv("INTERVAL", "5m")
	_ = os.Setenv("LIMIT", "100")
	_ = os.Setenv("SYNC_WINDOW_KEY", "test_key")

	defer func() {
		_ = os.Unsetenv("REDIS_HOST")
		_ = os.Unsetenv("REDIS_PORT")
		_ = os.Unsetenv("REDIS_DATABASE")
		_ = os.Unsetenv("REDIS_TTL")
		_ = os.Unsetenv("REDIS_SYNC_INTERVAL")
		_ = os.Unsetenv("INTERVAL")
		_ = os.Unsetenv("LIMIT")
		_ = os.Unsetenv("SYNC_WINDOW_KEY")
	}()

	redisHost, redisPort, redisDatabase, redisTtl, syncInterval, interval, limit, syncWindowKey := controller.ParseEnvVariables()

	assert.Equal(t, "localhost", redisHost)
	assert.Equal(t, 6379, redisPort)
	assert.Equal(t, 0, redisDatabase)
	assert.Equal(t, 10*time.Second, redisTtl)
	assert.Equal(t, 1*time.Minute, syncInterval)
	assert.Equal(t, 5*time.Minute, interval)
	assert.Equal(t, int64(100), limit)
	assert.Equal(t, "test_key", syncWindowKey)
}
