package config

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
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

type JWTConfig struct {
	SecretKey   string `mapstructure:"secret_key"`
	ExpiresTime int    `mapstructure:"expires_time"`
}

type GatewayConfig struct {
	Addr string `mapstructure:"addr"`
}

type ServerConfig struct {
	Name string
	Mode string
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
	} `mapstructure:"video"`
}

type VideoConfig struct {
	Name    string
	RPCAddr string `mapstructure:"rpc_addr"`
}

type EtcdConfig struct {
	Addr string `mapstructure:"addr"`
}

var (
	Server  ServerConfig
	MySQL   MySQLConfig
	Redis   RedisConfig
	JWT     JWTConfig
	Upload  UploadConfig
	User    UserConfig
	Video   VideoConfig
	Gateway GatewayConfig
	Etcd    EtcdConfig
)

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.UnmarshalKey("server", &Server); err != nil {
		panic(err)
	}
	if err := viper.UnmarshalKey("mysql", &MySQL); err != nil {
		panic(err)
	}
	if err := viper.UnmarshalKey("redis", &Redis); err != nil {
		panic(err)
	}
	if err := viper.UnmarshalKey("jwt", &JWT); err != nil {
		panic(err)
	}
	if err := viper.UnmarshalKey("upload", &Upload); err != nil {
		panic(err)
	}
	if err := viper.UnmarshalKey("user", &User); err != nil {
		panic(err)
	}
	if err := viper.UnmarshalKey("video", &Video); err != nil {
		panic(err)
	}
	if err := viper.UnmarshalKey("gateway", &Gateway); err != nil {
		panic(err)
	}
	if err := viper.UnmarshalKey("etcd", &Etcd); err != nil {
		panic(err)
	}
}

func GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MySQL.Username,
		MySQL.Password,
		MySQL.Host,
		MySQL.Port,
		MySQL.Database,
	)
}

func GetClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", Redis.Host, Redis.Port),
		Password: Redis.Password,
		DB:       Redis.DB,
	})

	return client
}
