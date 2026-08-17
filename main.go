package main

import (
	"fmt"
	libhttp "restapi/libHTTP"
	"restapi/library"
)

func main() {
	library := library.NewLibrary()
	hadlers := libhttp.NewHTTPHandlers(library)
	server := libhttp.NewHTTPServer(hadlers)

	if err := server.StartServer(); err != nil {
		fmt.Println("Fail to start server:", err)
	}

}
