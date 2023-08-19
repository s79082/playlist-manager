package main

import (
        //"fmt"
        //"log"
        "net/http"
        //"os"
)
func enableCors(w *http.ResponseWriter) {
        (*w).Header().Set("Access-Control-Allow-Origin", "*")
        
}

/*func main() {
        // Define the port number
        const port = 80

        fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

        // Setup a basic endpoint
        http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
                enableCors(&w)
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

*/