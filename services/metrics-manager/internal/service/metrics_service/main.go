package metrics_service

import (
	"database/sql"
	"metricsTerrarium/services/metrics-manager/internal/general_types"
	"metricsTerrarium/services/metrics-manager/internal/storage/metrics_storage"
	"metricsTerrarium/services/metrics-manager/internal/util"
	"metricsTerrarium/services/metrics-manager/pkg/api"
	"time"
)

func GetPreparedMetrics(req *api.PreparedMetricsRequestMessage, dbConnection *sql.DB) (*api.MetricsResponse, error) {
	avgSpeed, avgAvailability, err := metrics_storage.GetPreparedMetrics(req, dbConnection)
	return &api.MetricsResponse{
		Availability: avgAvailability,
		Speed:        avgSpeed,
		Type:         "prepared",
	}, err
}

func GetRawMetrics(rawMetrics map[string]general_types.RawMetric, rawLifePeriod int) (*api.MetricsResponse, error) {
	avgSpeed, avgAvailability := util.GetPreparedMetrics(rawMetrics, time.Now().Add(time.Duration(-rawLifePeriod)*time.Second), time.Now(), rawLifePeriod)
	return &api.MetricsResponse{
		Availability: avgAvailability,
		Speed:        avgSpeed,
		Type:         "raw",
	}, nil
}
