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

	// add routing
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot))))
	mux.HandleFunc("/healthz", handlerReadiness)

	serv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	// log confirmation that the port is running
	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)

	// if the server fails it will log the error
	log.Fatal(serv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
