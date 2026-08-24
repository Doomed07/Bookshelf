# Library — REST API for Library Management

A REST API in Go for managing a library's book catalog, with a custom-built HTTP handling layer instead of relying on a full web framework.

## Features

- Create, retrieve, update (mark read/unread), and delete books
- Filter books by read status or by author
- Thread-safe in-memory storage (`sync.RWMutex`)
- Custom `libHTTP` package wrapping [gorilla/mux](https://github.com/gorilla/mux) for routing and response handling

## Tech stack

Go 1.26 · gorilla/mux

## Project structure

```
.
├── main.go
├── libHTTP/    # HTTP layer: routing, handlers, DTOs
└── library/    # domain logic: Book, Library (storage)
```

## API

| Method | Endpoint                  | Description                     |
|--------|----------------------------|-----------------------------------|
| POST   | `/books`                  | Add a new book                   |
| GET    | `/books`                  | Get all books                    |
| GET    | `/books?read=true`        | Get books marked as read         |
| GET    | `/books?read=false`       | Get unread books                 |
| GET    | `/books/{title}`          | Get a book by title              |
| GET    | `/books/author/{author}`  | Get all books by an author       |
| PATCH  | `/books/{title}`          | Mark a book as read/unread       |
| DELETE | `/books/{title}`          | Delete a book                    |

## Running locally

```bash
go run main.go
```

Server starts on `localhost:9999`.

## Possible next steps

- Persist books in PostgreSQL instead of in-memory storage
- Add unit tests
- Add pagination for large catalogs
