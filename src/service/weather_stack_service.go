// The service package provides a boundary to the backend, exposed through a set of interfaces.
package service

import (
	"context"
	"errors"

	"github.com/ColinSchofield/zai-weather/src/config"
	"github.com/ColinSchofield/zai-weather/src/model"

	resty "github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

// The WeatherFetcher interface provides an HTTP client for the third party weather stack service.
type WeatherFetcher interface {
	FetchWeather(ctx context.Context, location string) (*model.Weather, error)
}

//go:generate mockgen -source=weather_stack_service.go -destination=../mock/mock_weather_fetcher.go

type DefaultWeatherFetcher struct {
	cfg *config.WeatherConfig
	log *logrus.Logger

	client *resty.Client
}

var _ WeatherFetcher = (*DefaultWeatherFetcher)(nil)

// NewWeatherStack returns the default struct for the weather stack service.
func NewWeatherStack(cfg *config.WeatherConfig, log *logrus.Logger) *DefaultWeatherFetcher {
	return &DefaultWeatherFetcher{
		cfg: cfg,
		log: log,

		client: resty.New(),
	}
}

// The FetchWeather method returns the weather information (in degrees celsius) and wind speed (in Km/hr).
func (s *DefaultWeatherFetcher) FetchWeather(ctx context.Context, location string) (*model.Weather, error) {
	var response model.StackResponse

	queryParams := map[string]string{
		"access_key": s.cfg.PrimaryAccessKey,
		"query":      location,
	}

	resp, err := s.client.R().
		SetContext(ctx).
		SetQueryParams(queryParams).
		SetResult(&response).
		Get(s.cfg.PrimaryEndPoint)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 || response.Error.Code > 0 {
		return nil, errors.New("weather stack did not return any results")
	}

	return &model.Weather{
		Data: &model.Data{
			Temperature: response.Current.Temperature,
			WindSpeed:   response.Current.WindSpeed,
		},
	}, nil
}
