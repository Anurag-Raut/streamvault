package main

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gorilla/websocket"
	
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Home Page")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	cmd := exec.Command("ffmpeg",
		"-i", "pipe:0",
		"-c:v", "libx264", "-preset", "ultrafast", "-tune", "zerolatency",
		"-c:a", "aac", "-ar", "44100", "-b:a", "64k",
		"-f", "hls",
		"-g", "20",
		"-hls_time", "2",
		"-hls_list_size", "0",
		"/home/anurag/projects/streamvault/packages/backend/hls/output.m3u8",
	)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("Error getting stdin pipe:", err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	// Start a goroutine to read messages from the WebSocket connection
	go func() {
		defer stdin.Close()
		defer conn.Close() // Close the WebSocket connection when this goroutine exits
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Println("Error reading message:", err)
				}
				break
			}
			fmt.Printf("Received message: %s\n", message)

			if _, err := stdin.Write(message); err != nil {
				fmt.Println("Error writing message:", err)
				break
			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for command:", err)
		return
	}
}

func startStream(w http.ResponseWriter, r *http.Request) {



}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
	http.HandleFunc("startStream", startStream)
}

func main() {
	setupRoutes()
	fmt.Println("Hello, World!")
	http.ListenAndServe(":8080", nil)
}
