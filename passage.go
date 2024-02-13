package main

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"
)

func FormatReference(book string, startChapter int32, startVerse int32, endChapter int32, endVerse int32) string {
	if startChapter == endChapter && startVerse == endVerse {
		return fmt.Sprintf("%s %d:%d", book, startChapter, startVerse)
	} else if startChapter == endChapter {
		return fmt.Sprintf("%s %d:%d-%d", book, startChapter, startVerse, endVerse)
	} else {
		return fmt.Sprintf("%s %d:%d-%d:%d", book, startChapter, startVerse, endChapter, endVerse)
	}
}

type ParsedReference struct {
	Book         string
	StartChapter int32
	StartVerse   int32
	EndChapter   int32
	EndVerse     int32
}

var referenceRegexp = regexp.MustCompile(`(.+?)\s*(\d+)[.:](\d+)(?:\s*-\s*(?:(\d+)[.:])?(\d+))?`)

func ParseReference(referenceString string) (ParsedReference, error) {
	match := referenceRegexp.FindStringSubmatch(referenceString)
	if match == nil {
		return ParsedReference{}, nil
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

	return ParsedReference{
			Book:         match[1],
			StartChapter: int32(startChapter),
			StartVerse:   int32(startVerse),
			EndChapter:   int32(endChapter),
			EndVerse:     int32(endVerse),
		},
		nil
}

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
			return int(float64(*currentInterval) + 2.0*extension)
		}
	}

	return 1
}
