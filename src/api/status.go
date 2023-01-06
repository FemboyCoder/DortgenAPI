package api

import (
	"encoding/json"
	"net/http"
)

type StatusResponse struct {
	Stock int `json:"stock"`
}

var StatusFunc = func(writer http.ResponseWriter, request *http.Request) {
	// get stock amount
	response := StatusResponse{
		Stock: 100,
	}

	responsePayload, err := json.Marshal(response)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(responsePayload)
	if err != nil {
		return
	}
}
