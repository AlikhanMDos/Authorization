package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type User struct {
	Name     string `json:"name"`
	Phone    string `json:"phoneNumber"`
	Password string `json:"password"`
}

var userCollection *mongo.Collection

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
	log.Println("Connected to MongoDB!")
}
