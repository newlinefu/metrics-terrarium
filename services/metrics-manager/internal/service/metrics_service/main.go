package metrics_service

import (
	"database/sql"
	"log"
	"metricsTerrarium/services/metrics-manager/internal/general_types"
	"metricsTerrarium/services/metrics-manager/internal/storage/metrics_storage"
	"metricsTerrarium/services/metrics-manager/internal/util"
	"metricsTerrarium/services/metrics-manager/pkg/api"
)

func GetPreparedMetrics(req *api.PreparedMetricsRequestMessage, dbConnection *sql.DB) (*api.MetricsResponse, error) {
	log.Printf("Get Prepared metrics request handled successfuly")

	avgSpeed, avgAvailability, counter, err := metrics_storage.GetPreparedMetrics(req, dbConnection)
	return &api.MetricsResponse{
		Availability: avgAvailability / float32(counter),
		Speed:        avgSpeed / float32(counter),
	}, err
}

func GetRawMetrics(rawMetrics *map[string]general_types.RawMetric) (*api.MetricsResponse, error) {
	avgSpeed, lastAvailability := util.GetPreparedMetrics(rawMetrics)
	availabilityResponse := 0.0
	if lastAvailability {
		availabilityResponse = 1.0
	}
	return &api.MetricsResponse{
		Availability: float32(availabilityResponse),
		Speed:        avgSpeed,
	}, nil
}
