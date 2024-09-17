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
	logfile  = "logviewer2/logfile.text"
	clients  = make(map[*websocket.Conn]bool)
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins
		},
	}
	broadCast = make(chan string)
)

func main() {
	http.HandleFunc("/ws", handleConnections)

	go broadCastLogUpdates()
	go handleLogUpdates()

	fmt.Println("server started on : 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("failed to upgrade websocket", err)
		return
	}
	defer ws.Close()

	clients[ws] = true
	sendLast10Lines(ws)

	for {
		_, _, err = ws.ReadMessage()
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
		broadCast <- line
	}
}

func sendLast10Lines(ws *websocket.Conn) {
	file, err := os.Open(logfile)
	if err != nil {
		fmt.Println("failed to open file, ", err)
		return
	}
	defer file.Close()

	var lines []string
	var bufferSize = 1024
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("Failed to get file info: %v\n", err)
		return
	}
	fileSize := fileInfo.Size()
	var offset = 0
	var linebuffer []byte

	for {
		if int(fileSize)-offset < bufferSize {
			bufferSize = int(fileSize) - offset
		}
		offset -= bufferSize
		file.Seek(int64(offset), os.SEEK_SET)
		buffer := make([]byte, bufferSize)
		bytesRead, err := file.Read(buffer)
		if err != nil {
			fmt.Printf("Failed to read file: %v\n", err)
			return
		}

		linebuffer = append(buffer[:bytesRead], linebuffer...)
		lines = strings.Split(string(linebuffer), "\n")
		if len(lines) > 10 {
			lines = lines[len(lines)-10:]
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

func broadCastLogUpdates() {
	for {
		msg := <-broadCast
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				fmt.Printf("Error broadcasting message: %v\n", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
