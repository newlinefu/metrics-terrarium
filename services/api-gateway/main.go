package main

import (
	"context"
	"encoding/json"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"metricsTerrarium/lib"
	"metricsTerrarium/services/api-gateway/pkg/api"
	"metricsTerrarium/services/api-gateway/pkg/util"
	"net/http"
)

type TimeBracket struct {
	FromTime string `json:"fromTime"`
	ToTime   string `json:"toTime"`
}

type MetricsResponse struct {
	Availability float32 `json:"availability"`
	Speed        float32 `json:"speed"`
	Type         string  `json:"type"`
}

var client api.MetricsGetterClient

var config *lib.Config

func main() {
	logFile, err := lib.SetupLogs()
	defer logFile.Close()

	if err != nil {
		log.Fatalf("Error with logs initialization occured. Err: %s", err)
	}

	err = godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error with env variables file definition. Err: %s", err)
	}
	config = lib.CreateConfig()
	conn, err := grpc.Dial("localhost"+config.MetricsManagerPort, grpc.WithTransportCredentials(insecure.NewCredentials()))

	log.Printf("GRPC dial initialized with MetricsManagerPort=[ %s].", config.MetricsManagerPort)
	if err != nil {
		log.Fatalf("Error with GRPC dial initialization. Err: %s", err)
	}
	client = api.NewMetricsGetterClient(conn)

	http.HandleFunc("/api/v1/get-metrics", handleGetMetrics)

	err = http.ListenAndServe(config.ApiGatewayHTTPPort, nil)
	if err != nil {
		log.Fatalf("Failed to listrening server at port %s", config.ApiGatewayHTTPPort)
	}
}

func handleGetMetrics(w http.ResponseWriter, r *http.Request) {

	log.Printf("\"/get-metrics\" %s request handled", r.Method)

	var tb TimeBracket
	err := json.NewDecoder(r.Body).Decode(&tb)
	if err != nil {
		log.Printf("Error with body parsing. Err: %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fromTimeParsed, toTimeParsed, err := util.ParseTimeBracket(tb.FromTime, tb.ToTime)
	if err != nil {
		log.Printf("Error with dates parsing in body. Body: %s. Err: %s\n", tb, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isRawTimeBracket := util.GetIsRawTimeBracket(fromTimeParsed, toTimeParsed, config.RawLifePeriod)

	var metrics *api.MetricsResponse

	if isRawTimeBracket {
		log.Println("Raw case handled")
		metrics, err = client.GetRawMetrics(context.Background(), &api.RawMetricsRequestMessage{})
		if err != nil {
			log.Printf("Error with GRPC request. Err: %s\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		log.Println("Prepared case handled")
		metrics, err = client.GetPreparedMetrics(context.Background(), &api.PreparedMetricsRequestMessage{
			FromTime: tb.FromTime,
			ToTime:   tb.ToTime,
		})
		if err != nil {
			log.Printf("Error with GRPC request. Err: %s\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	log.Printf("Response created. Availability: %b | Speed: %f\n", metrics.Availability, metrics.Speed)
	err = json.NewEncoder(w).Encode(MetricsResponse{
		Availability: metrics.Availability,
		Speed:        metrics.Speed,
		Type:         metrics.Type,
	})
	if err != nil {
		log.Printf("Error sending response. Err: %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
