package simulation

import (
	"bufio"
	// "bytes"
	"context"
	"fmt"
	"math/rand/v2"
	// "net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"streamvault/chat"
	"streamvault/postgres"
	"streamvault/utils"
	"time"

	cohere "github.com/cohere-ai/cohere-go/v2"
	cohereclient "github.com/cohere-ai/cohere-go/v2/client"
	"github.com/gofor-little/env"
	"github.com/goombaio/namegenerator"
)

func generateChat(username string, videoId string) (string, error) {
	videodata, err := postgres.GetVideoDataFromDatabase(videoId)
	if err != nil {
		fmt.Println("Erorr: ",err.Error())
		return "",err
	}
	var chats []postgres.Chat = postgres.GetChatsFromDatabase(videoId, 10)
	chatString := "{"
	for _, chat := range chats {
		fmt.Println(chat.Message,"asdasd",chat.User.Username)
		chatString += fmt.Sprintf("[Username: %s, Chat: %s]", chat.User.Username, chat.Message)
	}
	chatString += "}"
	
	client := cohereclient.NewClient(cohereclient.WithToken(env.Get("COHERE_TOKEN", "")))
	fmt.Println("reached",env.Get("COHERE_TOKEN", ""))
	response, err := client.Chat(
		context.TODO(),
		&cohere.ChatRequest{
			Message: fmt.Sprintf(`You are a viewer in a stream .
		  Your username is : %s
		  Here are the details of the stream:
		  (
			Title:%s,
			Description:%s
			Category:%s
			streamer's Username:%s
		  )
		  /n

		  these are the previous chats -
		  %s.
		  \n

		 


		  now your task is to write small chat with the other users in the chat or talk something related to stream.
		  chat should be short (1-10) words, use words that chatters use on twitch and youtube and intersting and like real humans on internet .
		  you can  comment about video , ask related questions , address other chatter , or any other thing , dont do all in single chat only one thing make sure to keep it short like about 5 to 10 words.
		  do not repeat messages , make new chats
		  only return the text, dont give username of the chatter
		  output format : chat
		  `, username, videodata.Title, videodata.Description, videodata.Category, videodata.User.Username, chatString),
		},
	)

	if err != nil {
		fmt.Println("Error occued incoher token : ", err.Error())
		return "", err
	}
	
	fmt.Println("response", response.Text)
	return response.Text, nil

}

var uuids []string
var userIds []string
var backendUrl = env.Get("BACKEND_URL", "https://api.echon.live")

var users []postgres.UserDetails

var videoIds []string

type Stream struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Thumbnail   string `json:"thumbnail"`
	UserId      string `json:"userId"`
	FileName    string `json:"fileName"`
}

var streams = []Stream{
	{
		Title:       "test1",
		Description: "description",
		Category:    "cat1",
		Thumbnail:   fmt.Sprintf("%s/hls/thumbnail/lTRiuFIWV54-HD.jpg", backendUrl),

		FileName: "file1.mp4",
	},
	{
		Title: "How AI Was Stolen",

		Description: "description",
		Category:    "cat1",
		Thumbnail:   fmt.Sprintf("%s/hls/thumbnail/BQTXv5jm6s4-HD.jpg", backendUrl),

		FileName: "file2.mp4",
	},
	{
		Title: "Interview with Senior Rust Developer",

		Description: "description",
		Category:    "cat1",
		Thumbnail:   fmt.Sprintf("%s/hls/thumbnail/TGfQu0bQTKc-HD.jpg", backendUrl),

		FileName: "file3.mp4",
	},
	{
		Title: "The Only Database Abstraction You Need | Prime Reacts",

		Description: "description",
		Category:    "cat1",
		Thumbnail:   fmt.Sprintf("%s/hls/thumbnail/nWchov5Do-o-HD.jpg", backendUrl),

		FileName: "file4.mp4",
	},
	{
		Title: "Ludwig and Squeex conquer Elden Ring and the Wheel of Punishment (Day 1 - Part 2)",

		Description: "description",
		Category:    "cat1",
		Thumbnail:   fmt.Sprintf("%s/hls/thumbnail/S4c1KAI81CE-HD.jpg", backendUrl),

		FileName: "file5.mp4",
	},
}


func deleteStreamsByTitle(){
	for _,stream:=range(streams) {
		fmt.Println("deleteing ")
		postgres.DeleteStreamsByTitle(stream.Title)
	}

}

