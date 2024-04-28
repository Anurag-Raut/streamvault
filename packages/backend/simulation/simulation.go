package simulation

import (
	"fmt"
	"os"
	"os/exec"
	"streamvault/postgres"
)

var uuids []string

var files = []string{"file1.mp4"}

func StartSimulation() {
	fmt.Println("simulation started")

	


	inputPath := fmt.Sprintf("/home/anurag/s3mnt/simulation/%s", files[0])
	fmt.Println("making user")

	userId, err := postgres.CreateUserWithPassword("test User", "test User")

	fmt.Println("user made", userId)

	// os.RemoveAll(dirPath)

	if err != nil {
		fmt.Println("error hello", err.Error())
		return
	}

	fmt.Println("adding stream")
	id,err:=postgres.AddStream("test1", "description", "cat1", "thumbnail", userId)
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
