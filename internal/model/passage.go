package model

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"
)

type Passage struct {
    Id int
    Reference Reference `db:""`
	Text         string
    Owner       int `db:"user_id"`
	Interval     int
	ReviewedAt   time.Time
    NextReview   time.Time `db:"review_at"`
}

func CreatePassage(referenceString string, text string, userId int) (Passage, error) {
    passage := Passage{
        Text: text,
        Owner: userId,
    }

    reference, err := ParseReference(referenceString)
    if err != nil {
        return passage,err
    }
    passage.Reference = reference

    return passage, nil
}

func (p *Passage) NextInterval(grade int, reviewedAt time.Time) int {
    	if p.Interval == 0 || p.ReviewedAt.IsZero() {
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
		reviewInterval := math.Ceil(reviewedAt.Sub(p.ReviewedAt).Hours() / 24.0)
		extension := min(float64(p.Interval), reviewInterval) + 0.5*max(0.0, reviewInterval-float64(p.Interval))

		switch grade {
		case 1:
			return max(1, int(extension/2.0))
		case 2:
			return int(extension)
		case 3:
			return int(float64(p.Interval) + extension)
		case 4:
			return int(float64(p.Interval) + 1.5*extension)
		}
	}

	return 1
}

func (p *Passage) Review(grade int, reviewedAt time.Time) {
    p.Interval = p.NextInterval(grade, reviewedAt) 
    p.NextReview = reviewedAt.AddDate(0, 0, p.Interval)
    p.ReviewedAt = reviewedAt
}


type Reference struct {
    Book         string
	StartChapter int
	StartVerse   int
	EndChapter   int
	EndVerse     int
}

var ReferenceFormat = regexp.MustCompile(`(.+?)\s*(\d+)[.:](\d+)(?:\s*-\s*(?:(\d+)[.:])?(\d+))?`)

func ParseReference(str string) (Reference, error) {
	match := ReferenceFormat.FindStringSubmatch(str)
	if match == nil {
		return Reference{}, fmt.Errorf("Failed to parse reference")
	}

	startChapter, _ := strconv.ParseInt(match[2], 10, 32)
	startVerse, _ := strconv.ParseInt(match[3], 10, 32)

	endChapter := startChapter
	if match[4] != "" {
		endChapter, _ = strconv.ParseInt(match[4], 10, 32)
	}

	endVerse := startVerse
	if match[5] != "" {
		endVerse, _ = strconv.ParseInt(match[5], 10, 32)
	}

	return Reference{
			Book:         match[1],
			StartChapter: int(startChapter),
			StartVerse:   int(startVerse),
			EndChapter:   int(endChapter),
			EndVerse:     int(endVerse),
		},
		nil
}

func (r Reference) String() string {
    if r.StartChapter == r.EndChapter && r.StartVerse == r.EndVerse {
		return fmt.Sprintf("%s %d:%d", r.Book, r.StartChapter, r.StartVerse)
	} else if r.StartChapter == r.EndChapter {
		return fmt.Sprintf("%s %d:%d-%d", r.Book, r.StartChapter, r.StartVerse, r.EndVerse)
	} else {
		return fmt.Sprintf("%s %d:%d-%d:%d", r.Book, r.StartChapter, r.StartVerse, r.EndChapter, r.EndVerse)
	}
}
