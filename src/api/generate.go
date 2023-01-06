package api

import (
	"DortgenAPI/src/database"
	"encoding/json"
	"net/http"
)

type GenerateResponse struct {
	Success bool         `json:"success"`
	Data    GenerateData `json:"data,omitempty"`
}

type GenerateData struct {
	Error    string `json:"error,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Combo    string `json:"combo,omitempty"`
}

var (
	CurrentRequests = map[string]struct{}{}
)

var GenerateFunc = func(writer http.ResponseWriter, request *http.Request) {

	if isRequesting(request.RemoteAddr) {
		response := GenerateResponse{
			Success: false,
			Data: GenerateData{
				Error: "already requesting",
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
	
	addRequest(request.RemoteAddr)
	defer removeRequest(request.RemoteAddr)

	// get key from query string
	key := request.URL.Query().Get("key")
	// check if key is set
	if key == "" {
		// if not, return an error
		response := GenerateResponse{
			Success: false,
			Data: GenerateData{
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
		response := GenerateResponse{
			Success: false,
			Data: GenerateData{
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
		// if not, return an error
		response := GenerateResponse{
			Success: false,
			Data: GenerateData{
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

	err = database.UpdateCooldown(key)
	if err != nil {
		response := GenerateResponse{
			Success: false,
			Data: GenerateData{
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

	// get stock amount
	response := GenerateResponse{
		Success: true,
		Data: GenerateData{
			Email:    "testemail@hotmail.com",
			Password: "password@123!",
			Combo:    "testemail@hotmail.com:password@123!",
		},
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

func addRequest(ip string) {
	CurrentRequests[ip] = struct{}{}
}

func removeRequest(ip string) {
	delete(CurrentRequests, ip)
}

func isRequesting(ip string) bool {
	_, ok := CurrentRequests[ip]
	return ok
}