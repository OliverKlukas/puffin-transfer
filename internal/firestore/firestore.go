package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"os"
	"strings"
)

var client *firestore.Client

func Initialize(projectID, credsFile string) error {
	// Connect to GCP Firestore.
	c, err := firestore.NewClient(context.Background(), projectID,
		option.WithCredentialsFile(credsFile))
	if err != nil {
		return fmt.Errorf("connection to Firestore failed: %v", err)
	}
	client = c
	return nil
}

func Transfer(filepath string) error {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("could not read %s", filepath)
	}
	_, _, err = client.Collection("files").Add(context.Background(), map[string]interface{}{
		"filepath": filepath,
		"content":  file,
	})
	if err != nil {
		return fmt.Errorf("could not transfer %s to db", filepath)
	}
	return nil
}

func Retrieve(filepath string) (string, error) {
	q := client.Collection("files").Where("filepath", "==", filepath)
	iter := q.Documents(context.Background())
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			return "", fmt.Errorf("could not retrieve %s from db", filepath)
		}
		if err != nil {
			return "", fmt.Errorf("could not retrieve %s from db", filepath)
		}
		data := doc.Data()
		content, ok := data["content"].([]byte)
		if !ok {
			return "", fmt.Errorf("could not retrieve %s from db", filepath)
		}
		return string(content), nil
	}
}

func Run(inputChan chan string) {
	for {
		select {
		case cmd := <-inputChan:
			parts := strings.Split(cmd, " ")
			if len(parts) < 2 {
				fmt.Println("Invalid command!")
				fmt.Println("Usage: store <transfer|retrieve> <filepath>")
				continue
			}
			switch parts[0] {
			case "transfer":
				err := Transfer(parts[1])
				if err != nil {
					fmt.Printf("Error transferring file: %v\n", err)
				}
			case "retrieve":
				content, err := Retrieve(parts[1])
				if err != nil {
					fmt.Printf("Error retrieving file: %v\n", err)
				} else {
					fmt.Printf("File content: %s\n", content)
				}
			default:
				fmt.Println("Invalid command!")
				fmt.Println("Usage: store <transfer|retrieve> <filepath>")
			}
		}
	}
}
