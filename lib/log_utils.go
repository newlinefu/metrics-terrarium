package lib

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func SetupLogs() (*os.File, error) {
	reg := regexp.MustCompile(`[-:+]`)
	timeStr := reg.ReplaceAllString(time.Now().Format(time.RFC3339), "_")

	logPath := filepath.Join(".", "logs", timeStr)
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	log.SetFlags(log.Ldate | log.Ltime)
	log.SetOutput(f)
	log.Println("Setup logs finished")

	return f, nil
}
