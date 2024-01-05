package util

import "time"

func ParseTimeBracket(fromTime string, toTime string) (time.Time, time.Time, error) {
	fromTimeParsed, err := time.Parse(time.RFC3339, fromTime)
	if err != nil {
		return time.Now(), time.Now(), err
	}
	toTimeParsed, err := time.Parse(time.RFC3339, toTime)
	if err != nil {
		return time.Now(), time.Now(), err
	}
	return fromTimeParsed, toTimeParsed, nil
}

func GetIsRawTimeBracket(toTimeParsed time.Time, fromTimeParsed time.Time, rawLifePeriod int) bool {
	return time.Now().Sub(toTimeParsed).Seconds() < float64(rawLifePeriod) && fromTimeParsed.Sub(toTimeParsed).Seconds() < float64(rawLifePeriod)
}
