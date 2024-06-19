package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type User struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type Post struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title    string             `json:"title" bson:"title"`
	Text     string             `json:"text" bson:"text"`
	Image    string             `json:"image" bson:"image"`
	Author   string             `json:"author" bson:"author"`
	Date     time.Time          `json:"date" bson:"date"`
	AuthorID primitive.ObjectID `json:"author_id" bson:"author_id"`
}

var userCollection *mongo.Collection
var postCollection *mongo.Collection

func initDB() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://220727:1234567899@cluster0.authau1.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	userCollection = client.Database("auth").Collection("users")
	postCollection = client.Database("auth").Collection("posts")

	log.Println("Connected to MongoDB!")
}
