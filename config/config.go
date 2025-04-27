package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

type MySQLConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type KafkaConfig struct {
	Network  string
	Address  string
	User     string
	Password string
}

type JWTConfig struct {
	SecretKey   string `mapstructure:"secret"`
	ExpiresTime int    `mapstructure:"expires_time"`
}

type GatewayConfig struct {
	Addr string `mapstructure:"addr"`
}

type ServerConfig struct {
	Name           string
	Mode           string
	MaxConnections int64 `mapstructure:"max_connections"`
	MaxQPS         int64 `mapstructure:"max_qps"`
}

type OtelConfig struct {
	CollectorAddr string `mapstructure:"collector_addr"`
}

type UserConfig struct {
	Name          string
	RPCAddr       string `mapstructure:"rpc_addr"`
	DefaultAvatar string `mapstructure:"default_avatar"`
}

type UploadConfig struct {
	Avatar struct {
		MaxSize      int      `mapstructure:"max_size"`
		AllowedTypes []string `mapstructure:"allowed_types"`
		UploadDir    string   `mapstructure:"upload_dir"`
		BaseURL      string   `mapstructure:"base_url"`
	} `mapstructure:"avatar"`
	Video struct {
		MaxSize      int      `mapstructure:"max_size"`
		AllowedTypes []string `mapstructure:"allowed_types"`
		UploadDir    string   `mapstructure:"upload_dir"`
		BaseURL      string   `mapstructure:"base_url"`
		CoverDir     string   `mapstructure:"cover_dir"`
	} `mapstructure:"video"`
}

type VideoConfig struct {
	Name    string
	RPCAddr string `mapstructure:"rpc_addr"`
}

type SocialConfig struct {
	Name     string
	RPCAddr  string `mapstructure:"rpc_addr"`
	HTTPAddr string `mapstructure:"http_addr"`
	Chat     struct {
		MaxGroupMembers   int `mapstructure:"max_group_members"`
		MaxMessageLength  int `mapstructure:"max_message_length"`
		MaxGroupsPerUser  int `mapstructure:"max_groups_per_user"`
		MaxFriendsPerUser int `mapstructure:"max_friends_per_user"`
		MessagePageSize   int `mapstructure:"message_page_size"`
		HistoryDays       int `mapstructure:"history_days"`
	} `mapstructure:"chat"`
	File struct {
		MaxSize      int      `mapstructure:"max_size"`
		AllowedTypes []string `mapstructure:"allowed_types"`
		UploadDir    string   `mapstructure:"upload_dir"`
		BaseURL      string   `mapstructure:"base_url"`
	} `mapstructure:"file"`
}

type EtcdConfig struct {
	Addr string `mapstructure:"addr"`
}

type VideoInteractionsConfig struct {
	Name    string
	RPCAddr string `mapstructure:"rpc_addr"`
}

var (
	Server            *ServerConfig
	MySQL             *MySQLConfig
	Redis             *RedisConfig
	Kafka             *KafkaConfig
	JWT               *JWTConfig
	Upload            *UploadConfig
	User              *UserConfig
	Video             *VideoConfig
	Social            *SocialConfig
	Gateway           *GatewayConfig
	Etcd              *EtcdConfig
	VideoInteractions *VideoInteractionsConfig
	Otel              *OtelConfig
	runtimeViper      = viper.New()
)

const (
	remoteProvider = "etcd3"
	remotePath     = "/config"
	remoteFileName = "config"
	remoteFileType = "yaml"
)

type Config struct {
	Server            ServerConfig
	MySQL             MySQLConfig
	Redis             RedisConfig
	Kafka             KafkaConfig
	JWT               JWTConfig
	Upload            UploadConfig
	User              UserConfig
	Video             VideoConfig
	Social            SocialConfig
	Gateway           GatewayConfig
	Etcd              EtcdConfig
	VideoInteractions VideoInteractionsConfig
	Otel              OtelConfig
}

// Init 初始化配置，支持从etcd远程获取配置
func Init(service string) {
	// 从环境变量获取etcd地址
	etcdAddr := os.Getenv("ETCD_ADDR")
	if etcdAddr == "" {
		panic("ETCD_ADDR environment variable is not set")
	}
	fmt.Printf("Using etcd address: %s\n", etcdAddr)

	// 配置远程provider
	err := runtimeViper.AddRemoteProvider(remoteProvider, etcdAddr, remotePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to add remote provider: %v", err))
	}

	runtimeViper.SetConfigName(remoteFileName)
	runtimeViper.SetConfigType(remoteFileType)

	// 读取远程配置
	if err := runtimeViper.ReadRemoteConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			panic("Could not find config files")
		}
		panic(fmt.Sprintf("Failed to read remote config: %v", err))
	}

	// 映射配置到结构体
	configMapping()

	// 设置配置变更监听
	runtimeViper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file changed: %s\n", e.Name)
		configMapping() // 重新加载配置
	})
	runtimeViper.WatchConfig()
}

// configMapping 将配置映射到全局变量
func configMapping() {
	conf := &Config{}
	if err := runtimeViper.Unmarshal(conf); err != nil {
		panic(fmt.Sprintf("Failed to unmarshal config: %v", err))
	}

	// 更新全局配置变量
	Server = &conf.Server
	MySQL = &conf.MySQL
	Redis = &conf.Redis
	Kafka = &conf.Kafka
	JWT = &conf.JWT
	Upload = &conf.Upload
	User = &conf.User
	Video = &conf.Video
	Social = &conf.Social
	Gateway = &conf.Gateway
	Etcd = &conf.Etcd
	VideoInteractions = &conf.VideoInteractions
	Otel = &conf.Otel
}

// GetDSN 获取MySQL连接字符串
func GetDSN() string {
	if MySQL == nil {
		panic("MySQL config is not initialized")
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MySQL.Username,
		MySQL.Password,
		MySQL.Host,
		MySQL.Port,
		MySQL.Database,
	)
}

// GetRedisClient 获取Redis客户端实例
func GetRedisClient() *redis.Client {
	if Redis == nil {
		panic("Redis config is not initialized")
	}
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", Redis.Host, Redis.Port),
		Password: Redis.Password,
		DB:       Redis.DB,
	})
}
