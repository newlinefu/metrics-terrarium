package main

import (
	"encoding/json"
	"github.com/sony/gobreaker"
	"log"
	"metricsTerrarium/hardware-circuit-breaker/pkg/types"
	"metricsTerrarium/hardware-circuit-breaker/pkg/util"
	"net/http"
)

var cb *gobreaker.CircuitBreaker
var config *types.Config

type AvailabilityMetricResponse struct {
	Availability bool `json:"availability"`
}

type SpeedMetricResponse struct {
	Speed float32 `json:"speed"`
}

func initCircuitBreaker() {
	var st gobreaker.Settings
	st.Name = "Hardware Mimic Metrics Breaker"
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return failureRatio >= 0.5
	}
	st.OnStateChange = func(name string, from gobreaker.State, to gobreaker.State) {
		log.Printf("State of circuit breaker changed from [%s] to [%s]\n", from, to)
	}

	cb = gobreaker.NewCircuitBreaker(st)
}

func HandleAvailabilityMetric(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start of handling availability metric request")
	w.Header().Set("Content-Type", "application/json")
	body, err := cb.Execute(func() (interface{}, error) {
		response, err := http.Get("http://localhost" + config.HardwareMimicPort + "/api/v1/trigger-availability-metric")
		if err != nil {
			return nil, err
		}
		var availabilityBody *AvailabilityMetricResponse
		err = json.NewDecoder(response.Body).Decode(&availabilityBody)
		if err != nil {
			return nil, err
		}
		return availabilityBody, nil
	})
	err = json.NewEncoder(w).Encode(body)
	if err != nil {
		log.Printf("Error during encoding. Err: %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func HandleSpeedMetric(w http.ResponseWriter, r *http.Request) {
	log.Printf("Start of handling speed metric request")
	w.Header().Set("Content-Type", "application/json")

	body, err := cb.Execute(func() (interface{}, error) {
		response, err := http.Get("http://localhost" + config.HardwareMimicPort + "/api/v1/trigger-speed-metric")
		if err != nil {
			return nil, err
		}
		var speedBody *SpeedMetricResponse
		err = json.NewDecoder(response.Body).Decode(&speedBody)
		if err != nil {
			return nil, err
		}
		return speedBody, nil
	})
	err = json.NewEncoder(w).Encode(body)
	if err != nil {
		log.Printf("Error during encoding. Err: %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func main() {
	config = types.CreateConfig()

	logFile, err := util.SetupLogs()
	defer logFile.Close()
	if err != nil {
		log.Fatalf("Error with logs initialization occured. Err: %s", err)
	}

	initCircuitBreaker()
	http.HandleFunc("/api/v1/trigger-availability-metric", HandleAvailabilityMetric)
	http.HandleFunc("/api/v1/trigger-speed-metric", HandleSpeedMetric)

	err = http.ListenAndServe(config.HardwareCircuitBreakerPort, nil)
	if err != nil {
		log.Fatalf("Error during listening start occured. Err: %s", err)
	}
}
