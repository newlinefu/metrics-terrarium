package types

import (
	"os"
)

type Config struct {
	HardwareMimicPort string
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func CreateConfig() *Config {
	return &Config{
		HardwareMimicPort: getEnv("HARDWARE_MIMIC_PORT", ":8891"),
	}
}
