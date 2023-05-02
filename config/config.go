package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DatabaseType             string
	FirestoreProjectID       string
	FirestoreCredentialsFile string
}

func createConfigFile() (Config, error) {
	config := Config{}

	// Prompt user to choose database type
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Choose database type:")
	fmt.Println("1. Firestore")

	var choice int
	_, err := fmt.Fscan(reader, &choice)
	if err != nil {
		return config, err
	}

	switch choice {
	case 1:
		config.DatabaseType = "Firestore"

		// Prompt user for Firestore project ID
		fmt.Print("Enter Firestore project ID: ")
		projectID, err := reader.ReadString('\n')
		if err != nil {
			return config, err
		}
		config.FirestoreProjectID = projectID[:len(projectID)-1] // remove newline character

		// Prompt user for service account key file location
		fmt.Print("Enter service account authentication key JSON file location: ")
		keyFileLocation, err := reader.ReadString('\n')
		if err != nil {
			return config, err
		}
		config.FirestoreCredentialsFile = keyFileLocation[:len(keyFileLocation)-1] // remove newline character
	default:
		fmt.Println("Invalid choice. Please choose a valid option.")
		return createConfigFile()
	}

	// Write config to file
	file, err := os.Create("config.json")
	if err != nil {
		return config, err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func GetConfig() (Config, error) {
	config, err := LoadConfig()
	if err != nil {
		return config, err
	}

	if config.DatabaseType != "Firestore" {
		return config, fmt.Errorf("unsupported database type: %s", config.DatabaseType)
	}

	return config, nil
}

func LoadConfig() (Config, error) {
	config := Config{}

	// Check if config file exists
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		// Config file does not exist, create it
		newConfig, err := createConfigFile()
		if err != nil {
			return config, err
		}
		config = newConfig
	} else {
		// Config file exists, load it
		file, err := os.Open("config.json")
		if err != nil {
			return config, err
		}
		defer file.Close()

		err = json.NewDecoder(file).Decode(&config)
		if err != nil {
			return config, err
		}
	}

	return config, nil
}
