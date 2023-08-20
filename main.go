package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")

}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Client Connected")

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error during message reading:", err)
			break
		}

		fmt.Printf("Received: %s\n", p)

		echo := string(p) + " echoed"

		if err := conn.WriteMessage(messageType, []byte(echo)); err != nil {
			fmt.Println("Error during message writing:", err)
			break
		}
	}
}

func main() {

	run()
	return

	Init()
	// Define the port number
	const port = 80

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	createPlaylist(&Playlist{Name: "adsddsad", Songs: []Song{}})

	// Setup a basic endpoint
	http.HandleFunc("/api/playlist", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		if r.Method == http.MethodGet {

			dtoBs, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}

			var dto struct {
				Id string `json:"id"`
			}

			err = json.Unmarshal(dtoBs, &dto)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}

			pl, err := getPlaylistById(dto.Id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}

			plBs, err := json.Marshal(pl)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}

			w.Write(plBs)

		} else if r.Method == http.MethodPost {

			var pl Playlist

			err := json.NewDecoder(r.Body).Decode(&pl)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}

			err = createPlaylist(&pl)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}

			w.Write(nil)
			return

		}

		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/api/playlists", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		if r.Method == http.MethodGet {

			pls := listPlaylists()

			plBs, err := json.Marshal(pls)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}

			w.Write(plBs)
			return

		}

		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/api/msg", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		d, _ := io.ReadAll(r.Body)

		log.Printf("new msg: %s", string(d))

		fmt.Fprint(w, "Hello, world!")
	})
	http.HandleFunc("/api/msg/ws", handleConnection)

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
