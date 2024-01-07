package metrics

import (
	"database/sql"
	"log"
	"metricsTerrarium/services/metrics-manager/internal/general_types"
	"metricsTerrarium/services/metrics-manager/internal/service/metrics_service"
	"metricsTerrarium/services/metrics-manager/pkg/api"
)

func GetRawMetrics(rawMetrics map[string]general_types.RawMetric, rawLifePeriod int) (*api.MetricsResponse, error) {
	log.Printf("Get Raw metrics request handled successfuly")

	return metrics_service.GetRawMetrics(rawMetrics, rawLifePeriod)
}

func GetPreparedMetrics(req *api.PreparedMetricsRequestMessage, dbConnection *sql.DB) (*api.MetricsResponse, error) {
	log.Printf("Get Prepared metrics request handled successfuly")

	return metrics_service.GetPreparedMetrics(req, dbConnection)
}
