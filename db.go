package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// This is a user defined method to close resources.
// This method closes mongoDB connection and cancel context.
func close(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc) {

	// CancelFunc to cancel to context
	defer cancel()

	// client provides a method to close
	// a mongoDB connection.
	defer func() {

		// client.Disconnect method also has deadline.
		// returns error if any,
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

// This is a user defined method that returns mongo.Client,
// context.Context, context.CancelFunc and error.
// mongo.Client will be used for further database operation.
// context.Context will be used set deadlines for process.
// context.CancelFunc will be used to cancel context and
// resource associated with it.

func connect(uri string) (*mongo.Client, context.Context,
	context.CancelFunc, error) {

	// ctx will be used to set deadline for process, here
	// deadline will of 30 seconds.
	ctx, cancel := context.WithTimeout(context.Background(),
		30*time.Second)

	// mongo.Connect return mongo.Client method
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}

// This is a user defined method that accepts
// mongo.Client and context.Context
// This method used to ping the mongoDB, return error if any.
func ping(client *mongo.Client, ctx context.Context) error {

	// mongo.Client has Ping to ping mongoDB, deadline of
	// the Ping method will be determined by cxt
	// Ping method return error if any occurred, then
	// the error can be handled.
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("connected successfully")
	return nil
}

var client *mongo.Client
var ctx context.Context

var DB *mongo.Database

type Link struct {
	Ref       string `bson:"ref"`
	Timestamp string `bson:"timestamp"`
}

type Song struct {
	Name  string `bson:"name"`
	Links []Link `bson:"links"`
}

type Playlist struct {
	Name  string `bson:"name"`
	Songs []Song `bson:"songs"`
}

func connectDB(dbName string) (*mongo.Database, context.Context, context.CancelFunc) {

	client, ctx, cancel, err := connect("mongodb://mongodb:27017")
	if err != nil {
		panic(err)
	}

	db := client.Database(dbName)

	return db, ctx, cancel
}

func listPlaylists() []Playlist {

	coll := DB.Collection("playlists")

	cur, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return nil
	}

	var result []Playlist
	cur.All(ctx, &result)

	return result

}

func createPlaylist(pl *Playlist) error {

	DB.Collection("playlists").InsertOne(ctx, pl)
	return nil
}

func getPlaylistById(id string) (*Playlist, error) {

	coll := DB.Collection("playlists")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res := coll.FindOne(ctx, bson.M{"_id": objectId})
	if res == nil {
		return nil, errors.New("not found")
	}

	var result Playlist

	res.Decode(&result)

	return &result, nil

}

func addSongsToPlaylist(pid string, songs []Song) error {

	pl, err := getPlaylistById(pid)

	coll := DB.Collection("playlists")

	objectId, err := primitive.ObjectIDFromHex(pid)
	if err != nil {
		return err
	}

	pl.Songs = append(pl.Songs, songs...)

	coll.UpdateByID(ctx, objectId, pl)

	return nil
}

func Init() {

	DB, ctx, _ = connectDB("playlists_db")

}
