package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"userProfile/internal/app/provider"
	"userProfile/internal/pkg/config"
	"userProfile/internal/pkg/grpcserver"
	"userProfile/internal/pkg/prometheus_metirc"
	"userProfile/pkg/cache"
	"userProfile/pkg/cache/adaptors"
	"userProfile/pkg/cache/middlewares"
	"userProfile/pkg/cache/multilayercache"
	"userProfile/pkg/errors"
	"userProfile/pkg/metrics/prometheus"
	"userProfile/pkg/userProfile"

	"github.com/allegro/bigcache/v2"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"userProfile/pkg/sql"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start Server",
	Run:   serve,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serve(cmd *cobra.Command, args []string) {
	printVersion()

	serverCtx, serverCancel := makeServerCtx()
	defer serverCancel()

	server, err := CreateServer(serverCtx, cmd)
	if err != nil {
		panicWithError(err, "failed to create server")
	}

	var serverWaitGroup sync.WaitGroup

	serverWaitGroup.Add(1)
	go func() {
		defer serverWaitGroup.Done()

		if err := server.Serve(); err != nil {
			panicWithError(err, "failed to serve")
		}
	}()

	liveness := Liveness{}
	liveness.startServer()

	<-serverCtx.Done()

	liveness.stopServer()
	server.Stop()

	serverWaitGroup.Wait()
}

func makeServerCtx() (context.Context, context.CancelFunc) {
	gracefulStop := make(chan os.Signal)

	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-gracefulStop
		cancel()
	}()

	return ctx, cancel
}

func provideServer(server userProfile.UserProfileServer, config *config.Config, logger *logrus.Logger) (*grpcserver.Server, error) {
	return grpcserver.New(server, logger, config.ListenPort)
}

func provideProvider(config *config.Config, _ *logrus.Logger, _ *prometheus.Server) provider.ClientInfoProvider {
	dbWrite, err := sql.GetDatabase(config.Database.WriteClient)
	if err != nil {
		logrus.WithError(err).WithField(
			"database write", config.Database.WriteClient).Panic("failed to connect to DB")
		return nil
	}
	dbRead, err := sql.GetDatabase(config.Database.ReadClients)
	if err != nil {
		logrus.WithError(err).WithField(
			"database read", config.Database.WriteClient).Panic("failed to connect to DB")
		return nil
	}
	providerInstance := provider.NewSQLWithReadAndWrite(dbWrite, dbRead)
	err = providerInstance.(sql.Migrate).Migrate()
	if err != nil {
		logrus.WithError(err).WithField(
			"database", config.Database.WriteClient).Panic("failed to Migrate")
		return nil
	}
	providerInstance = provider.NewInstrumentationMiddleware(
		providerInstance, prometheus_metirc.UserProviderMetrics.With(map[string]string{
			"provider_type": "postgres",
		}))

	return providerInstance
}

func provideCache(config *config.Config) (cache.Layer, error) {
	var cacheLayers []cache.Layer
	if config.Cache.Redis.Enabled {
		redisConfig := adaptors.InstanceOptions{
			Address: adaptors.InstanceOptionsAddress{
				Master:   config.Cache.Redis.HostMaster,
				Replicas: config.Cache.Redis.HostReadOnly,
			},
			Password: config.Cache.Redis.Password,
			DBNumber: config.Cache.Redis.DB,
			Expire:   config.Cache.Redis.ExpirationTime,
		}
		redisCache, err := adaptors.NewRedisAdaptor(
			&redisConfig)
		if err != nil {
			return nil, errors.Wrap(err, "fail to redis cache")
		}
		cacheLayers = append(cacheLayers, redisCache)
	}

	if config.Cache.BigCache.Enabled {
		bigCacheInstance, err := bigcache.NewBigCache(bigcache.Config{
			Shards:             config.Cache.BigCache.Shards,
			LifeWindow:         config.Cache.BigCache.ExpirationTime,
			MaxEntriesInWindow: config.Cache.BigCache.MaxEntriesInWindow,
			MaxEntrySize:       config.Cache.BigCache.MaxEntrySize,
			Verbose:            config.Cache.BigCache.Verbose,
			HardMaxCacheSize:   config.Cache.BigCache.HardMaxCacheSize,
		})
		if err != nil {
			return nil, errors.Wrap(err, "fail to initialize big cache")
		}

		cacheLayers = append(cacheLayers, adaptors.NewBigCacheAdaptor(bigCacheInstance))

	}

	cacheInstance := multilayercache.New(cacheLayers...)

	cacheInstance = middlewares.NewInstrumentationMiddleware(
		cacheInstance, prometheus_metirc.CacheMetrics.With(map[string]string{
			"cache_type": "multilayer",
		}))

	return cacheInstance, nil
}

func panicWithError(err error, format string, args ...interface{}) {
	logrus.WithError(err).Panicf(format, args...)
}

func providePrometheus(config *config.Config) *prometheus.Server {
	if config.Prometheus.Enabled {
		server := prometheus.NewServer(config.Prometheus.Port)
		go func() {
			err := server.Serve()
			if err != nil {
				panicWithError(err, "failed to start prometheus server")
			}
		}()
		return server
	}
	return nil
}

func provideConfig(cmd *cobra.Command) (*config.Config, error) {
	serviceConfig, err := config.LoadConfig(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load configurations.")
	}
	return serviceConfig, nil
}

func provideLogger(config *config.Config) (*logrus.Logger, error) {
	logger := logrus.New()
	if config.Logging.Level != "" {
		level, err := logrus.ParseLevel(config.Logging.Level)
		if err != nil {
			return nil, err
		}
		logger.SetLevel(level)
	}

	logger.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: false,
	})
	return logger, nil
}
