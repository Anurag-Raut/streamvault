package utils

import (
	"encoding/json"
	"net/http"

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