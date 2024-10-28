package main

import (
	"net/http"
	// "os"
	"sync"

	"github.com/gorilla/websocket"
	// "github.com/rs/zerolog"
	// "github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// WebSocketLogger broadcasts logs to all WebSocket clients
type WebSocketLogger struct {
	clients map[*websocket.Conn]bool
	lock    sync.Mutex
}

func NewWebSocketLogger() *WebSocketLogger {
	return &WebSocketLogger{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (wsl *WebSocketLogger) AddClient(conn *websocket.Conn) {
	wsl.lock.Lock()
	defer wsl.lock.Unlock()
	wsl.clients[conn] = true
}

func (wsl *WebSocketLogger) RemoveClient(conn *websocket.Conn) {
	wsl.lock.Lock()
	defer wsl.lock.Unlock()
	delete(wsl.clients, conn)
	conn.Close()
}

func (wsl *WebSocketLogger) Write(p []byte) (n int, err error) {
	wsl.lock.Lock()
	defer wsl.lock.Unlock()
	for client := range wsl.clients {
		if err := client.WriteMessage(websocket.TextMessage, p); err != nil {
			wsl.RemoveClient(client)
		}
	}
	return len(p), nil
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (wsl *WebSocketLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upgrade connection")
		return
	}
	wsl.AddClient(conn)
}

func initWsLog() *WebSocketLogger{
	// Create WebSocket logger
	wsLogger := NewWebSocketLogger()

	// Configure zerolog to write to both the WebSocket logger and console
	// log.Logger = log.Output(wsLogger)

	// Set up WebSocket route
	http.Handle("/ws", wsLogger)

	// go func() {
	// 	// Simulate log messages every second
	// 	for {
	// 		log.Info().Msg("Info log message for WebSocket clients")
	// 		log.Error().Msg("Error log message for WebSocket clients")
	// 		time.Sleep(1 * time.Second)
	// 	}
	// }()

	// Start HTTP server for WebSocket connections
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal().Err(err).Msg("listener failed")
		}
	}()

	return wsLogger
}
