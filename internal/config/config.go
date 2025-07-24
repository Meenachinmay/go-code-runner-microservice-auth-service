package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strconv"
)

type RawConfig struct {
	GrpcServerPort string `yaml:"grpc_server_port"`
	RequestTimeout int    `yaml:"request_timeout"`
	DBHost         string `yaml:"db_host"`
	DBPort         string `yaml:"db_port"`
	DBUser         string `yaml:"db_user"`
	DBPassword     string `yaml:"db_password"`
	DBName         string `yaml:"db_name"`
}

type Config struct {
	GrpcServerPort string
	RequestTimeout int
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBConnStr      string
}

func Load() (*Config, error) {
	env := os.Getenv("APP_ENVIRONMENT")
	if env == "" {
		env = "local"
	}

	cfgPath := filepath.Join("internal", "config", env+".yml")
	f, err := os.Open(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("open%s: %w", cfgPath, err)
	}
	defer f.Close()

	var raw RawConfig
	if err := yaml.NewDecoder(f).Decode(&raw); err != nil {
		return nil, fmt.Errorf("parse %s: %w", cfgPath, err)
	}

	if v := os.Getenv("GRPC_PORT"); v != "" {
		raw.GrpcServerPort = v
	}
	if v := os.Getenv("REQUEST_TIMEOUT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			raw.RequestTimeout = n
		}
	}

	// Database configuration from environment variables
	if v := os.Getenv("POSTGRES_HOST"); v != "" {
		raw.DBHost = v
	}
	if v := os.Getenv("POSTGRES_PORT"); v != "" {
		raw.DBPort = v
	}
	if v := os.Getenv("POSTGRES_USER"); v != "" {
		raw.DBUser = v
	}
	if v := os.Getenv("POSTGRES_PASSWORD"); v != "" {
		raw.DBPassword = v
	}
	if v := os.Getenv("POSTGRES_DB"); v != "" {
		raw.DBName = v
	}

	// Construct database connection string
	dbConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		raw.DBUser, raw.DBPassword, raw.DBHost, raw.DBPort, raw.DBName)

	return &Config{
		GrpcServerPort: raw.GrpcServerPort,
		RequestTimeout: raw.RequestTimeout,
		DBHost:         raw.DBHost,
		DBPort:         raw.DBPort,
		DBUser:         raw.DBUser,
		DBPassword:     raw.DBPassword,
		DBName:         raw.DBName,
		DBConnStr:      dbConnStr,
	}, nil
}
