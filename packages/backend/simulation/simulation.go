package simulation

import (
	"fmt"
	"math/rand/v2"
	"os"
	"os/exec"
	"streamvault/postgres"
	"time"

	"github.com/gofor-little/env"
	"github.com/goombaio/namegenerator"
)

var uuids []string
var userIds []string
var backendUrl = env.Get("BACKEND_URL", "http://localhost:8080")

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

func StartSimulation() {
	fmt.Println("simulation started")
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	fmt.Println("making user")

	fmt.Println("adding stream")
	for i := 0; i < 5; i++ {
		name := nameGenerator.Generate()
		userId, err := postgres.CreateUserWithPassword(name, "test User")
		if err != nil {
			fmt.Println("error adding user", err.Error())
			continue
		}
		userIds = append(userIds, userId)

	}

	fmt.Println("userids len", len(userIds))

	for _, stream := range streams {

		go func() {
			var r = rand.IntN(len(userIds))
			fmt.Println("r", r)
			userId := userIds[r]

			inputPath := fmt.Sprintf("/home/anurag/s3mnt/simulation/%s", stream.FileName)

			id, err := postgres.AddStream(stream.Title, stream.Description, stream.Category, stream.Thumbnail, userId)
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
			cmd := exec.Command("ffmpeg",

				"-i", inputPath,
				"-loop", "1",
				"-c:v", "libx264", "-preset", "ultrafast", "-tune", "zerolatency",
				"-c:a", "aac", "-ar", "44100", "-b:a", "64k",
				"-f", "hls",
				"-g", "20",
				"-hls_time", "5",
				"-hls_list_size", "50",
				"-hls_flags", "delete_segments",
				// "-progress", "pipe:1",
				fmt.Sprintf("%s/%s.m3u8", dirPath, id),
			)
			if err := cmd.Start(); err != nil {
				fmt.Println("Error starting command:", err)
				return
			}

			if err := cmd.Wait(); err != nil {
				fmt.Println("Error waiting for command:", err)
				return
			}
		}()
	}

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
