package main

import (
	"math"
	"time"
)

func GetNextInterval(reviewedAt time.Time, grade int, currentInterval *int, lastReview *time.Time) int {
	if currentInterval == nil || lastReview == nil {
		switch grade {
		case 1:
			return 1
		case 2:
			return 1
		case 3:
			return 2
		case 4:
			return 4
		}
	} else {
		reviewInterval := math.Ceil(reviewedAt.Sub(*lastReview).Hours() / 24.0)
		extension := min(float64(*currentInterval), reviewInterval) + 0.5*max(0.0, reviewInterval-float64(*currentInterval))

		switch grade {
		case 1:
			return max(1, int(extension/2.0))
		case 2:
			return int(extension)
		case 3:
			return int(float64(*currentInterval) + extension)
		case 4:
			return int(float64(*currentInterval) + 1.5*extension)
		}
	}

	return 1
}
