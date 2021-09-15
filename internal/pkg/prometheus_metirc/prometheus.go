package prometheus_metirc

import (
	"userProfile/pkg/metrics"
	prometheuspkg "userProfile/pkg/metrics/prometheus"
)

var (
	CacheMetrics metrics.Observer

	UserProviderMetrics metrics.Observer
)

func init() {
	CacheMetrics = prometheuspkg.NewHistogram("userProfile_cache",
		"view metrics about cache", "cache_type", "method", "ok", "success")

	UserProviderMetrics = prometheuspkg.NewHistogram("userProfile_provider",
		"view metrics about userProfile", "provider_type", "method", "ok", "success")
}
