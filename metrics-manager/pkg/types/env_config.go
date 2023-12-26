package types

import (
	"os"
	"strconv"
)

type Config struct {
	MetricsManagerPort      string
	MetricsDataBase         string
	MetricsDataBaseUser     string
	MetricsDataBasePassword string
	RawLifePeriod           int
	KafkaAddress            string
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

func CreateConfig() *Config {
	return &Config{
		MetricsManagerPort:      getEnv("METRICS_MANAGER_PORT", ":8890"),
		MetricsDataBase:         getEnv("METRICS_DATABASE", "metrics"),
		MetricsDataBaseUser:     getEnv("METRICS_DATABASE_USER", "postgres"),
		MetricsDataBasePassword: getEnv("METRICS_DATABASE_PASSWORD", "root"),
		RawLifePeriod:           getEnvAsInt("RAW_LIFE_PERIOD", 300),
		KafkaAddress:            getEnv("KAFKA_ADDRESS", "localhost:9092"),
	}
}
