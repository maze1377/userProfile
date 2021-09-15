package provider_test

import (
	"context"
	"errors"
	"testing"
	"userProfile/internal/app/provider"
	metricsMocks "userProfile/pkg/metrics/mocks"
	"userProfile/pkg/userProfile"

	providerMocks "userProfile/internal/app/provider/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ProviderInstrumentationMiddlewareTestSuite struct {
	suite.Suite
}

func TestProviderInstrumentationMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(ProviderInstrumentationMiddlewareTestSuite))
}

func (s *ProviderInstrumentationMiddlewareTestSuite) TestGetClientInfoShouldCountErrIsNil() {
	mockedProvider := &providerMocks.ClientInfoProvider{}
	mockedProvider.On("GetClientInfo", mock.Anything, mock.Anything).Return(nil, nil)

	mockedObserver := s.makeObserver(map[string]string{
		"method":  "GetClientInfo",
		"success": "true",
	})

	hooked := provider.NewInstrumentationMiddleware(mockedProvider, mockedObserver)
	_, err := hooked.GetClientInfo(context.Background(), &userProfile.ClientInfoRequest{})
	s.Nil(err)

	mockedObserver.AssertExpectations(s.T())
	mockedProvider.AssertExpectations(s.T())
}

func (s *ProviderInstrumentationMiddlewareTestSuite) TestGetClientInfoShouldCountErrIsNotNil() {
	mockedProvider := &providerMocks.ClientInfoProvider{}
	mockedProvider.On("GetClientInfo", mock.Anything, mock.Anything).Return(nil, errors.New("some err"))

	mockedObserver := s.makeObserver(map[string]string{
		"method":  "GetClientInfo",
		"success": "false",
	})

	hooked := provider.NewInstrumentationMiddleware(mockedProvider, mockedObserver)
	_, err := hooked.GetClientInfo(context.Background(), &userProfile.ClientInfoRequest{})
	s.NotNil(err)

	mockedObserver.AssertExpectations(s.T())
	mockedProvider.AssertExpectations(s.T())
}

func (s *ProviderInstrumentationMiddlewareTestSuite) TestAddClientInfoShouldCountErrIsNil() {
	mockedProvider := &providerMocks.ClientInfoProvider{}
	mockedProvider.On("RegisterClientInfo", mock.Anything, mock.Anything).Return(nil)

	mockedObserver := s.makeObserver(map[string]string{
		"method":  "RegisterClientInfo",
		"success": "true",
	})

	hooked := provider.NewInstrumentationMiddleware(mockedProvider, mockedObserver)
	err := hooked.RegisterClientInfo(context.Background(), nil)
	s.Nil(err)

	mockedObserver.AssertExpectations(s.T())
	mockedProvider.AssertExpectations(s.T())
}

func (s *ProviderInstrumentationMiddlewareTestSuite) TestAddClientInfoShouldCountErrIsNotNil() {
	mockedProvider := &providerMocks.ClientInfoProvider{}
	mockedProvider.On("RegisterClientInfo", mock.Anything, mock.Anything).Return(errors.New("some err"))

	mockedObserver := s.makeObserver(map[string]string{
		"method":  "RegisterClientInfo",
		"success": "false",
	})

	hooked := provider.NewInstrumentationMiddleware(mockedProvider, mockedObserver)
	err := hooked.RegisterClientInfo(context.Background(), nil)
	s.NotNil(err)

	mockedObserver.AssertExpectations(s.T())
	mockedProvider.AssertExpectations(s.T())
}

func (s *ProviderInstrumentationMiddlewareTestSuite) makeObserver(
	expectedLabels map[string]string) *metricsMocks.Observer {

	mockedObserver := &metricsMocks.Observer{}
	mockedObserver.On("Observe", mock.Anything).Once()
	mockedObserver.On("With", s.makeMatcher(expectedLabels)).Once().Return(mockedObserver)

	return mockedObserver
}

func (s *ProviderInstrumentationMiddlewareTestSuite) makeMatcher(
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
