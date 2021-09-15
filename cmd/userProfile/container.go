//+build wireinject

package main

import (
	"context"
	"userProfile/internal/app/core"
	"userProfile/internal/app/provider"
	"userProfile/internal/pkg/grpcserver"

	"github.com/google/wire"
	"github.com/spf13/cobra"
)

func CreateServer(ctx context.Context, cmd *cobra.Command) (*grpcserver.Server, error) {
	panic(wire.Build(
		provideConfig,
		provideLogger,
		provideProvider,
		provideCache,
		providePrometheus,
		provideServer,
		core.New,
	))
}

func CreateProvider(ctx context.Context, cmd *cobra.Command) (provider.ClientInfoProvider, error) {
	panic(wire.Build(
		provideConfig,
		provideLogger,
		provideProvider,
		providePrometheus,
	))
}
