package main

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"log"
	"metricsTerrarium/lib"
	"net/http"
	"time"
)

var config *lib.Config
var producer sarama.SyncProducer

type AvailabilityMetricResponse struct {
	Availability bool      `json:"availability"`
	Timestamp    time.Time `json:"timestamp"`
}

type SpeedMetricResponse struct {
	Speed     float32   `json:"speed"`
	Timestamp time.Time `json:"timestamp"`
}

type MetricMessage struct {
	MetricName string `json:"metric"`
	Value      string `json:"value"`
	Timestamp  string `json:"timestamp"`
}

func main() {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error with env variables file definition. Err: %s", err)
	} else {
		log.Printf("ENV variables initialized succesfully")
	}

	config = lib.CreateConfig()

	var kafkaConnectionErr error
	producer, kafkaConnectionErr = sarama.NewSyncProducer([]string{config.KafkaAddress}, nil)

	if kafkaConnectionErr != nil {
		log.Fatalf("Failed to create producer: %v", kafkaConnectionErr)
	}

	defer producer.Close()

	logFile, err := lib.SetupLogs()
	defer logFile.Close()
	if err != nil {
		log.Fatalf("Error with logs initialization occured. Err: %s", err)
	}

	startHardwareCheck()
}

func startHardwareCheck() {
	for {
		log.Printf("Start of hardware check")

		speedResponse, err := http.Get("http://localhost" + config.HardwareCircuitBreakerPort + "/api/v1/trigger-speed-metric")
		if err != nil {
			log.Printf("Request for speed metric response failed. Err: %s", err)
		} else {
			var speedBody *SpeedMetricResponse
			err = json.NewDecoder(speedResponse.Body).Decode(&speedBody)

			if err != nil {
				log.Printf("Error during speed metric decoding. Err: %s", err)
			} else {
				log.Printf("Speed metric recieved. Value is: %f | Time is: %s", speedBody.Speed, speedBody.Timestamp)

				requestID := uuid.New().String()

				bytes, err := json.Marshal(MetricMessage{
					MetricName: "speed",
					Value:      fmt.Sprintf("%f", speedBody.Speed),
					Timestamp:  speedBody.Timestamp.Format(time.RFC3339),
				})

				msg := &sarama.ProducerMessage{
					Topic: "metrics",
					Key:   sarama.StringEncoder(requestID),
					Value: sarama.ByteEncoder(bytes),
				}

				// отправка сообщения в Kafka
				partition, offset, err := producer.SendMessage(msg)
				log.Printf("Message send with partition: %b | offset is: %b", partition, offset)
				if err != nil {
					log.Printf("Failed to send message to Kafka: %s", err)
					return
				}
			}
		}

		availabilityResponse, err := http.Get("http://localhost" + config.HardwareCircuitBreakerPort + "/api/v1/trigger-availability-metric")
		if err != nil {
			log.Printf("Request for availability metric response failed. Err: %s", err)
		} else {
			var availabilityBody *AvailabilityMetricResponse
			err = json.NewDecoder(availabilityResponse.Body).Decode(&availabilityBody)

			if err != nil {
				log.Printf("Error during speed metric decoding. Err: %s", err)
			} else {
				log.Printf("Availability metric recieved. Value is: %v | Timestamp is %s", availabilityBody.Availability, availabilityBody.Timestamp)

				requestID := uuid.New().String()

				bytes, err := json.Marshal(MetricMessage{
					MetricName: "availability",
					Value:      fmt.Sprintf("%t", availabilityBody.Availability),
					Timestamp:  availabilityBody.Timestamp.Format(time.RFC3339),
				})

				msg := &sarama.ProducerMessage{
					Topic: "metrics",
					Key:   sarama.StringEncoder(requestID),
					Value: sarama.ByteEncoder(bytes),
				}

				partition, offset, err := producer.SendMessage(msg)
				log.Printf("Message send with partition: %b | offset is: %b", partition, offset)
				if err != nil {
					log.Printf("Failed to send message to Kafka: %s", err)
					return
				}
			}
		}

		time.Sleep(time.Duration(config.HardwareCheckInterval) * time.Second)
	}
}
