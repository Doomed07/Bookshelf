package library

import "sync"

type Library struct {
	books map[string]Book
	mtx   sync.RWMutex
}

func NewLibrary() *Library {
	return &Library{
		books: make(map[string]Book),
	}
}

func (l *Library) AddBook(book Book) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if _, ok := l.books[book.Title]; ok {
		return ErrBookAlreadyExists
	}

	l.books[book.Title] = book

	return nil
}

func (l *Library) ReadBook(title string) (Book, error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	book, ok := l.books[title]
	if !ok {
		return Book{}, ErrBookNotFound
	}

	book.ReadIt()

	l.books[book.Title] = book

	return book, nil
}

func (l *Library) UnReadBook(title string) (Book, error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	book, ok := l.books[title]
	if !ok {
		return Book{}, ErrBookNotFound
	}

	book.UnRead()

	l.books[book.Title] = book

	return book, nil
}

func (l *Library) GetBook(title string) (Book, error) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	book, ok := l.books[title]
	if !ok {
		return Book{}, ErrBookNotFound
	}

	return book, nil
}

func (l *Library) AllBooks() map[string]Book {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	temp := make(map[string]Book, len(l.books))

	for k, v := range l.books {
		temp[k] = v
	}

	return temp
}

func (l *Library) AllReadBooks() map[string]Book {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	readBooks := make(map[string]Book, len(l.books))

	for title, books := range l.books {
		if books.Read {
			readBooks[title] = books
		}
	}

	return readBooks
}

func (l *Library) AllUnReadBooks() map[string]Book {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	readBooks := make(map[string]Book, len(l.books))

	for title, books := range l.books {
		if !books.Read {
			readBooks[title] = books
		}
	}

	return readBooks
}

func (l *Library) AuthorBooks(author string) map[string]Book {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	authorBooks := make(map[string]Book, len(l.books))

	for title, book := range l.books {
		if book.Author == author {
			authorBooks[title] = book
		}
	}

	return authorBooks
}

func (l *Library) DeleteBook(title string) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	_, ok := l.books[title]
	if !ok {
		return ErrBookNotFound
	}

	delete(l.books, title)

	return nil
}
