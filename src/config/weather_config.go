package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type WeatherConfig struct {
	// See the Dockerfile for the Port mappings
	Port            string `env:"PORT" env-default:":8080"`
	CacheTTLSeconds int    `env:"CACHE_TTL_SECONDS" env-default:"3"`
	// The primary is the Weather Stack Service
	PrimaryTimeoutSeconds int    `env:"PRIMARY_TIMEOUT_SECONDS" env-default:"3"`
	PrimaryAccessKey      string `env:"PRIMARY_ACCESS_KEY" env-default:"1cadfad44c3387c66d14a12cb33f282e"`
	PrimaryEndPoint       string `env:"PRIMARY_END_POINT" env-default:"http://api.weatherstack.com/current"`
	// The failover is the Open Weather Map Service
	FailoverTimeoutSeconds int    `env:"FAILOVER_TIMEOUT_SECONDS" env-default:"3"`
	FailoverAccessKey      string `env:"FAILOVER_ACCESS_KEY" env-default:"fe0e197efcdefea9a19e9c4810f2801b"`
	FailoverEndPoint       string `env:"FAILOVER_END_POINT" env-default:"http://api.openweathermap.org/data/2.5/weather"`
	// Circuit Breaker used for both the primary and failover
	PrimaryRequests      uint32  `env:"PRIMARY_REQUESTS" env-default:"3"`
	PrimaryFailureRatio  float64 `env:"PRIMARY_FAILURE_RATIO" env-default:"0.6"`
	FailoverRequests     uint32  `env:"FAILOVER_REQUESTS" env-default:"3"`
	FailoverFailureRatio float64 `env:"FAILOVER_FAILURE_RATIO" env-default:"0.6"`
}

// LoadConfig reads the configuration from the system environment variables.
func LoadConfig() (*WeatherConfig, error) {
	var cfg WeatherConfig
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
