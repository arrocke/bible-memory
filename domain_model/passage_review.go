package domain_model

import (
	"math"
	"strconv"
	"time"
)

func ParseReviewInterval(str string) (uint, error) {
	parsedInterval, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(parsedInterval), nil
}

type NextReviewDate struct {
	Value time.Time
}

func ParseNextReviewDate(str string) (time.Time, error) {
	date, err := time.Parse("2006-01-02", str)
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}

type PassageReview struct {
	Interval   uint
	ReviewedAt *time.Time
	NextReview time.Time
}

func (r *PassageReview) Update(interval uint, nextReview time.Time) PassageReview {
    var reviewedAt *time.Time
    if r != nil {
        reviewedAt = r.ReviewedAt
    }

	return PassageReview{
		Interval:   interval,
		ReviewedAt: reviewedAt,
		NextReview: nextReview,
	}
}

func (r *PassageReview) nextInterval(grade uint, timestamp time.Time) (int, error) {
	if r == nil || r.ReviewedAt == nil {
		switch grade {
		case 1:
			return 1, nil
		case 2:
			return 1, nil
		case 3:
			return 2, nil
		case 4:
			return 4, nil
		}
	} else {
		reviewInterval := math.Ceil(timestamp.Sub(*r.ReviewedAt).Hours() / 24.0)
		extension := min(float64(r.Interval), reviewInterval) + 0.5*max(0.0, reviewInterval-float64(r.Interval.Value))

		switch grade {
		case 1:
			return max(1, int(extension/2.0)), nil
		case 2:
			return int(extension), nil
		case 3:
			return int(float64(r.Interval) + extension), nil
		case 4:
			return int(float64(r.Interval) + 1.5*extension), nil
		}
	}

	return 0, nil
}

func (r *PassageReview) Review(grade uint, timestamp time.Time) (PassageReview, error) {
    // Only review once a day.
    if r != nil && r.ReviewedAt != nil && r.ReviewedAt.Equal(timestamp) {
        return *r, nil
    }
    
    nextInterval, err := r.nextInterval(grade, timestamp)
    if err != nil {
        return PassageReview{}, err
    }
     
	nextReview := timestamp.AddDate(0, 0, nextInterval)

    return PassageReview {
        Interval: uint(nextInterval),
        ReviewedAt: &timestamp,
        NextReview: nextReview,
    }, nil
}
