package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// TODO:
type errorResponse struct {
	Error string `json:"error"`
}

func renderErrorResponse(rw http.ResponseWriter, msg string, status int) {
	renderResponse(rw, errorResponse{Error: msg}, status)
}

func renderResponse(rw http.ResponseWriter, res interface{}, status int) {
	rw.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(res)
	if err != nil {
		rw.WriteHeader(status)
		return
	}

	rw.WriteHeader(status)

	if _, err := rw.Write(data); err != nil {
		//TODO: implement this
		fmt.Println(err)
	}
}
