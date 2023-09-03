package config_test

import (
	"testing"

	"github.com/ColinSchofield/zai-weather/src/config"

	"github.com/stretchr/testify/assert"
)

func Test_ConfigUsingDefaultValues(t *testing.T) {
	cfg, err := config.LoadConfig()
	assert.NoError(t, err)
	assert.Equal(t, ":8080", cfg.Port)
	assert.Equal(t, 3, cfg.CacheTTLSeconds)
	assert.Equal(t, 3, cfg.PrimaryTimeoutSeconds)
	assert.Equal(t, "1cadfad44c3387c66d14a12cb33f282e", cfg.PrimaryAccessKey)
	assert.Equal(t, "http://api.weatherstack.com/current", cfg.PrimaryEndPoint)
	assert.Equal(t, 3, cfg.FailoverTimeoutSeconds)
	assert.Equal(t, "fe0e197efcdefea9a19e9c4810f2801b", cfg.FailoverAccessKey)
	assert.Equal(t, "http://api.openweathermap.org/data/2.5/weather", cfg.FailoverEndPoint)
	assert.Equal(t, uint32(3), cfg.PrimaryRequests)
	assert.Equal(t, 0.6, cfg.PrimaryFailureRatio)
	assert.Equal(t, uint32(3), cfg.FailoverRequests)
	assert.Equal(t, 0.6, cfg.FailoverFailureRatio)
}

func Test_ConfigFromEnviroment(t *testing.T) {
	t.Setenv("PORT", "1")
	t.Setenv("CACHE_TTL_SECONDS", "2")
	t.Setenv("PRIMARY_TIMEOUT_SECONDS", "3")
	t.Setenv("PRIMARY_ACCESS_KEY", "4")
	t.Setenv("PRIMARY_END_POINT", "5")
	t.Setenv("FAILOVER_TIMEOUT_SECONDS", "6")
	t.Setenv("FAILOVER_ACCESS_KEY", "7")
	t.Setenv("FAILOVER_END_POINT", "8")
	t.Setenv("PRIMARY_REQUESTS", "9")
	t.Setenv("PRIMARY_FAILURE_RATIO", "10")
	t.Setenv("FAILOVER_REQUESTS", "11")
	t.Setenv("FAILOVER_FAILURE_RATIO", "12")

	cfg, err := config.LoadConfig()
	assert.NoError(t, err)
	assert.Equal(t, "1", cfg.Port)
	assert.Equal(t, 2, cfg.CacheTTLSeconds)
	assert.Equal(t, 3, cfg.PrimaryTimeoutSeconds)
	assert.Equal(t, "4", cfg.PrimaryAccessKey)
	assert.Equal(t, "5", cfg.PrimaryEndPoint)
	assert.Equal(t, 6, cfg.FailoverTimeoutSeconds)
	assert.Equal(t, "7", cfg.FailoverAccessKey)
	assert.Equal(t, "8", cfg.FailoverEndPoint)
	assert.Equal(t, uint32(9), cfg.PrimaryRequests)
	assert.Equal(t, float64(10), cfg.PrimaryFailureRatio)
	assert.Equal(t, uint32(11), cfg.FailoverRequests)
	assert.Equal(t, float64(12), cfg.FailoverFailureRatio)
}
