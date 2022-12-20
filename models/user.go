package models

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID        primitive.ObjectID
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

func DbInsertUser(client *mongo.Client, user User) (interface{}, error) {
	collection := client.Database(os.Getenv("MONGODB_DATABASE")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if a user with the same username already exists
	filter := bson.M{"username": user.Username}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Return an error if a user with the same username already exists
	if count > 0 {
		return nil, fmt.Errorf("User with the same username already exists")
	}

	res, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Fatalln("Error inserting new user", err)
	}
	return res.InsertedID, nil
}

func DbUpdateUser(client *mongo.Client, filter bson.M, newUser User) (interface{}, error) {
	collection := client.Database(os.Getenv("MONGODB_DATABASE")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{"$set": newUser}
	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	} else if res.ModifiedCount == 0 {
		return nil, fmt.Errorf("User does not exist")
	}

	return res.ModifiedCount, nil
}

// DbDeleteUser deletes a user from the database with the given filter
func DbDeleteUser(client *mongo.Client, filter bson.M) (interface{}, error) {
	collection := client.Database(os.Getenv("MONGODB_DATABASE")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	} else if res.DeletedCount == 0 {
		return nil, fmt.Errorf("User does not exist")
	}

	return res.DeletedCount, nil
}
