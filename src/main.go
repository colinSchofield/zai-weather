package main

import (
	"github.com/ColinSchofield/zai-weather/src/config"
	"github.com/ColinSchofield/zai-weather/src/controller"
	"github.com/ColinSchofield/zai-weather/src/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/sony/gobreaker"
)

// Run a microservice to serve requests for temperature (in celsius) and wind speed (in km/hr).
// Code is separated into packages (i.e. controller, service, model etc) based upon the separation of concerns.
// Software cache the results, based upon a configured TTL.
// Use a primary and a fail-over 3rd party weather provider.
// Handling of the primary and fail-over 3rd party servers, is done by using the circuit breaker design pattern.
//
// See https://en.wikipedia.org/wiki/Circuit_breaker_design_pattern.
func main() {
	log := logrus.New()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.WithError(err).Fatal("failed to load the configuration")
	}

	weatherController := controller.NewWeatherController(
		cfg,
		log,
		service.NewWeatherStack(cfg, log),
		service.NewOpenWeatherMap(cfg, log),
		gobreaker.NewCircuitBreaker(
			gobreaker.Settings{
				Name: "Weather Stack (primary)",
				ReadyToTrip: func(counts gobreaker.Counts) bool {
					failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
					return counts.Requests >= cfg.PrimaryRequests && failureRatio >= cfg.PrimaryFailureRatio
				},
			},
		),
		gobreaker.NewCircuitBreaker(
			gobreaker.Settings{
				Name: "Open Weather Map (failover)",
				ReadyToTrip: func(counts gobreaker.Counts) bool {
					failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
					return counts.Requests >= cfg.FailoverRequests && failureRatio >= cfg.FailoverFailureRatio
				},
			},
		),
	)

	log.Info("Starting Zai Weather REST API Service on Port ", cfg.Port)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("v1/weather", weatherController.GetWeather)
	if err := router.Run(cfg.Port); err != nil {
		log.WithError(err).WithField("port_num", cfg.Port).Fatal("failed to run HTTP service")
	}
}
