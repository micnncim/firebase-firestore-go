package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
)

type State struct {
	Capital string
	Population float32
}

type FirestoreClient struct {
	*firestore.Client
}

func main() {
	ctx := context.Background()
	client, err := NewFirestoreClient(ctx)
	if err != nil {
        log.Fatalln(err)
	}
	s, err := client.Read(ctx, "NewYork")
	if err != nil {
		return
	}
	fmt.Printf("old value: %#v\n", s)
	s.Population += 1
	if err := client.Write(ctx, "NewYork", s); err != nil {
		return
	}
	fmt.Printf("new value: %#v\n", s)
}

func NewFirestoreClient(ctx context.Context) (*FirestoreClient, error) {
	projectID := os.Getenv("PROJECT_ID")
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
        return nil, err
	}
	return &FirestoreClient{
		client,
	}, nil
}

func (c *FirestoreClient) Write(ctx context.Context, id string, s *State) error {
	collection := c.Collection("States")
	doc := collection.Doc(id)
	if _, err := doc.Set(ctx, s); err != nil {
		return err
	}
	return nil
}

func (c *FirestoreClient) Read(ctx context.Context, id string) (*State, error) {
	collection := c.Collection("States")
	doc := collection.Doc(id)
	snapshot, err := doc.Get(ctx)
	if err != nil {
		return nil, err
	}
	var s State
	if err := snapshot.DataTo(&s); err != nil {
        return nil, err
	}
	return &s, nil
}
