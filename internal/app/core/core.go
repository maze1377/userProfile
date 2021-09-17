package core

import (
	"context"
	"strings"
	"userProfile/pkg/errors"

	"golang.org/x/xerrors"
	"google.golang.org/protobuf/types/known/timestamppb"

	"userProfile/internal/app/provider"
	"userProfile/pkg/cache"
	"userProfile/pkg/userProfile"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type core struct {
	provider provider.ClientInfoProvider
	cache    cache.Layer
}

func New(provider provider.ClientInfoProvider, cache cache.Layer) userProfile.UserProfileServer {
	return &core{
		provider: provider,
		cache:    cache,
	}
}

func (c *core) GetClientInfo(
	ctx context.Context, clientInfo *userProfile.ClientInfoRequest) (result *userProfile.ClientInfoResponse, errResult error) {
	defer func() {
		if errResult == nil {
			contains := clientInfo.GetContains()
			if !contains.GetClientInfo() {
				result.UserProfile.ClientInfo = nil
			}
			if !contains.GetAndroidInfo() {
				result.UserProfile.AndroidInfo = nil
			}
			if !contains.GetFeature() {
				result.UserProfile.Features = nil
			}
			if !contains.GetLibrary() {
				result.UserProfile.Libraries = nil
			}
		}
	}()
	profile, err := c.getClientInfoFromCache(ctx, clientInfo.GetClientID())
	if err != nil {
		logrus.WithError(err).WithFields(map[string]interface{}{
			"clientId": clientInfo.GetClientID(),
		}).Error("failed to load data from cache")
	} else {
		return &userProfile.ClientInfoResponse{
			ResponseTimestamp: timestamppb.Now(),
			UserProfile:       profile,
		}, nil
	}

	profile, err = c.provider.GetClientInfo(ctx, clientInfo)
	if err != nil {
		if xerrors.Is(err, provider.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "clientId not found")
		}

		return nil, errors.WrapWithExtra(err, "failed to acquire client profile", map[string]interface{}{
			"request": clientInfo,
		})
	}

	err = c.setUserInfoFromCache(ctx, clientInfo.GetClientID(), profile)
	if err != nil {
		logrus.WithError(err).WithFields(map[string]interface{}{
			"clientId": clientInfo.GetClientID(),
		}).Error("failed to set data in cache")
	}
	return &userProfile.ClientInfoResponse{
		UserProfile: profile,
	}, nil
}

func (c *core) RegisterClientInfo(ctx context.Context, clientInfo *userProfile.RegisterRequest) (*userProfile.RegisterResponse, error) {
	// todo validate clientInfo
	err := c.provider.RegisterClientInfo(ctx, clientInfo)
	if err != nil {
		logrus.WithError(err).WithFields(map[string]interface{}{
			"clientInfo": clientInfo,
		}).Error("failed to set data in provider")
		return nil, err
	}

	err = c.setUserInfoFromCache(ctx, clientInfo.GetUserProfile().GetClientID(), clientInfo.GetUserProfile())
	if err != nil {
		logrus.WithError(err).WithFields(map[string]interface{}{
			"clientId": clientInfo.GetUserProfile().GetClientID(),
		}).Error("failed to set data in cache")
	}
	return &userProfile.RegisterResponse{
		ResponseTimestamp: timestamppb.Now(),
	}, nil
}

func (c *core) getClientInfoFromCache(ctx context.Context, clientId string) (*userProfile.UserProfile, error) {
	key := strings.TrimSpace(clientId)
	if len(key) == 0 {
		return nil, errors.New("clientId not valid")
	}
	var result userProfile.UserProfile
	err := c.cache.Get(ctx, key, &result)
	if err != nil {
		return nil, errors.Wrap(err, "fail to get UserProfile from cache")
	}
	logrus.Info("load UserProfile from cache")
	return &result, nil
}

func (c *core) setUserInfoFromCache(ctx context.Context, clientId string, clientInfo *userProfile.UserProfile) (err error) {
	key := strings.TrimSpace(clientId)
	if len(key) == 0 {
		return errors.New("clientId not valid")
	}
	err = c.cache.Set(ctx, key, clientInfo)
	if err != nil {
		return errors.Wrap(err, "fail to set clientInfo in cache")
	}
	return nil
}
