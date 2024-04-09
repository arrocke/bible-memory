package view

import "fmt"

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
