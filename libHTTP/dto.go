package libhttp

import (
	"encoding/json"
	"errors"
	"time"
)

type BookDTO struct {
	Title  string
	Author string
	Sheets int
}

func (b *BookDTO) ValidateforNewBook() error {
	if b.Title == "" {
		return errors.New("title is empty")
	}
	if b.Author == "" {
		return errors.New("author is empty")
	}
	if b.Sheets == 0 {
		return errors.New("sheets is empty")
	}
	return nil
}

type ReadBookDTO struct {
	Read bool
}

type ErrorDTO struct {
	Message string
	Time    time.Time
}

func (e *ErrorDTO) ToString() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
