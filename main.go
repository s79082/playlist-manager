package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")

}

func main() {

	conn()

	// Define the port number
	const port = 80

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// Setup a basic endpoint
	http.HandleFunc("/api/playlist", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		if r.Method == http.MethodGet {
			return

		}

		fmt.Fprint(w, "Hello, world!")
	})
	http.HandleFunc("/api/msg", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		d, _ := io.ReadAll(r.Body)

		log.Printf("new msg: %s", string(d))

		fmt.Fprint(w, "Hello, world!")
	})

	// Start the HTTP server
	address := fmt.Sprintf(":%d", port)
	log.Printf("Starting HTTP server on %s", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		os.Exit(3)
		log.Fatalf("Error starting server: %v", err)
	}

	println("hewwlo")

}
