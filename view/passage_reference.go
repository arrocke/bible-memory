package view

import (
	"fmt"
	"regexp"
	"strconv"
)

type PassageReference struct {
	Book         string
	StartChapter int
	StartVerse   int
	EndChapter   int
	EndVerse     int
}

func (r PassageReference) String() string {
	if r.StartChapter == r.EndChapter && r.StartVerse == r.EndVerse {
		return fmt.Sprintf("%s %d:%d", r.Book, r.StartChapter, r.StartVerse)
	} else if r.StartChapter == r.EndChapter {
		return fmt.Sprintf("%s %d:%d-%d", r.Book, r.StartChapter, r.StartVerse, r.EndVerse)
	} else {
		return fmt.Sprintf("%s %d:%d-%d:%d", r.Book, r.StartChapter, r.StartVerse, r.EndChapter, r.EndVerse)
	}
}

var referenceRegexp = regexp.MustCompile(`(.+?)\s*(\d+)[.:](\d+)(?:\s*-\s*(?:(\d+)[.:])?(\d+))?`)

func ParsePassageReference(str string) (PassageReference, error) {
	match := referenceRegexp.FindStringSubmatch(str)
	if match == nil {
		return PassageReference{}, fmt.Errorf("Failed to parse reference")
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

	return PassageReference{
			Book:         match[1],
			StartChapter: int(startChapter),
			StartVerse:   int(startVerse),
			EndChapter:   int(endChapter),
			EndVerse:     int(endVerse),
		},
		nil
}
