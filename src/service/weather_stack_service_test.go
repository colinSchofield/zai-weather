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

type WeatherStackServiceTestSuite struct {
	suite.Suite

	ctx             context.Context
	clientSvc       *DefaultWeatherFetcher
	mockResponse    model.StackResponse
	mockBadResponse model.StackResponse
}

func TestWeatherStackServiceSuite(t *testing.T) {
	suite.Run(t, new(WeatherStackServiceTestSuite))
}

func (s *WeatherStackServiceTestSuite) SetupTest() {
	s.ctx = context.Background()
	log := logrus.New()
	cfg := &config.WeatherConfig{
		PrimaryEndPoint: "http://localhost",
	}
	httpmock.Activate()
	s.mockResponse = model.StackResponse{
		Error: model.Error{},
		Current: model.Current{
			Temperature: 5,
			WindSpeed:   36,
		},
	}
	s.mockBadResponse = model.StackResponse{
		Error: model.Error{
			Code: 42, // This service always returns with http.StatusOK, but with an error code above zero
		},
		Current: model.Current{},
	}
	s.clientSvc = NewWeatherStack(cfg, log)
	httpmock.ActivateNonDefault(s.clientSvc.client.GetClient())
}

func (suite *WeatherStackServiceTestSuite) TearDownTest() {
	httpmock.DeactivateAndReset()
}

func (s *WeatherStackServiceTestSuite) Test_WeatherStackServiceSuccessful() {
	// When
	httpmock.RegisterResponder("GET", "http://localhost", httpmock.NewJsonResponderOrPanic(http.StatusOK, s.mockResponse))
	res, err := s.clientSvc.FetchWeather(s.ctx, "Melbourne")
	// Then
	s.Suite.Assert().NoError(err)
	s.Suite.Assert().True(res.Data.Temperature == 5 && res.Data.WindSpeed == 36, "all values are in the correct units")
}

func (s *WeatherStackServiceTestSuite) Test_WeatherStackServiceUnsuccessful() {
	// When
	httpmock.RegisterResponder("GET", "http://localhost", httpmock.NewJsonResponderOrPanic(http.StatusOK, s.mockBadResponse))
	res, err := s.clientSvc.FetchWeather(s.ctx, "Melbourne")
	// Then
	s.Suite.Assert().Error(err)
	s.Suite.Assert().Nil(res)
}
