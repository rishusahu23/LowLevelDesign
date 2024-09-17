package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func watchLogFile(logFilePath string) {
	file, err := os.Open(logFilePath)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	// Move to the end of the file
	file.Seek(0, os.SEEK_END)
	reader := bufio.NewReader(file)

	fmt.Println("Watching log file for changes...")

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				// End of file, sleep and continue
				time.Sleep(1 * time.Second)
				continue
			}
			log.Fatalf("Error reading from log file: %v", err)
		}
		// Print the new line to the terminal
		fmt.Print(line)
	}
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}

	// Print the current working directory
	fmt.Println("Current working directory:", dir)
	logFilePath := "browserstack/logfile.log" // Replace with the path to your log file

	watchLogFile(logFilePath)
}
