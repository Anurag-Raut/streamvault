package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"os/exec"
	"streamvault/auth"
	"streamvault/chat"
	"streamvault/postgres"
	"streamvault/rmq"
	"streamvault/utils"

	"os"

	"github.com/gofor-little/env"
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

func SendToSubtitler(message, streamId string, duration, totalDuration float64, segmentNumber int) error {
	var response struct {
		StreamId      string  `json:"streamId"`
		Message       string  `json:"message"`
		Duration      float64 `json:"duration"`
		SegmentNumber int     `json:"segmentNumber"`
		TotalDuration float64 `json:"totalDuration"`
	}
	fmt.Println("Sending to subtitler:", message)

	response.StreamId = streamId
	response.Message = message
	response.Duration = duration
	response.SegmentNumber = segmentNumber
	response.TotalDuration = totalDuration

	jsonPayload, err := json.Marshal(response)
	if err != nil {
		return err
	}
	fmt.Println("jsonPayload:")
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/receive_text", env.Get("SUBTITLER_API_URL", "http://subtitler:5000")), bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return err
	}
	var responseText struct {
		Message string `json:"message"`
		Success bool   `json:"success"`
	}
	err = json.Unmarshal(body, &responseText)

	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return err
	}
	fmt.Println("Response from subtitler:", responseText.Message, responseText.Success)

	return nil

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

	postgres.UpdateStatus(streamId, true)

	cmd := exec.Command("ffmpeg",
		"-re",
		"-i", "pipe:0",
		"-c:v", "libx264", "-preset", "ultrafast", "-tune", "zerolatency",
		"-c:a", "aac", "-ar", "44100", "-b:a", "64k",
		"-f", "hls",
		"-g", "20",
		"-hls_time", "5",
		"-hls_list_size", "0",
		"-progress", "pipe:1",
		// `/home/anurag/projects/streamvault/packages/backend/hls/output.m3u8`,
		fmt.Sprintf("%s/%s.m3u8", dirPath, streamId),
	)
	stdin, err := cmd.StdinPipe()

	if err != nil {
		fmt.Println("Error getting stdin pipe:", err)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("Error getting stderr pipe:", err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	defer func() {
		conn.Close()
		postgres.UpdateStatus(streamId, false)
		fmt.Println("WebSocket connection closed")
	}()

	go func() {
		defer stderr.Close()
		scanner := bufio.NewScanner(stderr)
		var totalDuration float64
		for scanner.Scan() {
			line := scanner.Text()
			// Parse the progress information from the stderr output
			// Progress information typically starts with "frame="
			if strings.Contains(line, "Opening '") && strings.Contains(line, "' for writing") {
				// Extract the file path between the single quotes
				start := strings.Index(line, "'") + 1
				end := strings.LastIndex(line, "'")
				filePath := line[start:end]
				parts := strings.Split(filePath, "/")

				// Get the file name from the file path
				if strings.HasSuffix(filePath, ".ts") && !strings.HasSuffix(filePath, ".m3u8.tmp") {
					cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", filePath)
					output, err := cmd.CombinedOutput()
					if err != nil {
						fmt.Println("Error executing ffprobe command:", err)
						return
					}

					durSting := string(output)
					durSting = strings.TrimSuffix(durSting, "\n")
					duration, err := strconv.ParseFloat(durSting, 64)

					if err != nil {
						fmt.Println("Error converting duration:", err)

						return
					}

					fmt.Println("Duration:", duration)
					var a = strings.Split(parts[len(parts)-1], ".")[0]
					segmentNumberString := strings.TrimPrefix(a, streamId)
					segmentNumber, err := strconv.Atoi(segmentNumberString)
					if err != nil {
						fmt.Println("Error converting segment number:", err)
						return
					}

					fmt.Println("Segment Number:", segmentNumber)

					err = SendToSubtitler(parts[len(parts)-2]+"/"+parts[len(parts)-1], streamId, duration, totalDuration, segmentNumber)
					if err != nil {
						fmt.Println("Error sending to subtitler:", err)
						return
					}
					totalDuration += duration

				}

			}
		}
	}()

	go func() {
		defer stdin.Close()
		defer conn.Close() // Close the WebSocket connection when this goroutine exits
		defer func() {
			postgres.UpdateStatus(streamId, false)
			_, err := http.Post(fmt.Sprintf("%s/stop_transcription", env.Get("SUBTITLER_API_URL", "http://subtitler:5000")), "application/json", bytes.NewBuffer([]byte(fmt.Sprintf(`{"streamId": "%s"}`, streamId))))
			if err != nil {
				fmt.Println("Error stopping transcription:", err)
			}

		}()
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
		utils.SendError(w, "Invalid request bauthMiddleWareody", http.StatusBadRequest)
		return
	}
	userId := r.Context().Value("userId").(string)

	streamId, err := postgres.AddStream(streamRequest.Title, streamRequest.Description, streamRequest.Category, streamRequest.Thumbnail, userId)
	if err != nil {
		fmt.Println("Error adding stream:", err)
		utils.SendError(w, "Error adding stream", http.StatusInternalServerError)
		return
	}

	responseJson := fmt.Sprintf(`{"streamId": "%s"}`, streamId)
	_, err = http.Post(fmt.Sprintf("%s/start_transcription", env.Get("SUBTITLER_API_URL", "http://subtitler:5000")), "application/json", bytes.NewBuffer([]byte(responseJson)))
	if err != nil {
		fmt.Println("Error starting transcription:", err)
		utils.SendError(w, "Error starting transcription", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", env.Get("FRONTEND_URL", "https://streamvault.vercel.app"))
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
			utils.SendError(w, "No token found", http.StatusUnauthorized)
			return
		}
		tokenString := cookie.Value

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("eat shit"), nil
		})

		if err != nil {
			// http.Error(w, "Invalid token", http.StatusUnauthorized)
			utils.SendError(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			utils.SendError(w, "Invalid Token", http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Access the username claim
			if userId, exists := claims["userId"].(string); exists {
				// Now you have the username
				fmt.Println("userId:", userId)
				userExists, _ := postgres.UserExists(userId)

				if !userExists {
					utils.SendError(w, "User does not exits", http.StatusUnsupportedMediaType)
					return
				}

				ctx := context.WithValue(r.Context(), "userId", userId)
				r = r.WithContext(ctx)

			} else {
				// http.Error(w, "Username claim not found", http.StatusUnauthorized)
				utils.SendError(w, "UserId clain not found", http.StatusUnauthorized)
				return
			}
		} else {
			// http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			utils.SendError(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func getVideoDataMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the incoming request method and URL
		println("Incoming request:", r.Method, r.URL.Path)

		cookie, err := r.Cookie("jwt")
		if err != nil {
			fmt.Println("No token found")
			fmt.Println(err)
			next.ServeHTTP(w, r)
			return
		}
		tokenString := cookie.Value

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// You should provide the secret key or the key used for signing the token here
			return []byte("eat shit"), nil
		})
		fmt.Println(token, "token")

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
			if userId, exists := claims["userId"].(string); exists {
				// Now you have the username
				fmt.Println("userId:", userId)
				userExists, _ := postgres.UserExists(userId)

				if !userExists {

					http.Error(w, "User does not exits ", http.StatusUnauthorized)
					return
				}

				ctx := context.WithValue(r.Context(), "userId", userId)
				r = r.WithContext(ctx)

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
	err := r.ParseMultipartForm(100 << 20) // 10 MB max
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
	homeDir := "/home/anurag/"
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

	// fmt.Fprint(w, "thumbnail/"+imageName+thumbnailExtension)
	var thumbnailPathResponse struct {
		ThumbnailPath string `json:"thumbnailPath"`
	}
	thumbnailPathResponse.ThumbnailPath = "thumbnail/" + imageName + thumbnailExtension
	responseJson, err := json.Marshal(thumbnailPathResponse)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
	w.Write(responseJson)
}

func UploadProfileImage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(100 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	profileImage, profileImageHeader, err := r.FormFile("profileImage")
	if err != nil {
		http.Error(w, "Error parsing profileImage", http.StatusBadRequest)
		return
	}
	defer profileImage.Close()

	profileImageExtension := filepath.Ext(profileImageHeader.Filename)
	profileImageExtension = strings.ToLower(profileImageExtension)

	imageData, err := io.ReadAll(profileImage)
	if err != nil {
		http.Error(w, "Error reading profileImage", http.StatusBadRequest)
		return
	}

	// Get the home directory
	homeDir := "/home/anurag/"
	if homeDir == "" {
		http.Error(w, "Unable to get home directory", http.StatusInternalServerError)
		return
	}

	uploadDir := filepath.Join(homeDir, "s3mnt", "profileImage")
	err = os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating directory: %v", err), http.StatusInternalServerError)
		return
	}

	imageName := uuid.New().String()
	profileImagePath := filepath.Join(uploadDir, imageName+profileImageExtension)
	err = os.WriteFile(profileImagePath, imageData, os.ModePerm)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error writing profileImage: %v", err), http.StatusInternalServerError)
		return
	}

	// fmt.Fprint(w, "profileImage/"+imageName+profileImageExtension)
	var profileImagePathResponse struct {
		ProfileImagePath string `json:"profileImagePath"`
	}
	profileImagePathResponse.ProfileImagePath = "profileImage/" + imageName + profileImageExtension
	responseJson, err := json.Marshal(profileImagePathResponse)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
	w.Write(responseJson)
}
func uploadVideo(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	err := r.ParseMultipartForm(1 << 30) // 1 GB
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get file data from request
	file, handler, err := r.FormFile("file")
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	videoid := uuid.New().String()
	uploadPath := "/home/anurag/s3mnt/vod"

	// Create the uploads directory if it doesn't exist
	err = os.MkdirAll(uploadPath, 0777)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fileExt := filepath.Ext(handler.Filename)

	fileName := videoid + fileExt
	fmt.Println(fileName)

	// Create a new file in the server
	outFile, err := os.Create(fmt.Sprintf("%s/%s", uploadPath, fileName))
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer outFile.Close()
	videoData, err := io.ReadAll(file)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy the file data to the server file
	os.WriteFile(fmt.Sprintf("%s/%s", uploadPath, fileName), videoData, os.ModePerm)

	// Respond with success
	w.WriteHeader(http.StatusOK)
	var videoId = videoid
	responseJson, err := json.Marshal(videoId)
	if err != nil {
		utils.SendError(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	w.Write(responseJson)

}

func setupRoutes(mux *http.ServeMux) {
	mux.Handle("/", authMiddleWare(http.HandlerFunc(homePage)))
	mux.HandleFunc("/ws", wsEndpoint)
	// mux.Handle("/startStream", authMiddleWare(http.HandlerFunc(startStream)))
	mux.Handle("/startStream", authMiddleWare(http.HandlerFunc(startStream)))
	mux.Handle("/uploadThumbnail", authMiddleWare(http.HandlerFunc(uploadThumbnail)))
	mux.Handle("/getVideos", (http.HandlerFunc(postgres.GetStreams)))
	mux.Handle("/getUserId", (http.HandlerFunc(postgres.GetUserId)))
	mux.Handle("/getContent", (http.HandlerFunc(postgres.GetContent)))
	mux.Handle("/getDashboardContent", authMiddleWare(http.HandlerFunc(postgres.GetDashboardContent)))

	mux.Handle("/getVideoData", getVideoDataMiddleware(http.HandlerFunc(postgres.GetVideoData)))
	mux.Handle("/like", authMiddleWare(http.HandlerFunc(postgres.Like)))
	mux.Handle("/dislike", authMiddleWare(http.HandlerFunc(postgres.Dislike)))
	mux.Handle("/removeLike", authMiddleWare(http.HandlerFunc(postgres.RemoveLike)))
	mux.Handle("/subscribe", authMiddleWare(http.HandlerFunc(postgres.Subscribe)))
	mux.Handle("/unsubscribe", authMiddleWare(http.HandlerFunc(postgres.Unsubscribe)))
	mux.Handle("/getLoggedUserDetails", http.HandlerFunc(auth.GetUserDetails))
	mux.HandleFunc("/getChats", postgres.GetChats)
	mux.Handle("/chat", (http.HandlerFunc(chat.Chat)))
	mux.Handle("/getCommmentsForChannel", authMiddleWare(http.HandlerFunc(postgres.GetCommmentsForCreator)))
	mux.Handle("/getUserDetailsByUsername", (http.HandlerFunc(postgres.GetUserDetailsByUsername)))
	mux.Handle("/getChannelSummary", authMiddleWare(http.HandlerFunc(postgres.GetChannelSummary)))
	mux.Handle("/updateUserDetails", authMiddleWare(http.HandlerFunc(postgres.UpdateUserDetails)))
	mux.Handle("/uploadVideo", authMiddleWare(http.HandlerFunc(uploadVideo)))
	mux.Handle("/uploadProfileImage", authMiddleWare(http.HandlerFunc(UploadProfileImage)))
	mux.Handle("/saveVod", authMiddleWare(http.HandlerFunc(postgres.SaveVod)))
	mux.HandleFunc("/getGoogleUrl", auth.GetGoogleUrl)
	mux.HandleFunc("/loginWithGoogle", auth.LoginWithGoogle)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/signup", auth.SignUp)
	mux.HandleFunc("/signin", auth.SignIn)
	mux.HandleFunc("/signOut", auth.SignOut)
	mux.Handle("/hls/", http.StripPrefix("/hls/", corsFileServer(http.Dir("/home/anurag/s3mnt"))))
}

func main() {
	fmt.Println("heelo staring go")
	mux := http.NewServeMux()
	if err := env.Load(".env"); err != nil {
		panic(err)
	}

	rmq.ConnectRMQ()
	defer rmq.CloseConnection()

	go rmq.ConsumeMessages("vods")
	postgres.Connect()
	defer postgres.Disconnect()
	go chat.HandleMessages()

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
		if req.Method == http.MethodOptions {
			// Handle preflight OPTIONS request
			w.Header().Set("Access-Control-Allow-Origin", env.Get("FRONTEND_URL", "https://streamvault.vercel.app"))
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.WriteHeader(http.StatusOK)
			return
		}

		// Set CORS headers for non-preflight requests
		w.Header().Set("Access-Control-Allow-Origin", env.Get("FRONTEND_URL", "https://streamvault.vercel.app"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Call the next handler in the chain
		next.ServeHTTP(w, req)
	})
}

func corsFileServer(fs http.FileSystem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		fileServer := http.FileServer(fs)
		fileServer.ServeHTTP(w, r)
	})
}
