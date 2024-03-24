package domain_model

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

type ReviewGrade int

const GRADE_FAIL = ReviewGrade(1)
const GRADE_HARD = ReviewGrade(2)
const GRADE_OK = ReviewGrade(3)
const GRADE_EASY = ReviewGrade(4)

type InvalidReviewGradeError struct {
	Grade string
}

func (e InvalidReviewGradeError) Error() string {
	return fmt.Sprintf("Invalid review grade: %v", e.Grade)
}

func ParseReviewGrade(gradestr string) (ReviewGrade, error) {
    var grade int
    switch gradestr {
    case "1":
        grade = 1
    case "2":
        grade = 2
    case "3":
        grade = 3
    case "4":
        grade = 4
    default:
		return 0, InvalidReviewGradeError{Grade: gradestr}
    }

	return ReviewGrade(grade), nil
}

type ReviewInterval int

func (i ReviewInterval) Value() int {
    return int(i)
}

type InvalidReviewIntervalError struct {
	Interval string
}

func (e InvalidReviewIntervalError) Error() string {
	return fmt.Sprintf("Invalid review interval: %v", e.Interval)
}

func ParseReviewInterval(intervalstr string) (ReviewInterval, error) {
    interval, err := strconv.ParseUint(intervalstr, 10, 64)
	if err != nil {
		return 0, InvalidReviewIntervalError{Interval: intervalstr}
	}
	return ReviewInterval(interval), nil
}

type ReviewTimestamp time.Time

type InvalidReviewTimestampError struct {
	Timestamp string
}

func (e InvalidReviewTimestampError) Error() string {
	return fmt.Sprintf("Invalid review timestamp: %v", e.Timestamp)
}

func ParseReviewTimestamp(timestampstr string, format string) (ReviewTimestamp, error) {
    timestamp, err := time.Parse(format, timestampstr)
	if err != nil {
		return ReviewTimestamp{}, InvalidReviewIntervalError{Interval: timestampstr}
	}
	return ReviewTimestamp(timestamp), nil
}

func (i ReviewTimestamp) Value() time.Time {
    return time.Time(i)
}
func (t ReviewTimestamp) DifferenceInDays(other ReviewTimestamp) float64 {
	return math.Ceil(math.Abs(t.Value().Sub(other.Value()).Hours() / 24.0))
}
func (t ReviewTimestamp) Equal(other ReviewTimestamp) bool {
	return t.Value().Equal(other.Value())
}
func (t ReviewTimestamp) AddDays(days int) ReviewTimestamp {
	return ReviewTimestamp(t.Value().AddDate(0, 0, days))
}

type PassageReview struct {
	Interval   ReviewInterval
	NextReview ReviewTimestamp
	ReviewedAt *ReviewTimestamp
}

func (r *PassageReview) Update(interval ReviewInterval, nextReview ReviewTimestamp) PassageReview {
	var reviewedAt *ReviewTimestamp
	if r != nil {
		reviewedAt = r.ReviewedAt
	}

	return PassageReview{
		Interval:   interval,
		ReviewedAt: reviewedAt,
		NextReview: nextReview,
	}
}

func (r *PassageReview) nextInterval(grade ReviewGrade, timestamp ReviewTimestamp) ReviewInterval {
	if r == nil {
		switch grade {
		case GRADE_FAIL:
			return 1
		case GRADE_HARD:
			return 1
		case GRADE_OK:
			return 2
		case GRADE_EASY:
			return 4
		}
	} else {
		actualInterval := float64(r.Interval)
		if r.ReviewedAt == nil {
			actualInterval = timestamp.DifferenceInDays(*r.ReviewedAt)
		}
		expectedInterval := float64(r.Interval)
		extraInterval := max(0.0, actualInterval-expectedInterval)
		weightedInterval := min(expectedInterval, actualInterval) + 0.5*extraInterval

		switch grade {
		case GRADE_FAIL:
			return ReviewInterval(max(1, int(weightedInterval/2.0)))
		case GRADE_HARD:
			return ReviewInterval(weightedInterval)
		case GRADE_OK:
			return ReviewInterval(expectedInterval + weightedInterval)
		case GRADE_EASY:
			return ReviewInterval(expectedInterval + 1.5*weightedInterval)
		}
	}

	panic("unreachable code in nextInterval")
}

func (r *PassageReview) Review(grade ReviewGrade, timestamp ReviewTimestamp) PassageReview {
	// Only review once a day.
	if r != nil && r.ReviewedAt != nil && r.ReviewedAt.Equal(timestamp) {
		return *r
	}

	nextInterval := r.nextInterval(grade, timestamp)

	return PassageReview{
		Interval:   nextInterval,
		ReviewedAt: &timestamp,
		NextReview: timestamp.AddDays(nextInterval.Value()),
	}
}
