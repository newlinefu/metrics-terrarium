package metrics_storage

import (
	"strconv"
	"time"
)

func createGetPreparedMetrics(fromTime string, toTime string) string {
	return `
		select 	MR.record_id as "id",
				MR.record_time as "time",
				availability.availability_value as "availability",
				speed.speed_value as "speed_value"
		  from	metrics_record MR
		  join  availability ON availability.record_id = MR.record_id
		  join  speed ON speed.record_id = MR.record_id
		 where  MR.record_time >= timestamp '` + fromTime + `' and 
			  	MR.record_time <= timestamp '` + toTime + `';
	`
}

func createAddMetricRecord() string {
	return `
		INSERT INTO metrics_record(record_time)
		VALUES ('` + time.Now().Format(time.RFC3339) + `')
		RETURNING record_id
	`
}

func createAddAvailability(record_id int, avgAvailability float32) string {
	return `
		INSERT INTO availability (record_id, availability_value)
		VALUES (` + strconv.Itoa(record_id) + `, ` + strconv.FormatFloat(float64(avgAvailability), 'f', 6, 32) + `)
	`
}

func createAddSpeed(record_id int, avgSpeed float32) string {
	return `
		INSERT INTO speed (record_id, speed_value)
		VALUES (` + strconv.Itoa(record_id) + `, ` + strconv.FormatFloat(float64(avgSpeed), 'f', 6, 32) + `)
	`
}
