package domain_model

import (
	"math"
	"time"
)

type ReviewTimestamp time.Time

func NewReviewTimestamp(timestamp time.Time) ReviewTimestamp {
    return ReviewTimestamp(timestamp)
}

func NewReviewTimestampForToday(tz int) ReviewTimestamp {
	location := time.FixedZone("Temp", tz*60)
	now := time.Now().In(location)
	return ReviewTimestamp(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC))
}

func (i ReviewTimestamp) Value() time.Time {
	return time.Time(i)
}
func (t ReviewTimestamp) Equal(other ReviewTimestamp) bool {
	return t.Value().Equal(other.Value())
}
func (t ReviewTimestamp) IntervalToDate(other ReviewTimestamp) ReviewInterval {
	return ReviewInterval(math.Ceil(t.Value().Sub(other.Value()).Abs().Hours() / 24.0))
}
func (t ReviewTimestamp) AddInterval(interval ReviewInterval) ReviewTimestamp {
	return ReviewTimestamp(t.Value().AddDate(0, 0, interval.Value()))
}

