package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connectionString = "mongodb://mongodb:27017"
	dbName           = "playlists_db"
	collectionName   = "items"
)

func run() {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//mongo.Connect(ctx)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database(dbName).Collection(collectionName)

	// Create
	id := createItem(collection, Item{Name: "Sample Item"})
	fmt.Println("Created item with ID:", id)

	// Read
	readItem := getItem(collection, id)
	fmt.Println("Read item:", readItem)

	// Update
	updateItem(collection, id, "Updated Item")

	fmt.Println(getAllItems(collection))

	// Delete
	deleteItem(collection, id)
}

type Item struct {
	ID   string `bson:"_id,omitempty"`
	Name string `bson:"name,omitempty"`
}

func createItem(collection *mongo.Collection, item Item) string {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.InsertOne(ctx, item)
	if err != nil {
		log.Fatal("Error on inserting new item: ", err)
	}

	return result.InsertedID.(primitive.ObjectID).Hex()
}

func getItem(collection *mongo.Collection, id string) *Item {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil
	}

	result := collection.FindOne(ctx, bson.M{"_id": objectId})
	item := &Item{}
	err = result.Decode(item)
	if err != nil {
		log.Println("Error on getting one item: ", err)
		return nil
	}
	return item
}

func updateItem(collection *mongo.Collection, id, name string) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"name": name}})
	if err != nil {
		log.Println("Error on updating item: ", err)
	}
}

func getAllItems(collection *mongo.Collection) []Item {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, err := collection.Find(ctx, bson.M{}) // Find with an empty filter to match all documents
	if err != nil {
		log.Fatal("Error on getting all items: ", err)
	}
	defer cursor.Close(ctx)

	var items []Item
	for cursor.Next(ctx) {
		var item Item
		cursor.Decode(&item)
		items = append(items, item)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal("Cursor error: ", err)
	}
	return items
}

func deleteItem(collection *mongo.Collection, id string) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		log.Println("Error on deleting item: ", err)
	}
}
