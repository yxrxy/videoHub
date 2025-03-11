package config

import (
	"fmt"
	"strings"

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

type ServerConfig struct {
	Name string
	Mode string
}

type UserConfig struct {
	Name          string
	HTTPAddr      string `mapstructure:"http_addr"`
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
}

type VideoConfig struct {
	Name     string
	HTTPAddr string `mapstructure:"http_addr"`
	RPCAddr  string `mapstructure:"rpc_addr"`
}

var (
	Server ServerConfig
	MySQL  MySQLConfig
	Redis  RedisConfig
	JWT    JWTConfig
	Upload UploadConfig
	User   UserConfig
	Video  VideoConfig
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

func GetMySQLEnv() map[string]string {
	return map[string]string{
		"MYSQL_ROOT_PASSWORD": MySQL.Password,
		"MYSQL_DATABASE":      MySQL.Database,
		"MYSQL_HOST":          MySQL.Host,
		"MYSQL_PORT":          fmt.Sprintf("%d", MySQL.Port),
	}
}

func GetUserServiceEnv() map[string]string {
	return map[string]string{
		"SERVICE_NAME": User.Name,
		"SERVICE_PORT": strings.TrimPrefix(User.HTTPAddr, ":"),
	}
}
