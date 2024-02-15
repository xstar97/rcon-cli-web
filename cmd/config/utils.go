package config

import (
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
    "os/exec"
	"encoding/json"
	"log"
	"os"
	"github.com/fsnotify/fsnotify"
)

// represents the configuration for a server
type ServerConfig struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	Log      string `yaml:"log"`
	Type     string `yaml:"type"`
	Timeout  string `yaml:"timeout"`
}

// SavedData represents the saved data structure
type SavedData struct {
	Server string `json:"server"`
	Mode   string `json:"mode"`
}

// Function to read the YAML config file and return the content
func ReadConfig() (map[string]ServerConfig, error) {
	filePath := CONFIG.CLI_CONFIG

	// Log the file path
	log.Printf("Reading config from file: %s", filePath)

	// Check if the file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil, err
	}

	// Read YAML file
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal YAML into map
	var data map[string]ServerConfig
	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		return nil, err
	}

	// Create a new watcher to monitor changes to the file
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	defer watcher.Close()

	// Watch for changes to the file
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("Config file modified. Reloading config...")
					// Reload config from file
					reloadData, err := ReadConfig()
					if err != nil {
						log.Println("Error reloading config:", err)
						continue
					}
					log.Println("Config reloaded successfully")
					// Update the existing data with the reloaded data
					data = reloadData
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Error watching file:", err)
			}
		}
	}()

	// Add the file to the watcher
	err = watcher.Add(filePath)
	if err != nil {
		return nil, err
	}

	log.Println("Config read successfully from file.")
	return data, nil
}

// reads the YAML config file and returns the configuration for a specific server
func GetServer(serverName string) (ServerConfig, error) {
	data, err := ReadConfig()
	if err != nil {
		return ServerConfig{}, err
	}

	// Check if the server name exists
	config, ok := data[serverName]
	if !ok {
		return ServerConfig{}, fmt.Errorf("server '%s' not found", serverName)
	}

	return config, nil
}

// ExecuteShellCommand executes a shell command with provided arguments and returns its output
func ExecuteShellCommand(command string, args ...string) ([]byte, error) {
    // Set the command to execute
    cmd := exec.Command(command, args...)
    
    // Capture the output of the command
    output, err := cmd.CombinedOutput()
    if err != nil {
        return nil, err
    }

    return output, nil
}

// Function to read saved data from JSON file
func ReadSavedDataFromFile() (SavedData, error) {
	filePath := CONFIG.DB_JSON_FILE

	// Log the file path
	log.Printf("Reading data from file: %s", filePath)

	// Check if the file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// File doesn't exist, generate it with default values
		data := SavedData{
			Server: CONFIG.CLI_DEFAULT_SERVER,
			Mode:   CONFIG.MODE,
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

	// Create a new watcher to monitor changes to the file
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return SavedData{}, err
	}
	defer watcher.Close()

	// Watch for changes to the file
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("File modified. Reloading saved data...")
					// Reload saved data from file
					reloadData, err := ReadSavedDataFromFile()
					if err != nil {
						log.Println("Error reloading saved data:", err)
						continue
					}
					log.Println("Saved data reloaded:", reloadData)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Error watching file:", err)
			}
		}
	}()

	// Add the file to the watcher
	err = watcher.Add(filePath)
	if err != nil {
		return SavedData{}, err
	}

	log.Println("Data read successfully from file.")
	return data, nil
}

// Function to write saved data to JSON file
func WriteSavedDataToFile(data SavedData) error {
	filePath := CONFIG.DB_JSON_FILE

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