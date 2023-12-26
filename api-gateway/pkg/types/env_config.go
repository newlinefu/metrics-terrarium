package types

import (
	"os"
	"strconv"
)

type Config struct {
	ApiGatewayHTTPPort string
	MetricsManagerPort string
	RawLifePeriod      int
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
		ApiGatewayHTTPPort: getEnv("API_GATEWAY_HTTP_PORT", ":8889"),
		MetricsManagerPort: getEnv("METRICS_MANAGER_PORT", ":8890"),
		RawLifePeriod:      getEnvAsInt("RAW_LIFE_PERIOD", 300),
	}
}
