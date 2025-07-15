package config

import (
	"os"
	"strconv"
	"time"

	"zerofy.pro/rbac-collector/src/constants"
)

const (
	DefaultLogFormat = constants.LogFormatConsole
)

type Config struct {
	CollectionInterval time.Duration
	LogFormat          string
}

func Load() (*Config, error) {
	intervalStr := getEnv(constants.EnvCollectionIntervalSeconds, constants.DefaultCollectionIntervalSeconds)
	intervalSec, err := strconv.Atoi(intervalStr)
	if err != nil {
		return nil, err
	}

	logFormat := getEnv(constants.EnvLogFormat, DefaultLogFormat)

	return &Config{
		CollectionInterval: time.Duration(intervalSec) * time.Second,
		LogFormat:          logFormat,
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
