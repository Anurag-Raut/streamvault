package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"os/exec"
	"streamvault/auth"
	"streamvault/postgres"

	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type StreamRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Thumbnail   string `json:"thumbnail"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	Subprotocols:    []string{"streamId"}, // <-- add this line
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Home Page")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {

	var a string = r.Header.Get("Sec-WebSocket-Protocol")
	parts := strings.Split(a, " ")
	if len(parts) != 2 {
		fmt.Println("Invalid Sec-WebSocket-Protocol header")
		return
	}
	var streamId string = parts[1]
	fmt.Println(streamId)

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}

	failOnError(err, "Failed to declare a queue")
	// fmt.Println("Client connected %d", number)
	dirPath := fmt.Sprintf("/home/anurag/s3mnt/%s", streamId)

	_, osserr := os.Stat(dirPath)

	if os.IsNotExist(osserr) {
		err := os.Mkdir(dirPath, 0777)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
		println("Directory created")
	}

	defer conn.Close()

	cmd := exec.Command("ffmpeg",
		"-re",
		"-i", "pipe:0",
		"-c:v", "libx264", "-preset", "ultrafast", "-tune", "zerolatency",
		"-c:a", "aac", "-ar", "44100", "-b:a", "64k",
		"-f", "hls",
		"-g", "20",
		"-hls_time", "5",
		"-hls_list_size", "0",
		// `/home/anurag/projects/streamvault/packages/backend/hls/output.m3u8`,
		fmt.Sprintf("%s/%s.m3u8", dirPath, streamId),
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
	// get the stream data which is json
	var streamRequest StreamRequest
	if err := json.NewDecoder(r.Body).Decode(&streamRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	streamId, err := postgres.AddStream(streamRequest.Title, streamRequest.Description, streamRequest.Category, streamRequest.Thumbnail)
	if err != nil {
		http.Error(w, "Error adding stream", http.StatusInternalServerError)
		return
	}
	responseJson := fmt.Sprintf(`{"streamId": "%s"}`, streamId)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Adjust the allowed methods accordingly
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Write([]byte(responseJson))
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


func authMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the incoming request method and URL
		println("Incoming request:", r.Method, r.URL.Path)

		cookie, err := r.Cookie("jwt")
		if err != nil {
			http.Error(w, "No token found", http.StatusUnauthorized)
			return
		}
		tokenString := cookie.Value

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// You should provide the secret key or the key used for signing the token here
			return []byte("eat shit"), nil
		})

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Access the username claim
			if username, exists := claims["username"].(string); exists {
				// Now you have the username
				fmt.Println("Username:", username)
				userExists, _ := postgres.UserExists(username)

				if !userExists {
					http.Error(w, "User does not exits ", http.StatusUnauthorized)
					return
				}

			} else {
				http.Error(w, "Username claim not found", http.StatusUnauthorized)
				return
			}
		} else {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
func uploadThumbnail(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	thumbnail, thumbnailHeader, err := r.FormFile("thumbnail")
	if err != nil {
		http.Error(w, "Error parsing thumbnail", http.StatusBadRequest)
		return
	}
	defer thumbnail.Close()

	thumbnailExtension := filepath.Ext(thumbnailHeader.Filename)
	thumbnailExtension = strings.ToLower(thumbnailExtension)

	imageData, err := io.ReadAll(thumbnail)
	if err != nil {
		http.Error(w, "Error reading thumbnail", http.StatusBadRequest)
		return
	}

	// Get the home directory
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		http.Error(w, "Unable to get home directory", http.StatusInternalServerError)
		return
	}

	uploadDir := filepath.Join(homeDir, "s3mnt", "thumbnail")
	err = os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating directory: %v", err), http.StatusInternalServerError)
		return
	}

	imageName := uuid.New().String()
	thumbnailPath := filepath.Join(uploadDir, imageName+thumbnailExtension)
	err = os.WriteFile(thumbnailPath, imageData, os.ModePerm)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error writing thumbnail: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "thumbnail/"+imageName+thumbnailExtension)
}

func setupRoutes(mux *http.ServeMux) {
	mux.Handle("/", authMiddleWare(http.HandlerFunc(homePage)))
	mux.HandleFunc("/ws", wsEndpoint)
	// mux.Handle("/startStream", authMiddleWare(http.HandlerFunc(startStream)))
	mux.HandleFunc("/startStream", startStream)
	mux.HandleFunc("/uploadThumbnail", uploadThumbnail)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/signup", auth.SignUp)
	mux.HandleFunc("/streams", postgres.GetStreams)
	mux.HandleFunc("/signIn", auth.SignIn)
	mux.Handle("/hls/", http.StripPrefix("/hls/", http.FileServer(http.Dir("/home/anurag/s3mnt"))))

}

func main() {
	mux := http.NewServeMux()

	postgres.Connect()
	defer postgres.Disconnect()
	setupRoutes(mux)
	fmt.Println("Hello, World!")

	handler := corsMiddleware(mux)

	// conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	// failOnError(err, "Failed to connect to RabbitMQ")

	// defer conn.Close()

	// go subtitle.StartSubtitleServer()

	http.ListenAndServe(":8080", handler)

}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Allow preflight requests to any origin with the specified methods
		if req.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.WriteHeader(http.StatusOK)
			return
		}

		// Set CORS headers for non-preflight requests
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Call the next handler in the chain
		next.ServeHTTP(w, req)
	})
}
