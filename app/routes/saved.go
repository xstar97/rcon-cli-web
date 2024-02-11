package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"rcon-cli-web/config"
)

// Define a struct to represent the saved data
type SavedData struct {
	Server string `json:"server"`
	Mode   string `json:"mode"`
}

// Function to read saved data from JSON file
func ReadSavedDataFromFile() (*SavedData, error) {
	filePath := config.CONFIG.DB_JSON_FILE

	// Check if the file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// File doesn't exist, generate it with default values
		data := &SavedData{
			Server: config.CONFIG.CLI_DEFAULT_SERVER,
			Mode:   config.CONFIG.MODE,
		}
		err := WriteSavedDataToFile(data)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	// File exists, read data from the file
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the file content into SavedData struct
	data := &SavedData{}
	err = json.Unmarshal(fileContent, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Function to save data to JSON file
func SaveDataToFile(data *SavedData) error {
	filePath := config.CONFIG.DB_JSON_FILE
	fileContent, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filePath, fileContent, 0644)
	if err != nil {
		return err
	}
	return nil
}

// HandleSaved handles GET and POST requests for saving and retrieving data
func HandleSaved(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Handle GET request to retrieve saved data
		savedData, err := ReadSavedDataFromFile()
		if err != nil {
			http.Error(w, "Failed to read saved data", http.StatusInternalServerError)
			return
		}
		// Set the content type to application/json
		w.Header().Set("Content-Type", "application/json")
		// Encode the saved data to JSON and write it to the response writer
		json.NewEncoder(w).Encode(savedData)

	case "POST":
		// Handle POST request to update saved data
		// Read the request body to get the updated saved data
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		// Decode the request body into the struct
		var requestBody SavedData
		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}
		// Save the updated data
		err = SaveDataToFile(&requestBody)
		if err != nil {
			http.Error(w, "Failed to save data", http.StatusInternalServerError)
			return
		}
		// Return a success message
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Saved data updated successfully")

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
