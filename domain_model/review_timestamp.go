package domain_model

import (
	"fmt"
	"math"
	"time"
)

type ReviewTimestamp time.Time

func NewReviewTimestamp(timestamp time.Time) ReviewTimestamp {
    return ReviewTimestamp(timestamp)
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
func (t ReviewTimestamp) AfterInterval(interval ReviewInterval) ReviewTimestamp {
	return ReviewTimestamp(t.Value().AddDate(0, 0, interval.Value()))
}


