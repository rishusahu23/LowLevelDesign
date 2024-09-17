package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type FeedBack struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Source  string `json:"source"`
}

func handleFeedback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	var feedback FeedBack
	err := json.NewDecoder(r.Body).Decode(&feedback)
	if err != nil {
		http.Error(w, "Invalid feedback data", http.StatusBadRequest)
		return
	}

	processFeedback(feedback)
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, "Feedback received")
}

func processFeedback(feedback FeedBack) {
	fmt.Printf("Processing feedback: %+v\n", feedback)
}

func main() {
	http.HandleFunc("/feedback", handleFeedback)
	fmt.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
