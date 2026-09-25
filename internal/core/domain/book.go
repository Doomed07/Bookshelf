package domain

type Book struct {
	ID          int
	Title       string
	Author      string
	Year        int
	Pages       int
	Genres      []string
	Description string
	Score       *int
	ReadsCount  int
}

func NewBook(
	id int,
	title string,
	author string,
	year int,
	pages int,
	genres []string,
	description string,
	score *int,
	readsCount int,
) Book {
	return Book{
		ID:          id,
		Title:       title,
		Author:      author,
		Year:        year,
		Pages:       pages,
		Genres:      genres,
		Description: description,
		Score:       score,
		ReadsCount:  readsCount,
	}
}
