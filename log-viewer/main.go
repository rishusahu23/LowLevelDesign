package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type LogWatcher struct {
	clients     map[*websocket.Conn]bool
	logFilePath string
}

func NewLogWatcher(logFilePath string) *LogWatcher {
	return &LogWatcher{
		clients:     make(map[*websocket.Conn]bool),
		logFilePath: logFilePath,
	}
}

func (lw *LogWatcher) addClient(conn *websocket.Conn) {
	lw.clients[conn] = true
}

func (lw *LogWatcher) removeClient(conn *websocket.Conn) {
	delete(lw.clients, conn)
	conn.Close()
}

func (lw *LogWatcher) broadcast(message string) {
	for client := range lw.clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			lw.removeClient(client)
		}
	}
}

func (lw *LogWatcher) watchLogFile() {
	file, err := os.Open(lw.logFilePath)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	file.Seek(0, os.SEEK_END)
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				time.Sleep(1 * time.Second)
				continue
			}
			log.Fatalf("Error reading from log file: %v", err)
		}
		lw.broadcast(line)
	}
}

func (lw *LogWatcher) serveWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	lw.addClient(conn)
	defer lw.removeClient(conn)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func handleFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func main() {
	logFilePath := "logfile.log" // Replace with the path to your log file
	logWatcher := NewLogWatcher(logFilePath)

	go logWatcher.watchLogFile()

	http.HandleFunc("/ws", logWatcher.serveWebSocket)
	http.HandleFunc("/log", handleFile)

	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
