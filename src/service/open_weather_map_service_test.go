package service

import (
	"context"
	"net/http"
	"testing"

	"github.com/ColinSchofield/zai-weather/src/config"
	"github.com/ColinSchofield/zai-weather/src/model"

	"github.com/jarcoal/httpmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type OpenWeatherMapServiceTestSuite struct {
	suite.Suite

	ctx          context.Context
	clientSvc    *DefaultOpenWeatherMap
	mockResponse model.OpenMapResponse
}

func TestOpenWeatherMapServiceSuite(t *testing.T) {
	suite.Run(t, new(OpenWeatherMapServiceTestSuite))
}

func (s *OpenWeatherMapServiceTestSuite) SetupTest() {
	s.ctx = context.Background()
	log := logrus.New()
	cfg := &config.WeatherConfig{
		FailoverEndPoint: "http://localhost",
	}
	httpmock.Activate()
	s.mockResponse = model.OpenMapResponse{
		Main: model.Main{
			Temperature: 5,
		},
		Wind: model.Wind{
			WindSpeed: 10,
		},
	}
	s.clientSvc = NewOpenWeatherMap(cfg, log)
	httpmock.ActivateNonDefault(s.clientSvc.client.GetClient())
}

func (suite *OpenWeatherMapServiceTestSuite) TearDownTest() {
	httpmock.DeactivateAndReset()
}

func (s *OpenWeatherMapServiceTestSuite) Test_OpenWeatherMapServiceSuccessful() {
	// When
	httpmock.RegisterResponder("GET", "http://localhost", httpmock.NewJsonResponderOrPanic(http.StatusOK, s.mockResponse))
	res, err := s.clientSvc.FetchWeather(s.ctx, "Melbourne")
	// Then
	s.Suite.Assert().NoError(err)
	s.Suite.Assert().True(res.Data.Temperature == 5 && res.Data.WindSpeed == 36, "wind speed is converted to km/hr")
}

func (s *OpenWeatherMapServiceTestSuite) Test_OpenWeatherMapServiceUnsuccessful() {
	// When
	httpmock.RegisterResponder("GET", "http://localhost", httpmock.NewJsonResponderOrPanic(http.StatusNotFound, nil))
	res, err := s.clientSvc.FetchWeather(s.ctx, "Melbourne")
	// Then
	s.Suite.Assert().Error(err)
	s.Suite.Assert().Nil(res)
}
