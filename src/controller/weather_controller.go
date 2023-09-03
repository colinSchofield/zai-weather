// controller package has direct access to the web/http layer. Its purpose is to mediate access to the service layer.
package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/ColinSchofield/zai-weather/src/cache"
	"github.com/ColinSchofield/zai-weather/src/config"
	"github.com/ColinSchofield/zai-weather/src/model"
	"github.com/ColinSchofield/zai-weather/src/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/sony/gobreaker"
)

const (
	// These message are returne in the JSON message field
	MessageSuccess      = "Request successful"
	MessageSuccessCache = "Request successful (cached)"
	MessageFailureCache = "Request failure (cache is stale)"
	MessageFailure      = "Location could not be found"
)

// The WeatherController interface provides access to the current weather conditions.
type WeatherController interface {
	GetWeather(gCtx *gin.Context)
}

type DefaultWeatherController struct {
	cfg        *config.WeatherConfig
	log        *logrus.Logger
	primary    service.WeatherFetcher
	failover   service.WeatherFetcher
	cbPrimary  *gobreaker.CircuitBreaker
	cbFailover *gobreaker.CircuitBreaker

	weatherCache *cache.DefaultWeatherCache
}

var _ WeatherController = (*DefaultWeatherController)(nil)

// NewWeatherController returns the default struct for the weather controller.
func NewWeatherController(
	cfg *config.WeatherConfig,
	log *logrus.Logger,
	primary service.WeatherFetcher,
	failover service.WeatherFetcher,
	cbPrimary *gobreaker.CircuitBreaker,
	cbFailover *gobreaker.CircuitBreaker,
) *DefaultWeatherController {
	return &DefaultWeatherController{
		cfg:        cfg,
		log:        log,
		primary:    primary,
		failover:   failover,
		cbPrimary:  cbPrimary,
		cbFailover: cbFailover,

		weatherCache: cache.NewWeatherCache(time.Duration(cfg.CacheTTLSeconds) * time.Second),
	}
}

// GetWeather returns a JSON value containing the temperature (in degrees celsius) and the wind speed (in km/hr).
// TODO Add Prometheus metrics for monitoring and observability.
//
// The handling of the primary and fail-over 3rd party servers, is done using the circuit breaker design pattern.
//
// See https://en.wikipedia.org/wiki/Circuit_breaker_design_pattern.
func (w *DefaultWeatherController) GetWeather(gCtx *gin.Context) {
	location := gCtx.DefaultQuery("city", "Melbourne")

	// Load the weather information, if possible, from the cache.
	if cached, found := w.weatherCache.Get(location); found {
		weather := cached.(*model.Weather)
		weather.Message = MessageSuccessCache
		gCtx.JSON(http.StatusOK, weather)
		return
	}

	// Fetch from primary service.
	if ok := w.fetchWeather(gCtx, w.cfg.PrimaryTimeoutSeconds, w.cbPrimary, location, w.primary); ok {
		return
	}

	// Fetch from fail-over service.
	if ok := w.fetchWeather(gCtx, w.cfg.FailoverTimeoutSeconds, w.cbFailover, location, w.failover); ok {
		return
	}

	// Fallback to cached values.
	if fallback, found := w.weatherCache.GetIgnoreTTL(location); found {
		weather := fallback.(*model.Weather)
		weather.Message = MessageFailureCache
		gCtx.JSON(http.StatusOK, fallback)
		return
	}

	// Assumee that the location is invalid.
	weather := model.Weather{
		Status:  http.StatusNotFound,
		Message: MessageFailure,
	}
	gCtx.JSON(http.StatusNotFound, weather)
}

// Fetch the weather information from a weather service.
func (w *DefaultWeatherController) fetchWeather(
	gCtx *gin.Context,
	timeout int,
	cb *gobreaker.CircuitBreaker,
	location string,
	fetcher service.WeatherFetcher,
) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	if res, err := cb.Execute(func() (interface{}, error) {
		return fetcher.FetchWeather(ctx, location)
	}); err != nil {
		w.log.WithError(err).WithField("location", location).Warn("Failed to fetch from ", cb.Name())
		// TODO Add in some metrics here (perhaps based upon cb.Counts())
		return false
	} else {
		response := res.(*model.Weather)
		response.Status = http.StatusOK
		response.Message = MessageSuccess
		w.weatherCache.Set(location, response)
		gCtx.JSON(http.StatusOK, res)
		return true
	}
}
