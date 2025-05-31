package config

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

type AppConfig struct {
	Name string
	Port string
}

type DbConfig struct {
	Host            string
	Port            string
	User            string
	Pass            string
	Schema          string
	MaxIdleConn     int
	MaxOpenConn     int
	MaxConnLifetime time.Duration
	Debug           bool
}

type RedisConfig struct {
	Host               string
	Port               string
	Pass               string
	Db                 int
	MandatoryPrefix    string
	AccessUuidPrefix   string
	RefreshUuidPrefix  string
	UserPrefix         string
	PermissionPrefix   string
	UserCacheTTL       time.Duration
	PermissionCacheTTL time.Duration
}

type AsynqConfig struct {
	RedisAddr   string
	DB          int
	Pass        string
	Concurrency int
	Queue       string
	Retention   time.Duration // in hours
	RetryCount  int
	Delay       time.Duration // in seconds
}

type JwtConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

type LoggerConfig struct {
	Level    string
	FilePath string
}

type Config struct {
	App    *AppConfig
	DB     *DbConfig
	Redis  *RedisConfig
	Asynq  *AsynqConfig
	Logger *LoggerConfig
	Jwt    *JwtConfig
}

var config Config

func GetAll() Config {
	return config
}

func App() *AppConfig {
	return config.App
}

func Db() *DbConfig {
	return config.DB
}

func Redis() *RedisConfig {
	return config.Redis
}

func Asynq() *AsynqConfig {
	return config.Asynq
}

func Logger() *LoggerConfig {
	return config.Logger
}

func Jwt() *JwtConfig {
	return config.Jwt
}

func LoadConfig() {
	setDefaultConfig()

	const (
		ENV_CONSUL_URL  = "CONSUL_URL"
		ENV_CONSUL_PATH = "CONSUL_PATH"
	)

	_ = viper.BindEnv(ENV_CONSUL_URL)
	_ = viper.BindEnv(ENV_CONSUL_PATH)

	consulURL := viper.GetString(ENV_CONSUL_URL)
	consulPath := viper.GetString(ENV_CONSUL_PATH)

	if consulURL != "" && consulPath != "" {
		_ = viper.AddRemoteProvider("consul", consulURL, consulPath)

		viper.SetConfigType("json")
		err := viper.ReadRemoteConfig()

		if err != nil {
			log.Printf("%s named \"%s\"\n", err.Error(), consulPath)
			panic(err)
		}

		config = Config{}

		if err := viper.Unmarshal(&config); err != nil {
			panic(err)
		}

		if r, err := json.MarshalIndent(&config, "", "  "); err == nil {
			fmt.Println(string(r))
		}
	} else {
		log.Println("CONSUL_URL or CONSUL_PATH missing! Serving with default config...")
	}
}

func setDefaultConfig() {
	config.App = &AppConfig{
		Name: "app",
		Port: "8080",
	}

	config.DB = &DbConfig{
		Host:            "127.0.0.1",
		Port:            "3306",
		User:            "gladiator",
		Pass:            "gladius",
		Schema:          "event_managment",
		MaxIdleConn:     1,
		MaxOpenConn:     2,
		MaxConnLifetime: 30,
		Debug:           true,
	}

	config.Redis = &RedisConfig{
		Host:            "127.0.0.1",
		Port:            "6379",
		Pass:            "",
		Db:              2,
		MandatoryPrefix: "event_management_",
	}

	config.Asynq = &AsynqConfig{
		RedisAddr:   "127.0.0.1:6379",
		DB:          2,
		Pass:        "",
		Concurrency: 10,
		Queue:       "app",
		Retention:   168,
		RetryCount:  25,
		Delay:       120,
	}
}
