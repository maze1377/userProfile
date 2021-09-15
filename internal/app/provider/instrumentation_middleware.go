package provider

import (
	"context"
	"fmt"
	"time"
	"userProfile/pkg/metrics"

	"userProfile/pkg/userProfile"
)

type instrumentationMiddleware struct {
	next   ClientInfoProvider
	timing metrics.Observer
}

func NewInstrumentationMiddleware(
	next ClientInfoProvider, timing metrics.Observer) ClientInfoProvider {

	return instrumentationMiddleware{
		next:   next,
		timing: timing,
	}
}

func (m instrumentationMiddleware) GetClientInfo(
	ctx context.Context, clientInfo *userProfile.ClientInfoRequest) (result *userProfile.UserProfile, err error) {

	defer func(startTime time.Time) {
		m.timing.With(map[string]string{
			"success": fmt.Sprint(err == nil),
			"method":  "GetClientInfo",
		}).Observe(time.Since(startTime).Seconds())
	}(time.Now())

	return m.next.GetClientInfo(ctx, clientInfo)
}

func (m instrumentationMiddleware) RegisterClientInfo(
	ctx context.Context, clientInfo *userProfile.RegisterRequest) (err error) {

	defer func(startTime time.Time) {
		m.timing.With(map[string]string{
			"success": fmt.Sprint(err == nil),
			"method":  "RegisterClientInfo",
		}).Observe(time.Since(startTime).Seconds())
	}(time.Now())

	return m.next.RegisterClientInfo(ctx, clientInfo)
}
