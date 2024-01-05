package lib

import (
	"os"
	"strconv"
)

type Config struct {
	MetricsManagerPort         string
	MetricsDataBase            string
	MetricsDataBaseUser        string
	MetricsDataBasePassword    string
	RawLifePeriod              int
	HardwareCheckInterval      int
	KafkaAddress               string
	MetricsInterpreterPort     string
	HardwareCircuitBreakerPort string
	HardwareMimicPort          string
	ApiGatewayHTTPPort         string
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
		MetricsManagerPort:         getEnv("METRICS_MANAGER_PORT", ":8890"),
		MetricsDataBase:            getEnv("METRICS_DATABASE", "metrics"),
		MetricsDataBaseUser:        getEnv("METRICS_DATABASE_USER", "postgres"),
		MetricsDataBasePassword:    getEnv("METRICS_DATABASE_PASSWORD", "root"),
		RawLifePeriod:              getEnvAsInt("RAW_LIFE_PERIOD", 300),
		KafkaAddress:               getEnv("KAFKA_ADDRESS", "localhost:9092"),
		HardwareCheckInterval:      getEnvAsInt("HARDWARE_CHECK_INTERVAL", 5),
		MetricsInterpreterPort:     getEnv("METRICS_INTERPRETER_PORT", ":8893"),
		HardwareCircuitBreakerPort: getEnv("HARDWARE_CIRCUIT_BREAKER_PORT", ":8892"),
		HardwareMimicPort:          getEnv("HARDWARE_MIMIC_PORT", ":8891"),
		ApiGatewayHTTPPort:         getEnv("API_GATEWAY_HTTP_PORT", ":8889"),
	}
}
