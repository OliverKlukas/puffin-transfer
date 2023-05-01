package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"os"
)

func main() {
	// Connect to GCP firestore.
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "puffin-transfer",
		option.WithCredentialsFile("/home/olli/Coding/go-tui-file-project-secrets/puffin-transfer-6391b6e60b42.json"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connection to firestore failed")
		os.Exit(1)
	}

	// Retrieve cli file path.
	args := os.Args[1:]
	if len(args) != 2 || args[0] != "transfer" && args[0] != "retrieve" {
		fmt.Fprintf(os.Stderr, "Usage: %s <transfer|retrieve> <filepath>\n", os.Args[0])
		os.Exit(1)
	}

	// Transfer or read in a file.
	command, filepath := args[0], args[1]
	if command == "transfer" {
		file, err := os.ReadFile(filepath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not read %s", filepath)
			os.Exit(1)
		}
		_, _, err = client.Collection("files").Add(ctx, map[string]interface{}{
			"filepath": filepath,
			"content":  file,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not transfer %s to db", filepath)
			os.Exit(1)
		}
	} else {
		q := client.Collection("files").Where("filepath", "==", filepath)
		iter := q.Documents(ctx)
		defer iter.Stop()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not retrieve %s from db", filepath)
			}
			fmt.Println(doc.Data())
		}
	}

	fmt.Fprintf(os.Stdout, "Successfully executed %s of %s!", command, filepath)
	os.Exit(0)
}
