package cache_test

import (
	"testing"
	"time"

	"github.com/ColinSchofield/zai-weather/src/cache"

	"github.com/stretchr/testify/suite"
)

type WeatherCacheTestSuite struct {
	suite.Suite

	weatherCache cache.Weather
}

func TestWeatherCacheSuite(t *testing.T) {
	suite.Run(t, new(WeatherCacheTestSuite))
}

func (s *WeatherCacheTestSuite) SetupTest() {
	s.weatherCache = cache.NewWeatherCache(200 * time.Millisecond)
}

func (s *WeatherCacheTestSuite) Test_HappyPathReadBeforeTTLExpires() {
	// When
	s.weatherCache.Set("one", "1")
	value, ok := s.weatherCache.Get("one")
	// Then
	s.Assert().True(ok)
	s.Assert().Equal("1", value, "Value has not expired")
}

func (s *WeatherCacheTestSuite) Test_HappyPathReadButTTLHasExpired() {
	// When
	s.weatherCache.Set("one", "1")
	time.Sleep(300 * time.Millisecond)
	value, ok := s.weatherCache.Get("one")
	// Then
	s.Assert().False(ok)
	s.Assert().Nil(value)
	// When
	value, ok = s.weatherCache.GetIgnoreTTL("one")
	// Then
	s.Assert().True(ok)
	s.Assert().Equal("1", value, "Value is still in the cache")
}
