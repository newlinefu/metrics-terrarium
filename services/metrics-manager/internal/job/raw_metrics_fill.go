package job

import (
	"encoding/json"
	"log"
	"metricsTerrarium/services/metrics-manager/internal/general_types"
	"metricsTerrarium/services/metrics-manager/internal/transport/kafka_listener"
)

func ListenIncomingMetrics(kafka kafka_listener.Kafka, rawMetrics map[string]general_types.RawMetric) {
	for {
		select {
		case msg, ok := <-(*kafka.PartConsumer).Messages():
			if !ok {
				log.Println("Channel closed, exiting goroutine")
				return
			}
			msgKey := string(msg.Key[:])
			_, exists := rawMetrics[msgKey]
			if exists {
				delete(rawMetrics, msgKey)
				log.Printf("Key %s in metrics cache map was replaced. Already exists", msgKey)
			}
			var rawMetric general_types.RawMetric

			err := json.Unmarshal(msg.Value, &rawMetric)
			if err != nil {
				log.Printf("Unmarshaling error. Error: %s", err)
			}

			rawMetrics[msgKey] = rawMetric
			log.Printf("Metric { type: %s, value: %s } was inserted to cache with key: %s", rawMetric.MetricName, rawMetric.Value, msgKey)
		}
	}
}
