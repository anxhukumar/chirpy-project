package main

import (
	"log"
	"net/http"
)

func main() {
	// port config
	const filePathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()

	// add file config
	mux.Handle("/", http.FileServer(http.Dir(filePathRoot)))

	serv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	// log confirmation that the port is running
	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)

	// if the server fails it will log the error
	log.Fatal(serv.ListenAndServe())
}
