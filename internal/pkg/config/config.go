package config

import (
	"os"
	"strings"
	"sync"
	"time"
	"userProfile/pkg/sql"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type BasicConfigChangeListener interface {
	OnConfigChanged()
}

var changeListeners []BasicConfigChangeListener
var ChangeConfigMutex = &sync.RWMutex{}

// Config the application's configuration structure
type Config struct {
	Logging             LoggingConfig
	ConfigFile          string
	ListenPort          int
	Prometheus          PrometheusConfig
	GrpcHealthCheckPort int
	Database            DatabaseConfig
	Cache               CacheConfig
}

func (c *Config) OnConfigChanged() {
	c.updateSettings()
}

// LoggingConfig the logger's configuration structure
type LoggingConfig struct {
	SentryEnabled bool
	Level         string
}

type DatabaseConfig struct {
	WriteClient sql.PostgresConfig
	ReadClients sql.PostgresConfig
}

type CacheConfig struct {
	Redis    RedisConfig
	BigCache BigCacheConfig
}
type PrometheusConfig struct {
	Enabled bool
	Port    int
}

type RedisConfig struct {
	Enabled        bool
	Host           string
	Port           int
	DB             int
	Prefix         string
	Password       string
	ExpirationTime time.Duration
}

type BigCacheConfig struct {
	Enabled            bool
	ExpirationTime     time.Duration
	MaxSpace           int
	Shards             int
	LifeWindow         time.Duration
	MaxEntriesInWindow int
	MaxEntrySize       int
	Verbose            bool
	HardMaxCacheSize   int
}

// LoadConfig loads the config from a file if specified, otherwise from the environment
func LoadConfig(cmd *cobra.Command) (*Config, error) { // todo
	// Setting defaults for this application
	viper.SetDefault("logging.SentryEnabled", false)
	viper.SetDefault("logging.level", "error")
	viper.SetDefault("listenPort", 8080)
	viper.SetDefault("prometheus.metricListenPort", 8081)
	viper.SetDefault("prometheus.enabled", true)

	viper.SetDefault("database.write.host", "127.0.0.1")
	viper.SetDefault("database.write.port", 5432)
	viper.SetDefault("database.write.database", "userProfile")
	viper.SetDefault("database.write.ssl", false)
	viper.SetDefault("database.read.host", "127.0.0.1")
	viper.SetDefault("database.read.port", 5432)
	viper.SetDefault("database.read.database", "userProfile")
	viper.SetDefault("database.read.ssl", false)

	viper.SetDefault("cache.redis.host", "127.0.0.1")
	viper.SetDefault("cache.redis.port", 6379)
	viper.SetDefault("cache.redis.db", 0)
	viper.SetDefault("cache.redis.expirationTime", 3*time.Hour)
	viper.SetDefault("cache.redis.prefix", "USER_VIEW")
	viper.SetDefault("cache.redis.enabled", true)

	viper.SetDefault("cache.bigCache.shards", 1024)
	viper.SetDefault("cache.bigCache.maxEntriesInWindow", 1100*10*60)
	viper.SetDefault("cache.bigCache.maxEntrySize", 500)
	viper.SetDefault("cache.bigCache.verbose", true)
	viper.SetDefault("cache.bigCache.hardMaxCacheSize", 125)
	viper.SetDefault("cache.bigCache.enabled", true)

	// Read DbConfig from ENV
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("UserProfile")
	viper.AutomaticEnv()

	// Read DbConfig from Flags
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return nil, err
	}

	configFile, err := cmd.Flags().GetString("config-file")
	if err == nil && configFile != "" {
		viper.SetConfigFile(configFile)
		viper.WatchConfig()
		viper.OnConfigChange(configChanged)
		if err := viper.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	var config Config
	config.ConfigFile = configFile
	config.updateSettings()

	return &config, nil
}

func (c *Config) updateSettings() {

	DBUsername, ok := os.LookupEnv("DATABASE_USER")

	if !ok {
		log.Error("failed to get database username")
		DBUsername = ""
	}
	c.Database.WriteClient.Username = DBUsername
	c.Database.ReadClients.Username = DBUsername

	DBPassword, ok := os.LookupEnv("DATABASE_PASSWORD")

	if !ok {
		log.Error("failed to get database password")
		DBPassword = ""
	}
	c.Database.WriteClient.Password = DBPassword
	c.Database.ReadClients.Password = DBPassword

	c.Logging.SentryEnabled = viper.GetBool("logging.SentryEnabled")
	c.Logging.Level = viper.GetString("logging.level")
	c.ListenPort = viper.GetInt("listenPort")
	c.Prometheus.Port = viper.GetInt("prometheus.metricListenPort")
	c.Prometheus.Enabled = viper.GetBool("prometheus.enabled")

	c.Database.WriteClient.Host = viper.GetString("database.write.host")
	c.Database.WriteClient.Port = viper.GetInt("database.write.port")
	c.Database.WriteClient.Database = viper.GetString("database.write.database")
	c.Database.WriteClient.SSL = viper.GetBool("database.write.ssl")
	c.Database.ReadClients.Host = viper.GetString("database.read.host")
	c.Database.ReadClients.Port = viper.GetInt("database.read.port")
	c.Database.ReadClients.Database = viper.GetString("database.read.database")
	c.Database.ReadClients.SSL = viper.GetBool("database.read.ssl")

	c.Cache.Redis.Host = viper.GetString("cache.redis.host")
	c.Cache.Redis.Port = viper.GetInt("cache.redis.port")
	c.Cache.Redis.DB = viper.GetInt("cache.redis.db")
	c.Cache.Redis.ExpirationTime = viper.GetDuration("cache.redis.expirationTime")
	c.Cache.Redis.Prefix = viper.GetString("cache.redis.prefix")
	c.Cache.Redis.Enabled = viper.GetBool("cache.redis.enabled")

	c.Cache.BigCache.Shards = viper.GetInt("cache.bigCache.shards")
	c.Cache.BigCache.MaxEntriesInWindow = viper.GetInt("cache.bigCache.maxEntriesInWindow")
	c.Cache.BigCache.MaxEntrySize = viper.GetInt("cache.bigCache.maxEntrySize")
	c.Cache.BigCache.Verbose = viper.GetBool("cache.bigCache.verbose")
	c.Cache.BigCache.HardMaxCacheSize = viper.GetInt("cache.bigCache.hardMaxCacheSize")
	c.Cache.BigCache.Enabled = viper.GetBool("cache.bigCache.enabled")

}

func AddToChangeListener(listener BasicConfigChangeListener) {
	ChangeConfigMutex.Lock()
	defer ChangeConfigMutex.Unlock()
	changeListeners = append(changeListeners, listener)
}

func configChanged(fsnotify.Event) {
	log.Infof("Config Changed.... Reloading Config")
	ChangeConfigMutex.Lock()
	defer ChangeConfigMutex.Unlock()
	for _, listener := range changeListeners {
		listener.OnConfigChanged()
	}
}
