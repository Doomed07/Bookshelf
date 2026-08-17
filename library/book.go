package library

import "time"

type Book struct {
	Title   string
	Author  string
	Sheets  int
	Read    bool
	AddedAt time.Time
	ReadAt  *time.Time
}

func NewBook(title string, author string, sheets int) Book {
	return Book{
		Title:   title,
		Author:  author,
		Sheets:  sheets,
		Read:    false,
		AddedAt: time.Now(),
		ReadAt:  nil,
	}
}

func (b *Book) ReadIt() {
	b.Read = true
	time := time.Now()
	b.ReadAt = &time
}

func (b *Book) UnRead() {
	b.Read = false
	b.ReadAt = nil
}