func StartSimulation() {
	fmt.Println("simulation started")
	seed := time.Now().UTC().UnixNano()	
	nameGenerator := namegenerator.NewNameGenerator(seed)

	fmt.Println("making user")
	deleteStreamsByTitle()

	fmt.Println("adding stream")
	for i := 0; i < 5; i++ {
		name := nameGenerator.Generate()
		userId, err := postgres.CreateUserWithPassword(name, "test User")
		if err != nil {
			fmt.Println("error adding user", err.Error())
			continue
		}
		userIds = append(userIds, userId)
		userDetails, err := postgres.GetUserDetailsFromDatabase(userId)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		users = append(users, userDetails)

	}

	fmt.Println("userids len", len(userIds))

	for _, stream := range streams {

		go func() {
			var r = rand.IntN(len(userIds))
			fmt.Println("r", r)
			userId := userIds[r]

			inputPath := fmt.Sprintf("/home/anurag/s3mnt/simulation/%s", stream.FileName)

			id, err := postgres.AddStream(stream.Title, stream.Description, stream.Category, stream.Thumbnail, userId)
			videoIds = append(videoIds, id)
			fmt.Println("videoId", len(videoIds))
			if err != nil {
				fmt.Println("error adding stream", err.Error())
				return

			}
			uuids = append(uuids, id)
			dirPath := fmt.Sprintf("/home/anurag/s3mnt/%s", id)

			_, osserr := os.Stat(dirPath)

			if os.IsNotExist(osserr) {
				err := os.Mkdir(dirPath, 0777)
				if err != nil {
					fmt.Println("Error creating directory:", err)
					return
				}
				println("Directory created")
			}

			fmt.Println("stream added")
			fmt.Printf("%s/%s.m3u8", dirPath, id)
			fmt.Println(inputPath)
			// responseJson := fmt.Sprintf(`{"streamId": "%s"}`, id)
			// _, err = http.Post(fmt.Sprintf("%s/start_transcription", env.Get("SUBTITLER_API_URL", "http://subtitler:5000")), "application/json", bytes.NewBuffer([]byte(responseJson)))
			// if err != nil {
			// 	fmt.Println("Error start_transcription:", err)

			// 	return 
			// }
			fmt.Println("YAYYY")
				cmd := exec.Command("ffmpeg",

				"-i", inputPath,
				"-loop", "1",
				"-c:v", "libx264", "-preset", "veryfast", "-tune", "zerolatency",
				"-c:a", "aac", "-ar", "44100", "-b:a", "64k",
				"-f", "hls",
				"-g", "20",
				"-hls_time", "5",
				"-hls_list_size", "500",
				"-hls_flags", "delete_segments",
				// "-progress", "pipe:1",
				fmt.Sprintf("%s/%s.m3u8", dirPath, id),
			)



			


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
				
				postgres.UpdateStatus(id, false)
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
						fmt.Println("filepath",filePath)
						parts := strings.Split(filePath, "/")
						time.Sleep(time.Millisecond*200)
		
						// Get the file name from the file path
						if strings.HasSuffix(filePath, ".ts") && !strings.HasSuffix(filePath, ".m3u8.tmp") {
							cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", filePath)
							output, err := cmd.CombinedOutput()
							if err != nil {
								// fmt.Println("Error executing ffprobe command:", err.Error())
								// return
								continue
							}
		
							durSting := string(output)
							// fmt.Println(dur)
							durSting = strings.TrimSuffix(durSting, "\n")
							duration, err := strconv.ParseFloat(durSting, 64)
		
							if err != nil {
								fmt.Println("Error converting duration:", err)
		
								return
							}
		
							fmt.Println("Duration:", duration)
							var a = strings.Split(parts[len(parts)-1], ".")[0]
							segmentNumberString := strings.TrimPrefix(a, id)
							segmentNumber, err := strconv.Atoi(segmentNumberString)
							if err != nil {
								fmt.Println("Error converting segment number:", err)
								return
							}
		
							fmt.Println("Segment Number:", segmentNumber)
		
							err = utils.SendToSubtitler(parts[len(parts)-2]+"/"+parts[len(parts)-1], id, duration, totalDuration, segmentNumber)
							if err != nil {
								fmt.Println("Error sending to subtitler:", err)
								return
							}
							totalDuration += duration
		
						}
		
					}
				}
			}()


			if err := cmd.Wait(); err != nil {
				fmt.Println("Error waiting for command:", err)
				return
			}
		}()
	}
	time.Sleep(10 * time.Second)
	go StartChatBots()

}

func StopSimulation() {
	fmt.Println("simulation stopped")
	for _, uuid := range uuids {
		dirPath := fmt.Sprintf("/home/anurag/s3mnt/%s", uuid)
		err := os.RemoveAll(dirPath)
		if err != nil {
			fmt.Println("Error removing directory:", err)
		}
	}
}

func StartChatBots() {
	fmt.Println("adsadsadasdsad", len(videoIds))

	for {

		for _, videoId := range videoIds {
			fmt.Println("Hightassad")
			var r = rand.IntN(len(users))
			var randomUser postgres.UserDetails = users[r]

			var msg chat.Message

			msg.User = randomUser
			msg.StreamId = videoId
			fmt.Println(randomUser.Username,videoId,"adasd")

			chatMessage, err := generateChat(randomUser.Username, videoId)
			if err != nil {
				time.Sleep(1000 * time.Second)
				continue
			}
			msg.Message = chatMessage

			fmt.Println(msg.User.Username, msg.Message)

			chat.Broadcast <- msg

			time.Sleep(3 * time.Second)

			// chat.Broadcast

		}
	}

}
