package config

import (
    "time"
    "github.com/spf13/viper"
)

type HTTP struct {
    Port        string        `mapstructure:"port"`
    ReadTimeout time.Duration `mapstructure:"read_timeout"`
    WriteTimeout time.Duration `mapstructure:"write_timeout"`
    IdleTimeout time.Duration `mapstructure:"idle_timeout"`
}

type Database struct {
    DSN string `mapstructure:"dsn"`
}

type Auth struct {
    JWTSecret string `mapstructure:"jwt_secret"`
}

type Config struct {
    HTTPServer HTTP    `mapstructure:"http"`
    Database   Database `mapstructure:"database"`
    Auth       Auth     `mapstructure:"auth"`
}

func Load() (*Config, error) {
    v := viper.New()
    v.SetConfigName("config")
    v.SetConfigType("yaml")
    v.AddConfigPath("./config")
    v.SetDefault("http.port", ":8080")
    v.SetDefault("http.read_timeout", "10s")
    v.SetDefault("http.write_timeout", "10s")
    v.SetDefault("http.idle_timeout", "60s")
    v.SetDefault("auth.jwt_secret", "devsecret")

    if err := v.ReadInConfig(); err != nil {
        // allow missing file; env-only configs
    }
    v.AutomaticEnv()

    var cfg Config
    if err := v.Unmarshal(&cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}
