package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

func clearHandler(w http.ResponseWriter, r *http.Request) {
	app := r.URL.Query().Get("app")
	if app == "" {
		http.Error(w, "Missing 'app' parameter", http.StatusBadRequest)
		return
	}

	var cmd *exec.Cmd

	switch app {
	case "chrome":
		cmd = exec.Command("bash", "-c", "rm -rf ~/.config/google-chrome/Default/Cache ~/.config/google-chrome/Default/History")
	case "firefox":
		cmd = exec.Command("bash", "-c", "rm -rf ~/.mozilla/firefox/*/*cache* ~/.mozilla/firefox/*/places.sqlite")
	default:
		http.Error(w, "Unsupported browser", http.StatusBadRequest)
		return
	}

	if err := cmd.Run(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to clear cache and history: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Cache and history cleared for %s", app)
}

func openHandler(w http.ResponseWriter, r *http.Request) {
	app := r.URL.Query().Get("app")
	url := r.URL.Query().Get("url")
	if app == "" || url == "" {
		http.Error(w, "Missing 'app' or 'url' parameter", http.StatusBadRequest)
		return
	}

	var cmd *exec.Cmd

	switch app {
	case "chrome":
		cmd = exec.Command("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", url)
	case "firefox":
		cmd = exec.Command("firefox", url)
	default:
		http.Error(w, "Unsupported browser", http.StatusBadRequest)
		return
	}

	if err := cmd.Start(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to open URL: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "URL %s opened in %s", url, app)
}

func closeHandler(w http.ResponseWriter, r *http.Request) {
	app := r.URL.Query().Get("app")
	if app == "" {
		http.Error(w, "Missing 'app' parameter", http.StatusBadRequest)
		return
	}

	var cmd *exec.Cmd

	switch app {
	case "chrome":
		cmd = exec.Command("pkill", "chrome")
	case "firefox":
		cmd = exec.Command("pkill", "firefox")
	default:
		http.Error(w, "Unsupported browser", http.StatusBadRequest)
		return
	}

	if err := cmd.Run(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to close browser: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Browser %s closed", app)
}

func main() {
	http.HandleFunc("/clear", clearHandler)
	http.HandleFunc("/open", openHandler)
	http.HandleFunc("/close", closeHandler)

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
