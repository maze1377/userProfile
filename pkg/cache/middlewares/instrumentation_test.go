package middlewares

import (
	"context"
	"testing"
	"userProfile/pkg/cache/adaptors"
	metricsMocks "userProfile/pkg/metrics/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CacheInstrumentationMiddlewareTestSuite struct {
	suite.Suite
}

func TestCacheInstrumentationMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(CacheInstrumentationMiddlewareTestSuite))
}

func (s *CacheInstrumentationMiddlewareTestSuite) TestGetShouldCountOkIsFalse() {
	cache := adaptors.NewSynMapAdaptor()

	mockedObserver := s.makeObserver(map[string]string{
		"method": "Get",
	})

	hooked := NewInstrumentationMiddleware(cache, mockedObserver)
	var value interface{}
	err := hooked.Get(context.Background(), "", &value)
	s.NotNil(err)

	mockedObserver.AssertExpectations(s.T())
}

func (s *CacheInstrumentationMiddlewareTestSuite) TestGetShouldCountOkIsTrue() {
	cache := adaptors.NewSynMapAdaptor()
	err := cache.Set(context.Background(), "", "value")
	s.NoError(err)

	mockedObserver := s.makeObserver(map[string]string{
		"method": "Get",
	})

	hooked := NewInstrumentationMiddleware(cache, mockedObserver)
	var value string
	err = hooked.Get(context.Background(), "", &value)
	s.Nil(err)

	mockedObserver.AssertExpectations(s.T())
}

func (s *CacheInstrumentationMiddlewareTestSuite) makeObserver(
	expectedLabels map[string]string) *metricsMocks.Observer {

	mockedObserver := &metricsMocks.Observer{}
	mockedObserver.On("Observe", mock.Anything).Once()
	mockedObserver.On("With", s.makeMatcher(expectedLabels)).Once().Return(mockedObserver)

	return mockedObserver
}

func (s *CacheInstrumentationMiddlewareTestSuite) makeMatcher(
	expectedLabels map[string]string) interface{} {

	return mock.MatchedBy(func(labels map[string]string) bool {
		result := true

		for expectedKey, expectedValue := range expectedLabels {
			value, ok := labels[expectedKey]
			s.True(ok, "expected to find label %v", expectedKey)
			if !ok {
				result = false
				continue
			}

			s.Equal(expectedValue, value,
				"expected to find value %v for key %v", expectedValue, expectedKey)
			result = result && expectedValue == value
		}

		return result
	})
}
