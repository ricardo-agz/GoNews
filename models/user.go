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
	ID        primitive.ObjectID `bson:"_id"`
	Username  string             `bson:"username"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type Users []*User

// Returns all users in the database with the matching filter
func DbQueryUsers(client *mongo.Client, filter bson.M) (Users, error, int) {
	var users Users
	collection := client.Database(os.Getenv("MONGODB_DATABASE")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatal("Error retrieving documents", err)
		return nil, err, 0
	}

	var outerror error
	count := 0

	for cur.Next(ctx) {
		var user User
		count += 1
		err = cur.Decode(&user)
		if err != nil {
			log.Fatal("Error decoding document", err)
			outerror = err
		}
		users = append(users, &user)
	}
	return users, outerror, count
}

// DbCreateUser creates a user in the database with the given user data
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

	// Initialize user id
	user.ID = primitive.NewObjectID()

	res, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Fatalln("Error inserting new user", err)
	}
	return res.InsertedID, nil
}

// DbUpdateUser updates a user from the database with the given filter and new user data
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
