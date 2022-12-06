package models

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Username  string    `bson:"username"`
	Email     string    `bson:"email"`
	Password  string    `bson:"password"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type Users []*User

func DbQueryUsers(client *mongo.Client, filter bson.M) Users {
	var users Users
	collection := client.Database(os.Getenv("MONGODB_DATABASE")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatal("Error retrieving documents", err)
	}

	for cur.Next(ctx) {
		var user User
		err = cur.Decode(&user)
		if err != nil {
			log.Fatal("Error decoding document", err)
		}
		users = append(users, &user)
	}
	return users
}

func DbInsertUser(client *mongo.Client, user User) interface{} {
	collection := client.Database(os.Getenv("MONGODB_DATABASE")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Fatalln("Error inserting new user", err)
	}
	return res.InsertedID
}
