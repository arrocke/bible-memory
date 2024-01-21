package main

import (
	"fmt"
	"regexp"
	"strconv"
)

func FormatReference(book string, startChapter int32, startVerse int32, endChapter int32, endVerse int32) string {
	return fmt.Sprintf("%s %d:%d-%d:%d", book, startChapter, startVerse, endChapter, endVerse)
}

type ParsedReference struct {
	Book         string
	StartChapter int32
	StartVerse   int32
	EndChapter   int32
	EndVerse     int32
}

var referenceRegexp = regexp.MustCompile(`(.+?)\s*(\d+)[.:](\d+)\s*-\s*(\d+)[.:](\d+)`)

func ParseReference(referenceString string) (ParsedReference, error) {
	match := referenceRegexp.FindStringSubmatch(referenceString)
	if match == nil {
		return ParsedReference{}, nil
	}

	startChapter, _ := strconv.ParseInt(match[2], 10, 32)
	startVerse, _ := strconv.ParseInt(match[3], 10, 32)
	endChapter, _ := strconv.ParseInt(match[4], 10, 32)
	endVerse, _ := strconv.ParseInt(match[5], 10, 32)

	return ParsedReference{
			Book:         match[1],
			StartChapter: int32(startChapter),
			StartVerse:   int32(startVerse),
			EndChapter:   int32(endChapter),
			EndVerse:     int32(endVerse),
		},
		nil
}
