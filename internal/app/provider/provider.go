package provider

import (
	"context"

	"userProfile/pkg/userProfile"
)

// ClientInfoProvider specifies mechanism of retrieving ClientInfo.
// which can be either a DB, some microservice, etc.
type ClientInfoProvider interface {
	RegisterClientInfo(context.Context, *userProfile.RegisterRequest) error
	GetClientInfo(context.Context, *userProfile.ClientInfoRequest) (*userProfile.UserProfile, error)
}
