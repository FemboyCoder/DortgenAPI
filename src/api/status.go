package api

import (
	"DortgenAPI/src/database"
	"encoding/json"
	"log"
	"net/http"
)

type StatusResponse struct {
	Stock int `json:"stock"`
}

var StatusFunc = func(writer http.ResponseWriter, request *http.Request) {
	// get stock amount

	stock, err := database.Connection.GetStockAmount()
	if err != nil {
		log.Println("Error getting stock amount: " + err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := StatusResponse{
		Stock: stock,
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
