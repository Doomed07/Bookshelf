package libhttp

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	httpHandlers *HTTPHandlers
}

func NewHTTPServer(httpHandlers *HTTPHandlers) *HTTPServer {
	return &HTTPServer{
		httpHandlers: httpHandlers,
	}
}

func (s *HTTPServer) StartServer() error {
	router := mux.NewRouter()

	router.Path("/books").Methods("POST").HandlerFunc(s.httpHandlers.HandlerNewBook)
	router.Path("/books/{title}").Methods("PATCH").HandlerFunc(s.httpHandlers.HandlerReadBook)
	router.Path("/books/{title}").Methods("GET").HandlerFunc(s.httpHandlers.HandlerGetBook)
	router.Path("/books").Methods("GET").Queries("read", "true").HandlerFunc(s.httpHandlers.HandlerAllReadBook)
	router.Path("/books").Methods("GET").Queries("read", "false").HandlerFunc(s.httpHandlers.HandlerAllUnReadBook)
	router.Path("/books").Methods("GET").HandlerFunc(s.httpHandlers.HandlerAllBook)
	router.Path("/books/author/{author}").Methods("GET").HandlerFunc(s.httpHandlers.HandlerAllAuthorBook)
	router.Path("/books/{title}").Methods("DELETE").HandlerFunc(s.httpHandlers.HandlerDeleteBook)

	if err := http.ListenAndServe(":9999", router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}

	return nil
}
