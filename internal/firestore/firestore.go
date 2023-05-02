package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"os"
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

func Run(input chan string, out chan string) {
	for {
		select {
		case cmd := <-input:
			if cmd == "transfer" {
				filepath := <-filepathChan
				err := Transfer(filepath)
				errorChan <- err
			} else if cmd == "retrieve" {
				filepath := <-filepathChan
				content, err := Retrieve(filepath)
				if err != nil {
					errorChan <- err
				} else {
					fmt.Println(content)
					errorChan <- nil
				}
			} else {
				errorChan <- fmt.Errorf("invalid command: %s", cmd)
			}
		default:
			// Do nothing and wait for new commands
		}
	}
}
