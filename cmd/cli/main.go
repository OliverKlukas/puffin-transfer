package main

import (
	"bufio"
	"fmt"
	"go-tui-file-project/config"
	"go-tui-file-project/internal/firestore"
	"os"
	"strings"
)

func main() {
	fmt.Println("Welcome to Puffin Transfer Project!")

	// Get local config.
	conf, err := config.GetConfig()
	if err != nil {
		fmt.Println("Error getting conf:", err)
		os.Exit(1)
	}

	// Configure the firestore service.
	err = firestore.Initialize(conf.FirestoreProjectID, conf.FirestoreCredentialsFile)
	if err != nil {
		fmt.Println("Error initializing firestore service:", err)
		os.Exit(1)
	}

	// Run Firestore in the background.
	cmdChan := make(chan string)
	filepathChan := make(chan string)
	errorChan := make(chan error)
	go firestore.Run(cmdChan, filepathChan, errorChan)

	// Start the command loop
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if len(input) == 0 {
			continue
		}

		if strings.HasPrefix(input, "store") {
			parts := strings.Fields(input)
			if len(parts) != 3 {
				fmt.Println("Usage: store <transfer|retrieve> <filepath>")
				continue
			}

			if parts[1] == "transfer" {
				filepathChan <- parts[2]
				err := <-errorChan
				if err != nil {
					fmt.Println("Failed to transfer file:", err)
				} else {
					fmt.Println("File transferred successfully")
				}
			} else if parts[1] == "retrieve" {
				filepathChan <- parts[2]
				err := <-errorChan
				if err != nil {
					fmt.Println("Failed to retrieve file:", err)
				} else {
					fmt.Println("File transferred successfully")
				}
			} else {
				fmt.Println("Usage: store <transfer|retrieve> <filepath>")
			}
		} else {
			fmt.Println("Invalid command")
		}
	}
}
