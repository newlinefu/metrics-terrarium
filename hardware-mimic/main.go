package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"metricsTerrarium/hardware-mimic/pkg/types"
	"net/http"
)

var config *types.Config

type HardwareMimicAvailabilityMetric struct {
	Availability bool `json:"availability"`
}

type HardwareMimicSpeedMetric struct {
	Speed float32 `json:"speed"`
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
		Speed: rand.Float32() * 1000,
	})
	if err != nil {
		log.Printf("Error sending response. Err: %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func main() {
	config = types.CreateConfig()

	http.HandleFunc("/api/v1/trigger-availability-metric", handleTriggerAvailabilityMetrics)
	http.HandleFunc("/api/v1/trigger-speed-metric", handleTriggerSpeedMetrics)

	err := http.ListenAndServe(config.HardwareMimicPort, nil)
	if err != nil {
		print(err)
	}
}
