package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"metricsTerrarium/lib"
	"net/http"
	"time"
)

var config *lib.Config

type HardwareMimicAvailabilityMetric struct {
	Availability bool      `json:"availability"`
	Timestamp    time.Time `json:"timestamp"`
}

type HardwareMimicSpeedMetric struct {
	Speed     float32   `json:"speed"`
	Timestamp time.Time `json:"timestamp"`
}

func handleTriggerAvailabilityMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var availability bool
	randAvailability := rand.Float32()
	if randAvailability > 0.3 {
		availability = true
	} else {
		availability = false
	}
	err := json.NewEncoder(w).Encode(HardwareMimicAvailabilityMetric{
		Availability: availability,
		Timestamp:    time.Now(),
	})
	if err != nil {
		log.Printf("Error sending response. Err: %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func handleTriggerSpeedMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(HardwareMimicSpeedMetric{
		Speed:     rand.Float32() * 1000,
		Timestamp: time.Now(),
	})
	if err != nil {
		log.Printf("Error sending response. Err: %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func main() {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error with env variables file definition. Err: %s", err)
	} else {
		log.Printf("ENV variables initialized succesfully")
	}

	config = lib.CreateConfig()

	http.HandleFunc("/api/v1/trigger-availability-metric", handleTriggerAvailabilityMetrics)
	http.HandleFunc("/api/v1/trigger-speed-metric", handleTriggerSpeedMetrics)

	err = http.ListenAndServe(config.HardwareMimicPort, nil)
	if err != nil {
		print(err)
	}
}
