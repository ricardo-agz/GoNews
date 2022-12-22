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

type Post struct {
	ID        primitive.ObjectID `bson:"_id"`
	Author    string             `bson:"author"`
	Content   string             `bson:"content"`
	Tags      []string           `bson:"tags"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type Posts []*Post

// Given a list of postIds, returns a list of post objects
func DbDereferencePosts(dbClient *mongo.Client, postIds []primitive.ObjectID) ([]Post, error) {
	postCollection := dbClient.Database(os.Getenv("MONGODB_DATABASE")).Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	postsCursor, err := postCollection.Find(ctx, bson.M{"_id": bson.M{"$in": postIds}})
	if err != nil {
		return nil, err
	}

	var posts []Post
	for postsCursor.Next(ctx) {
		var post Post
		if err := postsCursor.Decode(&post); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := postsCursor.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// Returns all posts in the database with the matching filter
func DbQueryPosts(client *mongo.Client, filter bson.M) (Posts, error, int) {
	var posts Posts
	collection := client.Database(os.Getenv("MONGODB_DATABASE")).Collection("posts")
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
		var post Post
		count += 1
		err = cur.Decode(&post)
		if err != nil {
			log.Fatal("Error decoding document", err)
			outerror = err
		}
		posts = append(posts, &post)
	}
	return posts, outerror, count
}

// Creates a post in the database with the given post data
func DbInsertPost(client *mongo.Client, post Post) (interface{}, error) {
	postCollection := client.Database(os.Getenv("MONGODB_DATABASE")).Collection("posts")
	userCollection := client.Database(os.Getenv("MONGODB_DATABASE")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if the post's author exists
	filter := bson.M{"username": post.Author}
	count, err := userCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Return an error if author does not exist
	if count == 0 {
		return nil, fmt.Errorf("Invalid author")
	}

	// Initialize post id
	post.ID = primitive.NewObjectID()

	// Parse tags from content

	res, err := postCollection.InsertOne(ctx, post)
	if err != nil {
		log.Fatalln("Error inserting new post", err)
	}
	return res.InsertedID, nil
}

// DbDeletePost deletes a user from the database with the given filter
func DbDeletePost(client *mongo.Client, filter bson.M) (interface{}, error) {
	collection := client.Database(os.Getenv("MONGODB_DATABASE")).Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	} else if res.DeletedCount == 0 {
		return nil, fmt.Errorf("Post does not exist")
	}

	return res.DeletedCount, nil
}
