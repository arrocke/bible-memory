package domain_model

import "fmt"

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

