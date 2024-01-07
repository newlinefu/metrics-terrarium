package util

import (
	"log"
	"metricsTerrarium/services/metrics-manager/internal/general_types"
	"strconv"
	"time"
)

func GetPreparedMetrics(rawMetrics map[string]general_types.RawMetric, fromTime time.Time, toTime time.Time, rawLifePeriod int) (float32, float32) {
	var avgSpeed float32 = 0
	var avgAvailability float32 = 0
	index := 0
	for _, metric := range rawMetrics {

		metricTimestamp := metric.Timestamp
		parsedMetricTimestamp, err := time.Parse(time.RFC3339, metricTimestamp)
		if err != nil {
			log.Printf("ERROR during metric time parsed. Primary value is: [%s] | Err: %s", metricTimestamp, err)
		} else if GetIsRawTimeBracket(parsedMetricTimestamp, toTime, fromTime, rawLifePeriod) {
			index++
			if metric.MetricName == "availability" {
				if (len(rawMetrics) - 1) == index {
					var intervalAvailability bool
					var err error
					intervalAvailability, err = strconv.ParseBool(metric.Value)
					if err != nil {
						intervalAvailability = false
					}
					if intervalAvailability {
						avgAvailability += 1.0
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
	}

	log.Printf("Count of prepared metrics: [%v]", index)

	if len(rawMetrics) == 0 {
		return 0, 0
	}
	avgSpeed = avgSpeed / float32(len(rawMetrics))
	avgAvailability = avgAvailability / float32(len(rawMetrics))

	return avgSpeed, avgAvailability
}

func ClearPreparedMetricsFromCache(rawMetrics map[string]general_types.RawMetric, fromTime time.Time, toTime time.Time, rawLifePeriod int) {
	cleared := 0
	startMetricsCacheLength := len(rawMetrics)
	for key, metric := range rawMetrics {
		metricTimestamp := metric.Timestamp
		parsedMetricTimestamp, err := time.Parse(time.RFC3339, metricTimestamp)
		if err != nil {
			log.Printf("ERROR during metric time parsed. Primary value is: [%s] | Err: %s", metricTimestamp, err)
		} else if GetIsRawTimeBracket(parsedMetricTimestamp, toTime, fromTime, rawLifePeriod) {
			cleared++
			delete(rawMetrics, key)
		}
	}
	log.Printf("Clearing completed: %v -> %v", startMetricsCacheLength, len(rawMetrics))
}

func GetIsRawTimeBracket(metricTime time.Time, toTime time.Time, fromTime time.Time, rawLifePeriod int) bool {
	return metricTime.Sub(toTime).Seconds() <= float64(rawLifePeriod) && fromTime.Sub(metricTime).Seconds() <= float64(rawLifePeriod)
}
