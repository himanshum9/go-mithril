package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func main() {
	router := mux.NewRouter()

	// WebSocket endpoint
	router.HandleFunc("/ws", wsHandler)

	// HTTP endpoint to receive location data and broadcast to WebSocket clients
	router.HandleFunc("/stream", streamHandler).Methods("POST")

	log.Println("Starting streaming service on port 8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()
	log.Println("WebSocket client connected")
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
	clientsMu.Lock()
	delete(clients, conn)
	clientsMu.Unlock()
	conn.Close()
	log.Println("WebSocket client disconnected")
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	var msg map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	broadcastToWebSocketClients(msg)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Streaming location data..."))
}

func broadcastToWebSocketClients(msg interface{}) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for client := range clients {
		err := client.WriteJSON(msg)
		if err != nil {
			log.Printf("WebSocket write error: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}
