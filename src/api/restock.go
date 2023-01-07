package api

import (
	"DortgenAPI/src/database"
	"encoding/json"
	"log"
	"mime/multipart"
	"net/http"
)

type RestockResponse struct {
	Succes  bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

var RestockFunc = func(writer http.ResponseWriter, request *http.Request) {

	// get the key from query string
	key := request.URL.Query().Get("key")
	// check if key is set
	if key == "" {
		response := RestockResponse{
			Succes:  false,
			Message: "key not set",
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling restock response (key not set):", err)
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
		response := RestockResponse{
			Succes:  false,
			Message: err.Error(),
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling restock response (validate key):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write(responsePayload)
		return
	}

	// check if key is valid
	if !valid {
		response := RestockResponse{
			Succes:  false,
			Message: "invalid key",
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling restock response (invalid key):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write(responsePayload)
		return
	}

	// check if key is disabled or not
	disabled, err := database.Connection.IsKeyDisabled(key)
	if err != nil {
		log.Println("error checking if key is disabled:", err)
		response := RestockResponse{
			Succes:  false,
			Message: err.Error(),
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling restock response (check if key is disabled):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write(responsePayload)
		return
	}

	// check if key is disabled
	if disabled {
		response := RestockResponse{
			Succes:  false,
			Message: "key is disabled",
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling restock response (key is disabled):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write(responsePayload)
		return
	}

	// check if key owner is admin
	owner, err := database.Connection.GetOwnerFromKey(key)
	if err != nil {
		log.Println("error getting key owner:", err)
		response := RestockResponse{
			Succes:  false,
			Message: err.Error(),
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling restock response (get key owner):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write(responsePayload)
		return
	}

	// check if key owner is admin
	if owner != "admin" {
		response := RestockResponse{
			Succes:  false,
			Message: "key is not admin",
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling restock response (key is not admin):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write(responsePayload)
		return
	}

	// parse the form
	err = request.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Println("error parsing multipart form:", err)
		response := RestockResponse{
			Succes:  false,
			Message: "error parsing multipart form",
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling restock response (parse form):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write(responsePayload)
		return
	}

	// check if a file was sent
	if request.MultipartForm == nil {
		response := RestockResponse{
			Succes:  false,
			Message: "no file sent",
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling restock response (no file sent):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write(responsePayload)
		if err != nil {
			log.Println("error writing restock response:", err)
			return
		}
		return
	}

	// read the file sent in the request
	file, fileHeader, err := request.FormFile("altfile")
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		response := RestockResponse{
			Succes:  false,
			Message: "error reading file",
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling restock response (error reading file):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = writer.Write(responsePayload)
		if err != nil {
			log.Println("error writing restock response:", err)
			return
		}
		return
	}
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)

	// add the accounts to the database
	response, err := database.Connection.AddAccountsFromFile(file, fileHeader.Size)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		response := RestockResponse{
			Succes:  false,
			Message: "error adding accounts to database",
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling restock response (error adding accounts to database):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = writer.Write(responsePayload)
		if err != nil {
			log.Println("error writing restock response:", err)
			return
		}
		return
	}

	// return a success response
	restockResponse := RestockResponse{
		Succes:  true,
		Message: response,
	}
	responsePayload, err := json.Marshal(restockResponse)
	if err != nil {
		log.Println("error marshalling restock response (success):", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(responsePayload)
	if err != nil {
		log.Println("error writing restock response:", err)
		return
	}
}
