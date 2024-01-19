package main

import "fmt"

func FormatReference(book string, startChapter int32, startVerse int32, endChapter int32, endVerse int32) string {
	return fmt.Sprintf("%v %v:%v-%v:%v", book, startChapter, startVerse, endChapter, endVerse)
}
