package main

import (
        "fmt"
        "log"
        "net/http"
        "os"
)

func main() {
        // Define the port number
        const port = 80

        // Setup a basic endpoint
        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
}