package main

import (
	"log"
	"net/http"
)

func main() {
	// port config
	const port = "8080"

	mux := http.NewServeMux()

	serv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	// log confirmation that the port is running
	log.Printf("Server running on port: %s\n", port)

	// if the server fails it will log the error
	log.Fatal(serv.ListenAndServe())
}
