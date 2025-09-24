package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env        string     `yaml:"env"`
	JWTSecret  []byte     `yaml:"jwt_secret"`
	HTTPServer HTTPServer `yaml:"http_server"`
	Database   Database   `yaml:"database"`
}

type HTTPServer struct {
	Port         string        `yaml:"port"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type Database struct {
	DSN      string `yaml:"dsn"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SSLMode  string `yaml:"sslmode"`
}

var configPath string = "./config/config.yaml"

func NewConfig() *Config {

	if os.Getenv("CONFIG_PATH") != "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	v := viper.New()
	v.SetConfigFile(configPath)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetTypeByDefaultValue(true)

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	config := &Config{
		Env:       v.GetString("env"),
		JWTSecret: []byte(v.GetString("jwt_secret")),
		HTTPServer: HTTPServer{
			Port:         v.GetString("http_server.port"),
			WriteTimeout: v.GetDuration("http_server.write_timeout"),
			ReadTimeout:  v.GetDuration("http_server.read_timeout"),
			IdleTimeout:  v.GetDuration("http_server.idle_timeout"),
		},
		Database: Database{
			DSN:      v.GetString("database.dsn"),
			Host:     v.GetString("database.host"),
			Port:     v.GetInt("database.port"),
			User:     v.GetString("database.user"),
			Password: v.GetString("database.password"),
			Name:     v.GetString("database.name"),
			SSLMode:  v.GetString("database.sslmode"),
		},
	}

	return config
}
