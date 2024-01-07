package metrics

import (
	"database/sql"
	"log"
	"metricsTerrarium/services/metrics-manager/internal/general_types"
	"metricsTerrarium/services/metrics-manager/internal/service/metrics_service"
	"metricsTerrarium/services/metrics-manager/pkg/api"
)

func GetRawMetrics(rawMetrics map[string]general_types.RawMetric) (*api.MetricsResponse, error) {
	log.Printf("Get Raw metrics request handled successfuly")

	return metrics_service.GetRawMetrics(rawMetrics)
}

func GetPreparedMetrics(req *api.PreparedMetricsRequestMessage, dbConnection *sql.DB) (*api.MetricsResponse, error) {
	log.Printf("Get Prepared metrics request handled successfuly")

	return metrics_service.GetPreparedMetrics(req, dbConnection)
}
