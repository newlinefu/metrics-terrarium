package job

import (
	"database/sql"
	"log"
	"metricsTerrarium/services/metrics-manager/internal/general_types"
	"metricsTerrarium/services/metrics-manager/internal/storage/metrics_storage"
	"metricsTerrarium/services/metrics-manager/internal/util"
	"time"
)

func PrepareMetrics(connection *sql.DB, rawLifePeriod int, rawMetrics map[string]general_types.RawMetric) {
	var prevTime = time.Time{}
	for {
		time.Sleep(time.Duration(rawLifePeriod) * time.Second)

		actualTime := time.Now()
		log.Printf("Metrics preparation started. Length of raw metrics is: %v", len(rawMetrics))
		avgSpeed, avgAvailability := util.GetPreparedMetrics(rawMetrics, prevTime, actualTime, rawLifePeriod)

		metrics_storage.AddPreparedMetric(connection, avgSpeed, avgAvailability)

		util.ClearPreparedMetricsFromCache(rawMetrics, prevTime, actualTime, rawLifePeriod)

		prevTime = actualTime
	}
}
