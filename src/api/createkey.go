package api

import (
	"DortgenAPI/src/database"
	"encoding/json"
	"log"
	"net/http"
)

type CreateKeyResponse struct {
	Success bool          `json:"success"`
	Data    CreateKeyData `json:"data,omitempty"`
}

type CreateKeyData struct {
	Error string `json:"error,omitempty"`
	Key   string `json:"key,omitempty"`
	Owner string `json:"owner,omitempty"`
}

type CreateKeyRequest struct {
	Owner string `json:"owner"`
}

var CreateKeyFunc = func(writer http.ResponseWriter, request *http.Request) {
	// get the key from query string
	key := request.URL.Query().Get("key")
	// check if key is set
	if key == "" {
		response := CreateKeyResponse{
			Success: false,
			Data: CreateKeyData{
				Error: "key not set",
			},
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling create key response (key not set):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write(responsePayload)
		return
	}

	// get the owner from post data
	var requestData CreateKeyRequest
	err := json.NewDecoder(request.Body).Decode(&requestData)
	if err != nil {
		log.Println("error decoding create key request:", err)
		response := CreateKeyResponse{
			Success: false,
			Data: CreateKeyData{
				Error: err.Error(),
			},
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling create key response (decode request):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write(responsePayload)
		return
	}
	// check if owner is set
	if requestData.Owner == "" {
		response := CreateKeyResponse{
			Success: false,
			Data: CreateKeyData{
				Error: "owner not set",
			},
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling create key response (owner not set):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write(responsePayload)
		return
	}

	// check if owner already exists
	ownerExists, err := database.Connection.DoesOwnerExist(requestData.Owner)
	if err != nil {
		log.Println("error checking if owner exists:", err)
		response := CreateKeyResponse{
			Success: false,
			Data: CreateKeyData{
				Error: err.Error(),
			},
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling create key response (check if owner exists):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusInternalServerError)
		_, err = writer.Write(responsePayload)
		return
	}
	if ownerExists {
		// get key
		key, err := database.Connection.GetKeyFromOwner(requestData.Owner)
		if err != nil {
			log.Println("error getting key from owner:", err)
			response := CreateKeyResponse{
				Success: false,
				Data: CreateKeyData{
					Error: err.Error(),
				},
			}
			responsePayload, err := json.Marshal(response)
			if err != nil {
				log.Println("error marshalling create key response (get key from owner):", err)
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
			writer.WriteHeader(http.StatusInternalServerError)
			_, err = writer.Write(responsePayload)
			return
		}
		response := CreateKeyResponse{
			Success: true,
			Data: CreateKeyData{
				Key:   key,
				Owner: requestData.Owner,
			},
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling create key response (owner exists):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
		_, err = writer.Write(responsePayload)
		return
	}

	// create key
	err = database.Connection.CreateApiKey(requestData.Owner, 12)
	if err != nil {
		log.Println("error creating api key:", err)
		response := CreateKeyResponse{
			Success: false,
			Data: CreateKeyData{
				Error: err.Error(),
			},
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling create key response (create key):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusInternalServerError)
		_, err = writer.Write(responsePayload)
		return
	}

	// get key
	key, err = database.Connection.GetKeyFromOwner(requestData.Owner)
	if err != nil {
		log.Println("error getting key from owner:", err)
		response := CreateKeyResponse{
			Success: false,
			Data: CreateKeyData{
				Error: err.Error(),
			},
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling create key response (get key from owner):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusInternalServerError)
		_, err = writer.Write(responsePayload)
		return
	}

	response := CreateKeyResponse{
		Success: true,
		Data: CreateKeyData{
			Key:   key,
			Owner: requestData.Owner,
		},
	}
	responsePayload, err := json.Marshal(response)
	if err != nil {
		log.Println("error marshalling create key response:", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(responsePayload)
	if err != nil {
		log.Println("error writing create key response:", err)
	}
	return
}
