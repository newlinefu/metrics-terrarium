package types

import (
	"os"
	"strconv"
)

type Config struct {
	HardwareCircuitBreakerPort string
	MetricsInterpreterPort     string
	HardwareCheckInterval      int
	KafkaAddress               string
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
		HardwareCircuitBreakerPort: getEnv("HARDWARE_CIRCUIT_BREAKER_PORT", ":8892"),
		MetricsInterpreterPort:     getEnv("METRICS_INTERPRETER_PORT", ":8893"),
		HardwareCheckInterval:      getEnvAsInt("HARDWARE_CHECK_INTERVAL", 5),
		KafkaAddress:               getEnv("KAFKA_ADDRESS", "localhost:9092"),
	}
}
