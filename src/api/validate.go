package api

import (
	"DortgenAPI/src/database"
	"encoding/json"
	"log"
	"net/http"
)

type ValidateResponse struct {
	Success bool         `json:"success"`
	Data    ValidateData `json:"data,omitempty"`
}

type ValidateData struct {
	Error string `json:"error,omitempty"`
	Valid string `json:"valid,omitempty"`
}

var ValidateFunc = func(writer http.ResponseWriter, request *http.Request) {

	// get the key from query string
	key := request.URL.Query().Get("key")
	// check if key is set
	if key == "" {
		response := ValidateResponse{
			Success: false,
			Data: ValidateData{
				Error: "key not set",
			},
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling validate response (key not set):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write(responsePayload)
		return
	}

	// validate the key
	valid, err := database.Connection.DoesKeyExist(key)
	if err != nil {
		log.Println("error validating key:", err)
		response := ValidateResponse{
			Success: false,
			Data: ValidateData{
				Error: err.Error(),
			},
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling validate response (validate key):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write(responsePayload)
		return
	}

	if !valid {
		response := ValidateResponse{
			Success: false,
			Data: ValidateData{
				Error: "invalid key",
			},
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling validate response (invalid key):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write(responsePayload)
		return
	}

	// check if the key is disabled or not

	disabled, err := database.Connection.IsKeyDisabled(key)
	if err != nil {
		log.Println("error validating key:", err)
		response := ValidateResponse{
			Success: false,
			Data: ValidateData{
				Error: err.Error(),
			},
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling validate response (validate key):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write(responsePayload)
		return
	}

	if disabled {
		response := ValidateResponse{
			Success: false,
			Data: ValidateData{
				Error: "key disabled",
			},
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling validate response (invalid key):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write(responsePayload)
		return
	}

	// key is valid and not disabled
	response := ValidateResponse{
		Success: true,
		Data: ValidateData{
			Valid: "true",
		},
	}
	responsePayload, err := json.Marshal(response)
	if err != nil {
		log.Println("error marshalling validate response (valid key):", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(responsePayload)
	return
}
