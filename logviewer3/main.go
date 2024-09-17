package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	logfile  = "logviewer3/logfile"
	clients  = make(map[*websocket.Conn]bool)
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	broadcast = make(chan string)
)

func main() {
	http.HandleFunc("/ws", handleConnections)
	go broadCastLogUpdates()
	go handleLogUpdates()
	fmt.Println("server started on port: 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("failed to upgrade", err)
		return
	}
	defer ws.Close()
	clients[ws] = true
	sendLast10lines(ws)

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("error reading message", err)
			delete(clients, ws)
			break
		}
	}
}

func handleLogUpdates() {
	file, err := os.Open(logfile)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return
	}
	defer file.Close()
	file.Seek(0, os.SEEK_END)
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		broadcast <- line
	}
}

func broadCastLogUpdates() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func sendLast10lines(ws *websocket.Conn) {
	file, err := os.Open(logfile)
	if err != nil {
		fmt.Println("error in openning file", err)
		return
	}
	defer file.Close()
	var lines []string
	var bufferSize = 1024
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("error in getting file info", err)
		return
	}
	fileSize := fileInfo.Size()
	offset := int(fileSize)
	var lineBuffer []byte

	for {
		if offset < bufferSize {
			bufferSize = offset
		}
		offset -= bufferSize
		file.Seek(int64(offset), os.SEEK_SET)
		buffer := make([]byte, bufferSize)
		bytesRead, err := file.Read(buffer)
		if err != nil {
			fmt.Println("failed to read file", err)
			return
		}
		lineBuffer = append(buffer[:bytesRead], lineBuffer...)
		lines = strings.Split(string(lineBuffer), "\n")
		if len(lines) > 10 {
			lines = lines[len(lines)-9:]
			break
		}
		if offset == 0 {
			break
		}
	}
	for _, line := range lines {
		ws.WriteMessage(websocket.TextMessage, []byte(line))
	}
}
