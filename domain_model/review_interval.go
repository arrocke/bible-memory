package domain_model

import (
	"fmt"
)

type ReviewInterval int

func (i ReviewInterval) Value() int {
	return int(i)
}

type InvalidReviewIntervalError struct {
	Interval int
}

func (e InvalidReviewIntervalError) Error() string {
	return fmt.Sprintf("Invalid review interval: %v", e.Interval)
}

func NewReviewInterval(interval int) (ReviewInterval, error) {
    if interval <= 0 || interval > 4 {
        return 0, InvalidReviewIntervalError{ interval }
    }
    return ReviewInterval(interval), nil
}

func FirstInterval(grade ReviewGrade) ReviewInterval {
    switch grade {
    case GRADE_FAIL:
        return 1
    case GRADE_HARD:
        return 1
    case GRADE_OK:
        return 2
    case GRADE_EASY:
        return 4
    default:
        panic("this should be unreachable")
    }
}

func (expected ReviewInterval) NextAfterReview(grade ReviewGrade, actual ReviewInterval) ReviewInterval {
    extra := max(0, actual - expected)
    weighted := min(expected, actual) + extra / 2

    switch grade {
    case GRADE_FAIL:
        return max(1, weighted / 2)
    case GRADE_HARD:
        return weighted
    case GRADE_OK:
        return expected + weighted
    case GRADE_EASY:
        return expected + 3 * weighted / 2
    default:
        panic("this should be unreachable")
    }
}
