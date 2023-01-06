package api

import (
	"DortgenAPI/src/database"
	"encoding/json"
	"errors"
	"net/http"
	"time"
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
	// get key from query string
	key := request.URL.Query().Get("key")
	// check if key is set
	if key == "" {
		// if not, return an error
		response := ValidateResponse{
			Success: false,
			Data: ValidateData{
				Error: "key not set",
			},
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
		_, err = writer.Write(responsePayload)
		return
	}

	valid, err := ValidateKey(key)
	if err != nil {
		response := ValidateResponse{
			Success: false,
			Data: ValidateData{
				Error: err.Error(),
			},
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
		_, err = writer.Write(responsePayload)
		return
	}

	// check if key is valid
	if !valid {
		// return invalid key
		response := ValidateResponse{
			Success: false,
			Data: ValidateData{
				Error: "invalid key",
			},
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
		_, err = writer.Write(responsePayload)
		return
	}
}

func ValidateKey(key string) (bool, error) {

	// check if key is in database
	apiKey, err := database.DoesKeyExist(key)
	if err != nil {
		return false, err
	}
	if apiKey == nil {
		return false, errors.New("invalid key")
	}

	if apiKey.Disabled {
		return false, errors.New("key is disabled")
	}

	if apiKey.LastGenerated+database.GenerateCooldown > time.Now().Unix() {
		return false, errors.New("key is on cooldown for " + time.Until(time.Unix(apiKey.LastGenerated+database.GenerateCooldown, 0)).Round(time.Second).String())
	}

	return true, nil
}
