// The service package provides a boundary to the backend, exposed through a set of interfaces.
package service

import (
	"context"
	"fmt"

	"github.com/ColinSchofield/zai-weather/src/config"
	"github.com/ColinSchofield/zai-weather/src/model"

	resty "github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type DefaultOpenWeatherMap struct {
	cfg *config.WeatherConfig
	log *logrus.Logger

	client *resty.Client
}

var _ WeatherFetcher = (*DefaultOpenWeatherMap)(nil)

// NewOpenWeatherMap returns the default struct for the open weather map service.
func NewOpenWeatherMap(cfg *config.WeatherConfig, log *logrus.Logger) *DefaultOpenWeatherMap {
	return &DefaultOpenWeatherMap{
		cfg: cfg,
		log: log,

		client: resty.New(),
	}
}

// The FetchWeather method returns the weather information (in degrees celsius) and wind speed (in km/hr).
func (o *DefaultOpenWeatherMap) FetchWeather(ctx context.Context, location string) (*model.Weather, error) {
	var response model.OpenMapResponse

	queryParams := map[string]string{
		"appid": o.cfg.FailoverAccessKey,
		"q":     location + ",AU", // The country is assumed to be Australia
		"units": "metric",         // Otherwise results will be in Kelvin
	}

	resp, err := o.client.R().
		SetContext(ctx).
		SetQueryParams(queryParams).
		SetResult(&response).
		Get(o.cfg.FailoverEndPoint)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("open weather map returned an unexpected status code of %d", resp.StatusCode())
	}

	return &model.Weather{
		Data: &model.Data{
			Temperature: int(response.Main.Temperature),
			WindSpeed:   int(response.Wind.WindSpeed * 3.6), // need to convert from meters/sec to km/hr
		},
	}, nil
}
