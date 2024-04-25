package domain_model

import "fmt"

type PassageReference struct {
	Book         string
	StartChapter int
	StartVerse   int
	EndChapter   int
	EndVerse     int
}

func NewPassageReference(book string, startChapter int, startVerse int, endChapter int, endVerse int) (PassageReference, error) {
	if startChapter > endChapter {
		return PassageReference{}, fmt.Errorf("end chapter must be on or after start chapter")
	} else if startVerse > endVerse && startChapter == endChapter {
		return PassageReference{}, fmt.Errorf("end verse must be on or after start verse")
	}

	return PassageReference{
		Book:         book,
		StartChapter: startChapter,
		StartVerse:   startVerse,
		EndChapter:   endChapter,
		EndVerse:     endVerse,
	}, nil
}
