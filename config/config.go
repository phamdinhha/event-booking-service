package config

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

// App config struct
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Logger     Logger           `mapstructure:"logger"`
	Postgres   PostgresConfig   `mapstructure:"postgres"`
	Redis      RedisConfig      `mapstructure:"redis"`
	Migrations MigrationsConfig `mapstructure:"migrations"`
}

type PostgresConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type MigrationsConfig struct {
	Path string `mapstructure:"path"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// Server config struct
type ServerConfig struct {
	Development bool
	AppVersion  string
	Host        string
	Port        string
	CorsOrigins []string
}

// Logger config
type Logger struct {
	Encoding string
	Level    string
}

// Load config file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

// Parse config file
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}

// Get config
func GetConfig(configPath string) (*Config, error) {
	cfgFile, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	cfg, err := ParseConfig(cfgFile)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func GetEnvConfig() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	return &Config{
		Server: ServerConfig{
			AppVersion:  v.GetString("SERVER_APPVERSION"),
			Host:        v.GetString("SERVER_HOST"),
			Port:        v.GetString("SERVER_PORT"),
			Development: v.GetBool("SERVER_DEVELOPMENT"),
			CorsOrigins: v.GetStringSlice("SERVER_CORS_ORIGINS"),
		},
		Logger: Logger{
			Encoding: v.GetString("LOGGER_ENCODING"),
			Level:    v.GetString("LOGGER_LEVEL"),
		},
		Postgres: PostgresConfig{
			Host:     v.GetString("POSTGRES_HOST"),
			Port:     v.GetString("POSTGRES_PORT"),
			User:     v.GetString("POSTGRES_USER"),
			Password: v.GetString("POSTGRES_PASSWORD"),
			Database: v.GetString("POSTGRES_DATABASE"),
			Driver:   v.GetString("POSTGRES_DRIVER"),
		},
		Redis: RedisConfig{
			Host:     v.GetString("REDIS_HOST"),
			Port:     v.GetString("REDIS_PORT"),
			Password: v.GetString("REDIS_PASSWORD"),
			DB:       v.GetInt("REDIS_DB"),
		},
	}, nil
}
