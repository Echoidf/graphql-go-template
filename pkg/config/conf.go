package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	ProductionEnv  = "production"
	DevelopmentEnv = "development"

	DatabaseTimeout    = 5 * time.Second
	ProductCachingTime = 1 * time.Minute
)

type Config struct {
	ServerPort  int         `yaml:"server_port"`
	SocketPath  string      `yaml:"socket_path"`
	Environment string      `yaml:"environment"`
	Logger      Logger      `yaml:"logger"`
	Mysql       MysqlConfig `yaml:"mysql"`
}

type (
	Logger struct {
		Level      string `yaml:"level"`
		LogPath    string `yaml:"log_path"`
		MaxSize    int    `yaml:"max_size"`
		MaxAge     int    `yaml:"max_age"`
		MaxBackups int    `yaml:"max_backups"`
	}

	MysqlConfig struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	}
)

var (
	cfg  Config
	once sync.Once
)

func LoadConfig() *Config {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)

	yamlFile, err := os.ReadFile(filepath.Join(currentDir, "config.yaml"))
	if err != nil {
		log.Printf("Error on reading configuration file, error: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalf("Error on parsing configuration file, error: %v", err)
	}

	return &cfg
}

func GetConfig() *Config {
	once.Do(func() {
		cfg = *LoadConfig()
	})
	return &cfg
}
