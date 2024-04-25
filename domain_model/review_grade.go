package domain_model

import "fmt"

type ReviewGrade int

const GRADE_FAIL = ReviewGrade(1)
const GRADE_HARD = ReviewGrade(2)
const GRADE_OK = ReviewGrade(3)
const GRADE_EASY = ReviewGrade(4)

func NewReviewGrade(grade int) (ReviewGrade, error) {
    if grade <= 0 || grade > 4 {
        return 0, fmt.Errorf("Invalid review grade: %v", grade)
    }
	return ReviewGrade(grade), nil
}

