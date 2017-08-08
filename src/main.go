package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Save pointer to Websocket (connected clients)
var clients = make(map[*websocket.Conn]bool)

// bordacast channel
var broadcast = make(chan Message)

// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Message structure
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	// Simple file server for js and css
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	// Websocket route
	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	// Start server localhost:8000 and log any errors
	log.Println("http server started on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade first get request to websocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

	// Register new client
	clients[ws] = true

	for {
		var msg Message

		// Read message as Json and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error %v", err)
			delete(clients, ws)
			break
		}

		broadcast <- msg

	}
}

func handleMessages() {
	for {
		msg := <-broadcast

		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
