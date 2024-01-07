package metrics_storage

import (
	"context"
	"database/sql"
	"log"
	"metricsTerrarium/services/metrics-manager/pkg/api"
	"strconv"
)

func GetPreparedMetrics(req *api.PreparedMetricsRequestMessage, dbConnection *sql.DB) (float32, float32, error) {
	rows, err := dbConnection.Query(createGetPreparedMetrics(req.FromTime, req.ToTime))
	defer rows.Close()
	if err != nil {
		log.Printf("Error while execution query for prepared metrics: %s", err)
		return 0, 0, rows.Err()
	}
	var (
		avgAvailability float32
		avgSpeed        float32
		counter         int
	)
	var (
		id           string
		time         string
		availability float32
		speed_value  float32
	)
	for rows.Next() {
		counter++
		err := rows.Scan(&id, &time, &availability, &speed_value)
		if err != nil {
			log.Printf("Error during row scan: %s", err)
			return 0, 0, rows.Err()
		}
		avgAvailability += availability
		avgSpeed += speed_value
	}
	if rows.Err() != nil {
		log.Printf("Error during all rows reading: %s", rows.Err())
		return 0, 0, rows.Err()
	}

	if counter == 0 {
		return 0, 0, nil
	}
	return avgSpeed / float32(counter), avgAvailability / float32(counter), err
}

func AddPreparedMetric(dbConnection *sql.DB, avgSpeed float32, avgAvailability float32) {

	log.Printf("Starting to commit cache data to database")
	var query = createAddMetricRecord()
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
		_, err = tx.ExecContext(ctx, createAddAvailability(record_id, avgAvailability))
		if err != nil {
			tx.Rollback()
			log.Printf("Failed to begin context of transaction [Availability metric]. Err: %s", err)
		}
		_, err = tx.ExecContext(ctx, createAddSpeed(record_id, avgSpeed))
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
		}
	}
}
