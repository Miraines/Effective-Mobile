package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Server     Server     `mapstructure:"server"`
		Database   Database   `mapstructure:"database"`
		Log        Log        `mapstructure:"log"`
		Enrichment Enrichment `mapstructure:"enrichment"`
		AppEnv     string     `mapstructure:"-"`
	}

	Server struct {
		Host         string        `mapstructure:"host"`
		Port         int           `mapstructure:"port"`
		ReadTimeout  time.Duration `mapstructure:"read_timeout"`
		WriteTimeout time.Duration `mapstructure:"write_timeout"`
	}

	Database struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
	}

	Log struct {
		Level string `mapstructure:"level"`
	}

	Enrichment struct {
		Timeout time.Duration `mapstructure:"timeout"`
		Retry   int           `mapstructure:"retry"`
	}
)

func Load() (*Config, error) {
	_ = godotenv.Load()

	env := viper.GetString("APP_ENV")
	if env == "" {
		env = "default"
	}

	v := viper.New()
	v.SetConfigName("config.default")
	v.AddConfigPath("configs")
	if err := v.MergeInConfig(); err != nil {
		return nil, err
	}

	if env != "default" {
		v.SetConfigName("config." + env)
		if err := v.MergeInConfig(); err != nil {
			return nil, err
		}
	}

	v.SetEnvPrefix("")
	v.AutomaticEnv()

	var cfg Config
	if err := v.UnmarshalExact(&cfg); err != nil {
		return nil, err
	}
	cfg.AppEnv = env
	return &cfg, nil
}
