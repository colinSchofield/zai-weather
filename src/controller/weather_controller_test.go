package controller_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ColinSchofield/zai-weather/src/config"
	"github.com/ColinSchofield/zai-weather/src/controller"
	mock "github.com/ColinSchofield/zai-weather/src/mock"
	"github.com/ColinSchofield/zai-weather/src/model"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/suite"
)

type ControllerTestSuite struct {
	suite.Suite

	ctrl         *gomock.Controller
	ctx          context.Context
	log          *logrus.Logger
	cfg          *config.WeatherConfig
	cbP          *gobreaker.CircuitBreaker
	cbF          *gobreaker.CircuitBreaker
	mockPrimary  *mock.MockWeatherFetcher
	mockFailover *mock.MockWeatherFetcher
	controller   controller.WeatherController
	record       *httptest.ResponseRecorder
	gCtx         *gin.Context
}

func TestControllerSuite(t *testing.T) {
	suite.Run(t, new(ControllerTestSuite))
}

func (s *ControllerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.ctx = context.Background()
	s.log = logrus.New()
	s.cfg = &config.WeatherConfig{
		CacheTTLSeconds: 1,
	}
	s.cbP = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.TotalFailures >= 1
		},
	})
	s.cbF = gobreaker.NewCircuitBreaker(gobreaker.Settings{})
	s.mockPrimary = mock.NewMockWeatherFetcher(s.ctrl)
	s.mockFailover = mock.NewMockWeatherFetcher(s.ctrl)
	s.controller = controller.NewWeatherController(
		s.cfg,
		s.log,
		s.mockPrimary,
		s.mockFailover,
		s.cbP,
		s.cbF,
	)
	s.record = httptest.NewRecorder()
	s.gCtx, _ = gin.CreateTestContext(s.record)
}

func (s *ControllerTestSuite) Test_HappyPath() {
	// Given
	mockResponse := &model.Weather{
		Data: &model.Data{
			Temperature: 10,
			WindSpeed:   15,
		},
	}
	s.mockPrimary.EXPECT().FetchWeather(gomock.Any(), "Melbourne").Return(mockResponse, nil)
	// When
	s.controller.GetWeather(s.gCtx)
	// Then
	s.Assert().Equal(http.StatusOK, s.record.Code, "HTTP status of 200")
	s.Assert().NotNil(s.record.Body)
}

func (s *ControllerTestSuite) Test_HappyPathThenReadsFromTheCache() {
	// Given
	mockResponse := &model.Weather{
		Data: &model.Data{
			Temperature: 10,
			WindSpeed:   15,
		},
	}
	s.mockPrimary.EXPECT().FetchWeather(gomock.Any(), "Melbourne").Return(mockResponse, nil)
	// When
	s.controller.GetWeather(s.gCtx)
	// Then
	s.Assert().Equal(http.StatusOK, s.record.Code, "HTTP status of 200")
	s.Assert().NotNil(s.record.Body)
	// When
	s.controller.GetWeather(s.gCtx)
	// Then
	s.Assert().Equal(http.StatusOK, s.record.Code, "HTTP status of 200")
	s.Assert().NotNil(s.record.Body)
}

func (s *ControllerTestSuite) Test_FailedPrimaryFailoverSuccess() {
	// Given
	mockResponse := &model.Weather{
		Data: &model.Data{
			Temperature: 10,
			WindSpeed:   15,
		},
	}
	s.mockPrimary.EXPECT().FetchWeather(gomock.Any(), "Melbourne").Return(nil, errors.New("Server is down!"))
	s.mockFailover.EXPECT().FetchWeather(gomock.Any(), "Melbourne").Return(mockResponse, nil)
	// When
	s.controller.GetWeather(s.gCtx)
	// Then
	s.Assert().Equal(http.StatusOK, s.record.Code, "HTTP status of 200")
	s.Assert().NotNil(s.record.Body)
}

func (s *ControllerTestSuite) Test_BothPrimaryAndFailoverFail() {
	// Given
	mockResponse := &model.Weather{
		Data: &model.Data{
			Temperature: 10,
			WindSpeed:   15,
		},
	}
	s.mockPrimary.EXPECT().FetchWeather(gomock.Any(), "Melbourne").Return(mockResponse, nil)
	s.mockPrimary.EXPECT().FetchWeather(gomock.Any(), "Melbourne").Return(nil, errors.New("Server is down!"))
	s.mockFailover.EXPECT().FetchWeather(gomock.Any(), "Melbourne").Return(nil, errors.New("Server is down!"))
	// When
	s.controller.GetWeather(s.gCtx)
	time.Sleep(1100 * time.Millisecond)
	s.controller.GetWeather(s.gCtx)
	// Then
	s.Assert().Equal(http.StatusOK, s.record.Code, "HTTP status of 200 (falback to the cache!)")
}

func (s *ControllerTestSuite) Test_BothPrimaryAndFailoverFailSoFallbackToCache() {
	// Given
	s.mockPrimary.EXPECT().FetchWeather(gomock.Any(), "Melbourne").Return(nil, errors.New("Server is down!"))
	s.mockFailover.EXPECT().FetchWeather(gomock.Any(), "Melbourne").Return(nil, errors.New("Server is down!"))
	// When
	s.controller.GetWeather(s.gCtx)
	// Then
	s.Assert().Equal(http.StatusNotFound, s.record.Code, "HTTP status of 404 (as nothing is in the cache!)")
}

func (s *ControllerTestSuite) Test_PrimaryFailsAndCircuitBreakerIsTripped() {
	// Given
	mockResponse := &model.Weather{
		Data: &model.Data{
			Temperature: 10,
			WindSpeed:   15,
		},
	}
	s.mockPrimary.EXPECT().FetchWeather(gomock.Any(), "Melbourne").Return(nil, errors.New("Server is down!"))
	s.mockFailover.EXPECT().FetchWeather(gomock.Any(), "Melbourne").Times(2).Return(mockResponse, nil)
	// When
	s.controller.GetWeather(s.gCtx)
	// Then
	s.Assert().Equal(http.StatusOK, s.record.Code, "HTTP status of 200")
	s.Assert().NotNil(s.record.Body)
	// When
	time.Sleep(1100 * time.Millisecond)
	s.controller.GetWeather(s.gCtx)
	// Then
	s.Assert().Equal(http.StatusOK, s.record.Code, "HTTP status of 200")
	s.Assert().NotNil(s.record.Body)
}
