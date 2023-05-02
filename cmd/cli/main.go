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
	firestoreInputChan := make(chan string)
	go firestore.Run(firestoreInputChan)

	// Start the command loop.
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

		switch {
		case strings.HasPrefix(input, "store"):
			// handle "store" command
			parts := strings.Fields(input)
			if len(parts) != 3 {
				fmt.Println("Usage: store <transfer|retrieve> <filepath>")
				continue
			}
			firestoreInputChan <- strings.Join(parts[1:], " ")
		case strings.HasPrefix(input, "help"):
			// handle "help" command
			fmt.Println("Available commands:")
			fmt.Println("\thelp\t\t\t\t\tPrints this help message")
			fmt.Println("\tstore <command> <filepath>\t\tTransfer or retrieve a file from store")
		default:
			// handle invalid command
			fmt.Println("Invalid command!")
			fmt.Println("Type 'help' to see available commands.")
		}
	}
}
