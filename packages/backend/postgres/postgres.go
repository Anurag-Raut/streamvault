package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	// adapt "demo" to your module name if it differs
	"streamvault/postgres/db"
)

var client = db.NewClient()

// func main() {
// 	if err := run(); err != nil {
// 		panic(err)
// 	}
// }

// func run() error {

// 	if err := client.Prisma.Connect(); err != nil {
// 		return err
// 	}

// 	defer func() {
// 		if err := client.Prisma.Disconnect(); err != nil {
// 			panic(err)
// 		}
// 	}()

// 	ctx := context.Background()

// 	// create a post
// 	createdPost, err := client.Post.CreateOne(
// 		db.Post.Title.Set("Hi from Prisma!"),
// 		db.Post.Published.Set(true),
// 		db.Post.Desc.Set("Prisma is a database toolkit and makes databases easy."),
// 	).Exec(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	result, _ := json.MarshalIndent(createdPost, "", "  ")
// 	fmt.Printf("created post: %s\n", result)

// 	// find a single post
// 	post, err := client.Post.fin
// 	if err != nil {
// 		return err
// 	}

// 	result, _ = json.MarshalIndent(post, "", "  ")
// 	fmt.Printf("post: %s\n", result)

// 	// for optional/nullable values, you need to check the function and create two return values
// 	// `desc` is a string, and `ok` is a bool whether the record is null or not. If it's null,
// 	// `ok` is false, and `desc` will default to Go's default values; in this case an empty string (""). Otherwise,
// 	// `ok` is true and `desc` will be "my description".
// 	desc, ok := post.Desc()
// 	if !ok {
// 		return fmt.Errorf("post's description is null")
// 	}

// 	fmt.Printf("The posts's description is: %s\n", desc)

// 	return nil
// }

func Connect() error {
	if err := client.Prisma.Connect(); err != nil {
		return err
	}

	fmt.Println("Connected to the database")

	return nil
}

func Disconnect() error {
	if err := client.Prisma.Disconnect(); err != nil {
		return err
	}
	fmt.Println("Disconnected from the database")

	return nil
}

func AddStream(title string) (string, error) {

	ctx := context.Background()

	addedStream, err := client.Stream.CreateOne(
		db.Stream.Title.Set(title),
	).Exec(ctx)

	if err != nil {
		return "", err
	}

	result, _ := json.MarshalIndent(addedStream, "", "  ")
	fmt.Printf("Stream: %s\n", result)

	return string(addedStream.ID), nil

}
