package job

import (
	"database/sql"
	"metricsTerrarium/services/metrics-manager/internal/general_types"
	"metricsTerrarium/services/metrics-manager/internal/storage/metrics_storage"
	"metricsTerrarium/services/metrics-manager/internal/util"
	"time"
)

func PrepareMetrics(connection *sql.DB, rawLifePeriod int, rawMetrics map[string]general_types.RawMetric) {
	for {
		time.Sleep(time.Duration(rawLifePeriod) * time.Second)

		avgSpeed, lastAvailability := util.GetPreparedMetrics(rawMetrics)

		metrics_storage.AddPreparedMetric(connection, avgSpeed, lastAvailability, func() {
			rawMetrics = make(map[string]general_types.RawMetric)
		})
	}
}
