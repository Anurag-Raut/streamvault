package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofor-little/env"
)

func SendError(w http.ResponseWriter, err string, code int) {
	type ErrorResponse struct {
		Error string `json:"error"`
	}
	// w.WriteHeader(http.StatusInternalServerError)
	// w.Write([]byte(err.Error()))
	w.WriteHeader(code)
	errObj := ErrorResponse{Error: err}
	errResp, _ := json.MarshalIndent(errObj, "", "  ")

	w.Write([]byte(errResp))


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
	fmt.Println("sending to subtitler ")
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/receive_text", env.Get("SUBTITLER_API_URL", "http://localhost:5000")), bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}

	fmt.Println("done sennding")
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