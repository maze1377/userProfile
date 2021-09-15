package core_test

import (
	"context"
	"testing"
	"userProfile/internal/app/provider"
	"userProfile/pkg/cache/adaptors"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"userProfile/internal/app/core"
	"userProfile/pkg/userProfile"

	providerMocks "userProfile/internal/app/provider/mocks"

	"github.com/stretchr/testify/suite"
)

type CoreTestSuite struct {
	suite.Suite
}

func TestCoreTestSuite(t *testing.T) {
	suite.Run(t, new(CoreTestSuite))
}

func (s *CoreTestSuite) TestShouldReturnNotFoundIfProviderReturnsNotFound() {
	cache := adaptors.NewSynMapAdaptor()

	mockProvider := &providerMocks.ClientInfoProvider{}
	mockProvider.On("GetClientInfo", mock.Anything, mock.Anything).Once().Return(nil, provider.ErrNotFound)

	c := core.New(mockProvider, cache)
	_, err := c.GetClientInfo(context.Background(), &userProfile.ClientInfoRequest{})
	s.NotNil(err)

	grpcStatus, ok := status.FromError(err)
	s.True(ok)
	if !ok {
		return
	}
	s.Equal(codes.NotFound, grpcStatus.Code())
}

func (s *CoreTestSuite) TestShouldReturnFromCacheIfFound() {
	cache := adaptors.NewSynMapAdaptor()
	c := core.New(nil, cache)
	clientId := "token"
	err := cache.Set(context.Background(), clientId, &userProfile.UserProfile{
		ClientID: clientId,
	})
	if !s.NoError(err, "fail to marshalize clientInfo") {
		return
	}
	response, err := c.GetClientInfo(context.Background(), &userProfile.ClientInfoRequest{
		ClientID: clientId,
	})

	s.NoError(err)
	s.Equal(response.UserProfile.ClientID, clientId)
}

func (s *CoreTestSuite) TestShouldReturnFromProviderIfNotCached() {
	cache := adaptors.NewSynMapAdaptor()

	mockProvider := &providerMocks.ClientInfoProvider{}
	request := &userProfile.ClientInfoRequest{
		ClientID: "token",
	}
	mockProvider.On("GetClientInfo", mock.Anything, request).Once().Return(&userProfile.UserProfile{
		ClientID: "token",
	}, nil)

	c := core.New(mockProvider, cache)
	response, err := c.GetClientInfo(context.Background(), request)

	s.Nil(err)
	s.Equal(response.UserProfile.ClientID, "token")
}
