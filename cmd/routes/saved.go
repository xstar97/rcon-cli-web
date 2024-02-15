package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"rcon-cli-web/config"
)

// HandleGetSavedData handles GET requests to retrieve saved data
func HandleGetSavedData(w http.ResponseWriter, r *http.Request) {
	log.Println("GET request received to retrieve saved data.")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Handle GET request to retrieve saved data
	savedData, err := config.ReadSavedDataFromFile()
	if err != nil {
		http.Error(w, "Failed to read saved data", http.StatusInternalServerError)
		return
	}

	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")
	// Encode the saved data to JSON and write it to the response writer
	json.NewEncoder(w).Encode(savedData)

	log.Println("Saved data sent successfully.")
}

// HandleUpdateSavedData handles POST requests to update saved data
func HandleUpdateSavedData(w http.ResponseWriter, r *http.Request) {
	log.Println("POST request received to update saved data.")

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Handle POST request to update saved data
	// Read the request body to get the updated saved data
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	// Decode the request body into a SavedData object
	var requestBody config.SavedData
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	// Update the saved data
	err = config.WriteSavedDataToFile(requestBody)
	if err != nil {
		http.Error(w, "Failed to save data", http.StatusInternalServerError)
		return
	}
	// Return an acknowledgment
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "Data received and saved successfully")

	log.Println("Saved data updated successfully.")
}

// HandleSaved handles GET and POST requests for saving and retrieving data
func HandleSaved(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling request for saving/retrieving data.")

	switch r.Method {
	case "GET":
		HandleGetSavedData(w, r)
	case "POST":
		HandleUpdateSavedData(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	log.Println("Request handled successfully.")
}
