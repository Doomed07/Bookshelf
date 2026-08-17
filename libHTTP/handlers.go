package libhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"restapi/library"
	"time"

	"github.com/gorilla/mux"
)

type HTTPHandlers struct {
	library *library.Library
}

func NewHTTPHandlers(library *library.Library) *HTTPHandlers {
	return &HTTPHandlers{
		library: library,
	}
}

func (h *HTTPHandlers) writeError(w http.ResponseWriter, err error, status int) {
	errDTO := ErrorDTO{
		Message: err.Error(),
		Time:    time.Now(),
	}

	http.Error(w, errDTO.ToString(), status)
}

func (h *HTTPHandlers) isEmpty(w http.ResponseWriter, mapa map[string]library.Book) bool {
	if len(mapa) == 0 {
		msg := "Library is empty / Books not found"
		fmt.Println(msg)
		w.WriteHeader(http.StatusBadRequest)

		h.Output(w, msg)
		return true
	}
	return false
}

func (h *HTTPHandlers) Output(w http.ResponseWriter, t any) {
	b, err := json.MarshalIndent(t, "", "    ")
	if err != nil {
		fmt.Println("Fail to convertation information into JSON")
		h.writeError(w, err, http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(b); err != nil {
		fmt.Println("Fail to write information into http response")
		h.writeError(w, err, http.StatusInternalServerError)
		return
	}
}

/*
	pattern: /books
	method:  POST
	info:    JSON in HTTP request body

succeed:
  - status code:   201 Created
  - response body: JSON represent created book

failed:
  - status code:   400, 409, 500, ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandlerNewBook(w http.ResponseWriter, r *http.Request) {
	var BookDTO BookDTO
	if err := json.NewDecoder(r.Body).Decode(&BookDTO); err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		return
	}

	if err := BookDTO.ValidateforNewBook(); err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		return
	}

	newBook := library.NewBook(BookDTO.Title, BookDTO.Author, BookDTO.Sheets)

	if err := h.library.AddBook(newBook); err != nil {
		if errors.Is(err, library.ErrBookAlreadyExists) { // Обязательно затестить ошибку
			h.writeError(w, err, http.StatusConflict)
		} else {
			h.writeError(w, err, http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	h.Output(w, newBook)
}

/*
	pattern: /books/{title}
	method:  PATCH
	info:    patern + JSON in HTTP request body

succeed:
  - status code:   200 Ok
  - response body: JSON represent changed book: Read: true or false

failed:
  - status code:   400, 404, 500, ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandlerReadBook(w http.ResponseWriter, r *http.Request) {
	var ReadDTO ReadBookDTO
	if err := json.NewDecoder(r.Body).Decode(&ReadDTO); err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		return
	}

	title := mux.Vars(r)["title"]

	var (
		changedBook library.Book
		err         error
	)

	if ReadDTO.Read {
		changedBook, err = h.library.ReadBook(title)
	} else {
		changedBook, err = h.library.UnReadBook(title)
	}

	if err != nil {
		if errors.Is(err, library.ErrBookNotFound) { // Обязательно затестить ошибку
			h.writeError(w, err, http.StatusNotFound)
		} else {
			h.writeError(w, err, http.StatusInternalServerError)
		}
		return
	}

	h.Output(w, changedBook)
}

/*
	pattern: /books/{title}
	method:  GET
	info:    patern

succeed:
  - status code:   200 Ok
  - response body: JSON represent book

failed:
  - status code:   404, 500, ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandlerGetBook(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]

	book, err := h.library.GetBook(title)
	if err != nil {
		if errors.Is(err, library.ErrBookNotFound) {
			h.writeError(w, err, http.StatusNotFound)
			return
		}
		h.writeError(w, err, http.StatusInternalServerError)
		return
	}

	h.Output(w, book)
}

/*
	pattern: /books
	method:  GET
	info:    -

succeed:
  - status code:   200 Ok
  - response body: JSON represent All books

failed:
  - status code:   404, 500, ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandlerAllBook(w http.ResponseWriter, r *http.Request) {
	tLibrary := h.library.AllBooks()

	if h.isEmpty(w, tLibrary) {
		return
	}

	h.Output(w, tLibrary)
}

/*
	pattern: /books/read=true
	method:  GET
	info:    query param

succeed:
  - status code:   200 Ok
  - response body: JSON represent read book

failed:
  - status code:   404, 500, ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandlerAllReadBook(w http.ResponseWriter, r *http.Request) {
	allReadBooks := h.library.AllReadBooks()

	if h.isEmpty(w, allReadBooks) {
		return
	}

	h.Output(w, allReadBooks)
}

/*
	pattern: /books/read=false
	method:  GET
	info:    query param

succeed:
  - status code:   200 Ok
  - response body: JSON represent unread book

failed:
  - status code:   404, 500, ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandlerAllUnReadBook(w http.ResponseWriter, r *http.Request) {
	allUnReadBooks := h.library.AllUnReadBooks()

	if h.isEmpty(w, allUnReadBooks) {
		return
	}

	h.Output(w, allUnReadBooks)
}

/*
	pattern: /books/{author}
	method:  GET
	info:    patern

succeed:
  - status code:   200 Ok
  - response body: JSON represent author's book

failed:
  - status code:   400, 404, 500, ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandlerAllAuthorBook(w http.ResponseWriter, r *http.Request) {
	author := mux.Vars(r)["author"]

	authorBooks := h.library.AuthorBooks(author)

	if h.isEmpty(w, authorBooks) {
		return
	}

	h.Output(w, authorBooks)
}

/*
	pattern: /books/{title}
	method:  DELETE
	info:    -

succeed:
  - status code:   204 No Content
  - response body: JSON represent author's book

failed:
  - status code:   404, 500, ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandlerDeleteBook(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]

	err := h.library.DeleteBook(title)
	if err != nil {
		if errors.Is(err, library.ErrBookNotFound) {
			h.writeError(w, err, http.StatusNotFound)
			return
		}
		h.writeError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
