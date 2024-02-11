package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"rcon-cli-web/config"
)

// SavedData represents the saved data structure
type SavedData struct {
	Server string `json:"server"`
	Mode   string `json:"mode"`
}

// Function to read saved data from JSON file
func ReadSavedDataFromFile() (SavedData, error) {
	filePath := config.CONFIG.DB_JSON_FILE

	// Log the file path
	log.Printf("Reading data from file: %s", filePath)

	// Check if the file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// File doesn't exist, generate it with default values
		data := SavedData{
			Server: config.CONFIG.CLI_DEFAULT_SERVER,
			Mode:   config.CONFIG.MODE,
		}
		err := WriteSavedDataToFile(data)
		if err != nil {
			return SavedData{}, err
		}
		log.Println("File does not exist. Generated file with default values.")
		return data, nil
	}

	log.Println("File exists. Reading data from file.")

	// File exists, read data from the file
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return SavedData{}, err
	}

	// Unmarshal the file content into SavedData struct
	var data SavedData
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		return SavedData{}, err
	}

	log.Println("Data read successfully from file.")
	return data, nil
}

// Function to write saved data to JSON file
func WriteSavedDataToFile(data SavedData) error {
	filePath := config.CONFIG.DB_JSON_FILE

	// Log the file path
	log.Printf("Writing data to file: %s", filePath)

	// Check if the file exists and is writable
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// If the file doesn't exist, try to create it
		log.Println("File does not exist. Creating file...")
		file, err := os.Create(filePath)
		if err != nil {
			log.Printf("Error creating file: %v", err)
			return err
		}
		defer file.Close()
		log.Println("File created successfully.")
	} else if err != nil {
		log.Printf("Error checking file status: %v", err)
		return err
	}

	// Marshal the data into JSON
	fileContent, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("Error marshaling data to JSON: %v", err)
		return err
	}

	// Write the JSON data to the file
	err = ioutil.WriteFile(filePath, fileContent, 0644)
	if err != nil {
		log.Printf("Error writing data to file: %v", err)
		return err
	}

	log.Println("Data successfully written to file.")
	return nil
}

// HandleGetSavedData handles GET requests to retrieve saved data
func HandleGetSavedData(w http.ResponseWriter, r *http.Request) {
	log.Println("GET request received to retrieve saved data.")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
	var requestBody SavedData
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	// Update the saved data
	err = WriteSavedDataToFile(requestBody)
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
