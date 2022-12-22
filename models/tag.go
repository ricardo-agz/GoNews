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

type Tag struct {
	ID    primitive.ObjectID   `bson:"_id"`
	Name  string               `bson:"name"`
	Posts []primitive.ObjectID `bson:"posts"`
}

type Tags []*Tag

// Returns all tags in the database with the matching filter
func DbQueryTags(client *mongo.Client, filter bson.M) (Tags, error, int) {
	var tags Tags
	collection := client.Database(os.Getenv("MONGODB_DATABASE")).Collection("tags")
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
		var tag Tag
		count += 1
		err = cur.Decode(&tag)
		if err != nil {
			log.Fatal("Error decoding document", err)
			outerror = err
		}
		tags = append(tags, &tag)
	}
	return tags, outerror, count
}

// DbInsertTag creates a tag in the database with the given tagname
func DbInsertTag(client *mongo.Client, tagname string) (interface{}, error) {
	collection := client.Database(os.Getenv("MONGODB_DATABASE")).Collection("tags")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if a tag with the same name already exists
	filter := bson.M{"name": tagname}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Return an error if a the tag already exists
	if count > 0 {
		return nil, fmt.Errorf("Tag already exists")
	}

	// Initialize tag object
	tag := Tag{
		ID:    primitive.NewObjectID(),
		Name:  tagname,
		Posts: []primitive.ObjectID{},
	}

	res, err := collection.InsertOne(ctx, tag)
	if err != nil {
		log.Fatalln("Error inserting new tag", err)
	}
	return res.InsertedID, nil
}

// DbCreateUser creates a user in the database with the given user data
func DbAddPostToTag(client *mongo.Client, tagname string, postId primitive.ObjectID) error {
	collection := client.Database(os.Getenv("MONGODB_DATABASE")).Collection("tags")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Return tag with specified tagname
	filter := bson.M{"name": tagname}
	update := bson.M{"$push": bson.M{"posts": postId}}

	// Save tag to database
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
