package util

import (
	"metricsTerrarium/services/metrics-manager/internal/general_types"
	"strconv"
)

func GetPreparedMetrics(rawMetrics map[string]general_types.RawMetric) (float32, bool) {
	var avgSpeed float32 = 0
	var lastAvailability bool
	index := 0
	for _, metric := range rawMetrics {
		index++
		if metric.MetricName == "availability" {
			if (len(rawMetrics) - 1) == index {
				var err error
				lastAvailability, err = strconv.ParseBool(metric.Value)
				if err != nil {
					lastAvailability = false
				}
			}
		}

		if metric.MetricName == "speed" {
			actualSpeed, err := strconv.ParseFloat(metric.Value, 32)
			if err == nil {
				avgSpeed += float32(actualSpeed)
			}
		}
	}
	if len(rawMetrics) == 0 {
		return 0, lastAvailability
	}
	avgSpeed = avgSpeed / float32(len(rawMetrics))

	return avgSpeed, lastAvailability
}
