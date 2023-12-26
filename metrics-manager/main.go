package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"metricsTerrarium/metrics-manager/pkg/api"
	"metricsTerrarium/metrics-manager/pkg/types"
	"metricsTerrarium/metrics-manager/pkg/util"
	"net"
	"strconv"
	"time"
)

var dbConnection *sql.DB
var config *types.Config

type RawMetric struct {
	MetricName string `json:"metric"`
	Value      string `json:"value"`
}

var rawMetricsCache map[string]RawMetric

type metricsGetterServerImpl struct {
	api.UnimplementedMetricsGetterServer
}

func GetPreparedMetrics() (float32, bool) {
	var avgSpeed float32 = 0
	var lastAvailability bool
	index := 0
	for _, metric := range rawMetricsCache {
		index++
		if metric.MetricName == "availability" {
			if (len(rawMetricsCache) - 1) == index {
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
	if len(rawMetricsCache) == 0 {
		return 0, lastAvailability
	}
	avgSpeed = avgSpeed / float32(len(rawMetricsCache))

	return avgSpeed, lastAvailability
}

func (s metricsGetterServerImpl) GetRawMetrics(context.Context, *api.RawMetricsRequestMessage) (*api.MetricsResponse, error) {
	log.Printf("Get Raw metrics request handled successfuly")

	avgSpeed, lastAvailability := GetPreparedMetrics()
	availabilityResponse := 0.0
	if lastAvailability {
		availabilityResponse = 1.0
	}
	return &api.MetricsResponse{
		Availability: float32(availabilityResponse),
		Speed:        avgSpeed,
	}, nil
}

func (s metricsGetterServerImpl) GetPreparedMetrics(ctx context.Context, req *api.PreparedMetricsRequestMessage) (*api.MetricsResponse, error) {
	log.Printf("Get Prepared metrics request handled successfuly")

	rows, err := dbConnection.Query(`
		select 	MR.record_id as "id",
				MR.record_time as "time",
				availability.availability_value as "availability",
				speed.speed_value as "speed_value"
		  from	metrics_record MR
		  join  availability ON availability.record_id = MR.record_id
		  join  speed ON speed.record_id = MR.record_id
		 where  MR.record_time >= timestamp '` + req.FromTime + `' and 
			  	MR.record_time <= timestamp '` + req.ToTime + `';
	`)
	defer rows.Close()
	if err != nil {
		log.Printf("Error while execution query for prepared metrics: %s", err)
	}
	var (
		avgAvailability float32
		avgSpeed        float32
		counter         int
	)
	var (
		id           string
		time         string
		availability bool
		speed_value  float32
	)
	for rows.Next() {
		counter++
		err := rows.Scan(&id, &time, &availability, &speed_value)
		if err != nil {
			log.Printf("Error during row scan: %s", err)
		}
		if availability {
			avgAvailability += 1.0
		}
		avgSpeed += speed_value
	}
	if rows.Err() != nil {
		log.Printf("Error during all rows reading: %s", rows.Err())
	}
	return &api.MetricsResponse{
		Availability: avgAvailability / float32(counter),
		Speed:        avgSpeed / float32(counter),
	}, nil
}

func createDBConnect() error {
	connStr := "user=" + config.MetricsDataBaseUser + " password=" + config.MetricsDataBasePassword + " dbname=" + config.MetricsDataBase + " sslmode=disable"
	var err error
	dbConnection, err = sql.Open("postgres", connStr)
	return err
}

func main() {
	logFile, err := util.SetupLogs()
	defer logFile.Close()
	if err != nil {
		log.Fatalf("Error with logs initialization occured. Err: %s", err)
	}

	rawMetricsCache = map[string]RawMetric{}

	err = godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error with env variables file definition. Err: %s", err)
	} else {
		log.Printf("ENV variables initialized succesfully")
	}
	config = types.CreateConfig()

	consumer, err := sarama.NewConsumer([]string{config.KafkaAddress}, nil)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	} else {
		log.Printf("Consumer created succesfully")
	}
	defer consumer.Close()
	partConsumer, err := consumer.ConsumePartition("metrics", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}
	defer partConsumer.Close()

	go func() {
		for {
			select {
			case msg, ok := <-partConsumer.Messages():
				if !ok {
					log.Println("Channel closed, exiting goroutine")
					return
				}
				msgKey := string(msg.Key[:])
				_, exists := rawMetricsCache[msgKey]
				if exists {
					delete(rawMetricsCache, msgKey)
					log.Printf("Key %s in metrics cache map was replaced. Already exists", msgKey)
				}
				var rawMetric RawMetric

				err := json.Unmarshal(msg.Value, &rawMetric)
				if err != nil {
					log.Printf("Unmarshaling error. Error: %s", err)
				}

				rawMetricsCache[msgKey] = rawMetric
				log.Printf("Metric { type: %s, value: %s } was inserted to cache with key: %s", rawMetric.MetricName, rawMetric.Value, msgKey)
			}
		}
	}()

	go func() {
		for {
			time.Sleep(time.Duration(config.RawLifePeriod) * time.Second)

			avgSpeed, lastAvailability := GetPreparedMetrics()

			log.Printf("Starting to commit cache data to database")
			var query = `
				INSERT INTO metrics_record(record_time)
				VALUES ('` + time.Now().Format(time.RFC3339) + `')
				RETURNING record_id
			`
			var record_id int
			err := dbConnection.QueryRow(
				query,
			).Scan(&record_id)
			if err != nil {
				log.Printf("Failed to insert metrics_record record. Err: %s", err)
			} else {

				log.Printf("metrics_record table Row inserted succesfully")

				log.Printf(strconv.Itoa(record_id))

				ctx := context.Background()
				tx, err := dbConnection.BeginTx(ctx, nil)
				if err != nil {
					tx.Rollback()
					log.Printf("Failed to begin context of transaction. Err: %s", err)
				}
				_, err = tx.ExecContext(ctx, `
					INSERT INTO availability (record_id, availability_value)
					VALUES (`+strconv.Itoa(record_id)+`, `+strconv.FormatBool(lastAvailability)+`)
				`)
				if err != nil {
					tx.Rollback()
					log.Printf("Failed to begin context of transaction [Availability metric]. Err: %s", err)
				}
				_, err = tx.ExecContext(ctx, `
					INSERT INTO speed (record_id, speed_value)
					VALUES (`+strconv.Itoa(record_id)+`, `+strconv.FormatFloat(float64(avgSpeed), 'f', 6, 32)+`)
				`)
				if err != nil {
					tx.Rollback()
					log.Printf("2 Failed to begin context of transaction [Speed metric]. Err: %s", err)
				}
				err = tx.Commit()
				if err != nil {
					tx.Rollback()
					log.Fatal(err)
				} else {
					log.Printf("Rows inserted succesfully")
					rawMetricsCache = make(map[string]RawMetric)
				}
			}
		}
	}()

	err = createDBConnect()
	if err != nil {
		log.Fatalf("Error during db connection creation. Err: %s", err)
	} else {
		log.Printf("Connection to database created succesfully")
	}
	defer dbConnection.Close()

	lis, err := net.Listen("tcp", config.MetricsManagerPort)
	if err != nil {
		log.Fatalf("TCP Connection creation error. Err: %s", err)
	} else {
		log.Printf("Created GRPC listener at %s", config.MetricsManagerPort)
	}
	s := grpc.NewServer()
	service := &metricsGetterServerImpl{}
	api.RegisterMetricsGetterServer(s, service)
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve. Err: %s", err)
	}
}
