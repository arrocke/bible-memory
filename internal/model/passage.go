package model

import (
	"fmt"
	"time"
)

type Passage struct {
    Id int
    Book         string
	StartChapter int
	StartVerse   int
	EndChapter   int
	EndVerse     int
	Text         string
    Owner       int `db:"user_id"`
	Interval     *int
	ReviewedAt   *time.Time
    NextReview   *time.Time `db:"review_at"`
}

func (p Passage) Reference() string {
    if p.StartChapter == p.EndChapter && p.StartVerse == p.EndVerse {
		return fmt.Sprintf("%s %d:%d", p.Book, p.StartChapter, p.StartVerse)
	} else if p.StartChapter == p.EndChapter {
		return fmt.Sprintf("%s %d:%d-%d", p.Book, p.StartChapter, p.StartVerse, p.EndVerse)
	} else {
		return fmt.Sprintf("%s %d:%d-%d:%d", p.Book, p.StartChapter, p.StartVerse, p.EndChapter, p.EndVerse)
	}
}
