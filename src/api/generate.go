package api

import (
	"DortgenAPI/src/database"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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
	// check to see if they are already requesting an account
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
		writer.WriteHeader(http.StatusTooManyRequests)
		_, err = writer.Write(responsePayload)
		return
	}
	// add them to the list of current requests
	addRequest(request.RemoteAddr)
	defer removeRequest(request.RemoteAddr)

	// get the key from query string
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
			log.Println("error marshalling generate response (key set):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write(responsePayload)
		return
	}

	// validate the key
	keyExists, err := database.Connection.DoesKeyExist(key)
	if err != nil {
		response := GenerateResponse{
			Success: false,
			Data: GenerateData{
				Error: err.Error(),
			},
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling generate response (validate key):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write(responsePayload)
		return
	}
	if !keyExists {
		response := GenerateResponse{
			Success: false,
			Data: GenerateData{
				Error: "invalid key",
			},
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling generate response (invalid key):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write(responsePayload)
		return
	}

	// check to see if key is disabled
	keyDisabled, err := database.Connection.IsKeyDisabled(key)
	if err != nil {
		response := GenerateResponse{
			Success: false,
			Data: GenerateData{
				Error: err.Error(),
			},
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling generate response (check key disabled):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write(responsePayload)
		return
	}
	if keyDisabled {
		response := GenerateResponse{
			Success: false,
			Data: GenerateData{
				Error: "key disabled",
			},
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling generate response (key disabled):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write(responsePayload)
		return
	}

	// check to see if cooldown is over
	cooldown, err := database.Connection.GetCooldown(key)
	if err != nil {
		log.Println("error getting cooldown:", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if cooldown > 0 {
		response := GenerateResponse{
			Success: false,
			Data: GenerateData{
				Error: "cooldown not over (" + strconv.Itoa(cooldown) + "s)",
			},
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling generate response (cooldown not over):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write(responsePayload)
		return
	}

	// check to see if there is stock
	stock, err := database.Connection.GetStockAmount()
	if err != nil {
		log.Println("error getting stock amount:", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if stock <= 0 {
		response := GenerateResponse{
			Success: false,
			Data: GenerateData{
				Error: "out of stock",
			},
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling generate response (out of stock):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
		_, err = writer.Write(responsePayload)
		return
	}

	// generate the account
	alt, err := database.Connection.GetAltAndRemoveFromStock()
	if err != nil {
		response := GenerateResponse{
			Success: false,
			Data: GenerateData{
				Error: err.Error(),
			},
		}
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Println("error marshalling generate response (get alt):", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusInternalServerError)
		_, err = writer.Write(responsePayload)
		return
	}

	// return the account
	response := GenerateResponse{
		Success: true,
		Data: GenerateData{
			Email:    alt.Email,
			Password: alt.Password,
			Combo:    alt.Email + ":" + alt.Password,
		},
	}
	responsePayload, err := json.Marshal(response)
	if err != nil {
		log.Println("error marshalling generate response (alt response):", err)
		writer.WriteHeader(http.StatusInternalServerError)
		err = database.Connection.AddAltToStock(alt.Email, alt.Password)
		if err != nil {
			log.Println("error adding alt back to stock:", err)
		}
		return
	}
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(responsePayload)

	// set cooldown for the key
	err = database.Connection.SetCooldown(key)

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
