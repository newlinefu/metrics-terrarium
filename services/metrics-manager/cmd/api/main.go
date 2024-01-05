package main

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"metricsTerrarium/lib"
	"metricsTerrarium/services/metrics-manager/internal/database"
	"metricsTerrarium/services/metrics-manager/internal/general_types"
	"metricsTerrarium/services/metrics-manager/internal/job"
	"metricsTerrarium/services/metrics-manager/internal/transport/grpc_server"
	"metricsTerrarium/services/metrics-manager/internal/transport/kafka_listener"
)

var rawMetricsCache *map[string]general_types.RawMetric

func main() {

	*rawMetricsCache = map[string]general_types.RawMetric{}

	logFile, err := lib.SetupLogs()
	defer logFile.Close()
	if err != nil {
		log.Fatalf("Error with logs initialization occured. Err: %s", err)
	}
	err = godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error with env variables file definition. Err: %s", err)
	} else {
		log.Printf("ENV variables initialized succesfully")
	}

	config := lib.CreateConfig()

	kafka := kafka_listener.Kafka{}
	kafka = kafka.Start(kafka_listener.KafkaProperties{Config: config})

	go job.ListenIncomingMetrics(&kafka, rawMetricsCache)

	db := database.Db{}
	db = db.Start(database.DbProperties{Config: config})

	go job.PrepareMetrics(db.Connection, config.RawLifePeriod, rawMetricsCache)

	grpcServer := grpc_server.Server{}
	grpcServer.Start(grpc_server.ServerProperties{
		Config:       config,
		DbConnection: &db,
		RawMetrics:   rawMetricsCache,
	})
}
