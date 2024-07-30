package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string        `yaml:"env" env-default:"local"`
	AccessTokenTTL string `yaml:"access_token_ttl" env-required:"true"`
	RefreshTokenTTL string `yaml:"refresh_token_ttl" env-required:"true"`
	Database Database `yaml:"database" env-required:"true"`
	Grpc Grpc `yaml:"grpc" env-required:"true"`
}

type Database struct {
	MongoUser string `yaml:"mongo_user" env-required:"true"`
	MongoPassword string `yaml:"mongo_password" env-required:"true"`
	MongoHostname string `yaml:"mongo_hostname" env-required:"true"`
	MongoPort int `yaml:"mongo_port" env-required:"true"`
	MongoDBName string `yaml:"mongo_db_name" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-defailt:"1"`
}

type Grpc struct {
	Port int `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-required:"true"`
	GrpcHost string `yaml:"grpc_host" env-required:"true"`
}

func LoadConfig() (*Config, error) {
	var cfg Config


	err := fetchConfigByPath(&cfg, os.Getenv("CONFIG_PATH"))

	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func fetchConfigByPath(cfg *Config, path string ) error {
	const op = "config.fetchConfigByPath"
	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}




