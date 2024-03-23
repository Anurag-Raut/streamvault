package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"os/exec"
	"streamvault/postgres"

	// "github.com/rs/cors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	// "strings"
)

type StartStreamRequest struct {
	Title string `json:"title"`
}

type StreamRequest struct {
	StreamId string `json:"streamId"`
	// StreamData []byte `json:"data"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	Subprotocols:    []string{"Bearer"}, // <-- add this line
}

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Home Page")
}

var number = 1

func wsEndpoint(w http.ResponseWriter, r *http.Request) {

	var a string = r.Header.Get("Sec-WebSocket-Protocol")
	parts := strings.Split(a, " ")
	if len(parts) != 2 {
		fmt.Println("Invalid Sec-WebSocket-Protocol header")
		return
	}
	token := parts[1]
	fmt.Println(token)

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}

	fmt.Println("Client connected %d", number)

	defer conn.Close()

	cmd := exec.Command("ffmpeg",
		"-i", "pipe:0",
		"-c:v", "libx264", "-preset", "ultrafast", "-tune", "zerolatency",
		"-c:a", "aac", "-ar", "44100", "-b:a", "64k",
		"-f", "hls",
		"-g", "20",
		"-hls_time", "2",
		"-hls_list_size", "0",
		// `/home/anurag/projects/streamvault/packages/backend/hls/output.m3u8`,
		fmt.Sprintf("/home/anurag/projects/streamvault/packages/backend/hls/stream%d.ts", number),
	)
	number++
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

				fmt.Println("Error reading messagesss:", err)

				break
			}
			// fmt.Printf("Received message: %s\n", message)
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
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var streamReq StartStreamRequest
	if err := json.NewDecoder(r.Body).Decode(&streamReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("Starting stream with title:", streamReq.Title)

	fmt.Println("Starting stream")
	streamId, err := postgres.AddStream(streamReq.Title)

	if err != nil {
		http.Error(w, "Error starting stream", http.StatusInternalServerError)
		return
	}

	fmt.Println("Stream ID:", streamId)

	fmt.Println("Stream started")
	response := struct {
		StreamID string `json:"streamId"`
	}{
		StreamID: streamId,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Adjust the allowed methods accordingly
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Write(responseJSON)

}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("Username:", loginRequest.Username)
	fmt.Println("Password:", loginRequest.Password)

	if loginRequest.Username == "anurag" && loginRequest.Password == "password" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Login successful"))

	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Login failed"))
	}

}

func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("Username:", loginRequest.Username)
	fmt.Println("Password:", loginRequest.Password)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": loginRequest.Username,
		"password": loginRequest.Password,
	})

	tokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tokenString))

}
func authMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the incoming request method and URL
		println("Incoming request:", r.Method, r.URL.Path)

		authHeader := r.Header.Get("Authorization")

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// You should provide the secret key or the key used for signing the token here
			return []byte("secret"), nil
		})

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func setupRoutes(mux *http.ServeMux) {
	mux.Handle("/", authMiddleWare(http.HandlerFunc(homePage)))
	mux.HandleFunc("/ws", wsEndpoint)
	mux.HandleFunc("/startStream", startStream)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/signup", SignUp)
}

func main() {
	mux := http.NewServeMux()

	postgres.Connect()
	defer postgres.Disconnect()
	setupRoutes(mux)
	fmt.Println("Hello, World!")

	handler := cors.Default().Handler(mux)

	http.ListenAndServe(":8080", handler)

}
