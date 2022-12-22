package controllers

import (
	"gonews/models"
	"gonews/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreatePost(c *gin.Context, dbClient *mongo.Client, username string) {
	post := models.Post{}

	// Bind the request body to the Post struct
	if err := c.ShouldBindJSON(&post); err != nil {
		// If there is an error, return a Bad Request response
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post.CreatedAt, post.UpdatedAt = time.Now(), time.Now()
	post.Author = username

	// Parse hashtags from content
	tags := services.ParseHashtags(post.Content)

	// Iterate through tags and create them in the database if they don't already exist
	for _, tag := range tags {
		models.DbInsertTag(dbClient, tag)
	}

	post.Tags = tags

	// Insert post to database
	dbPost, err := models.DbInsertPost(dbClient, post)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert dbPost to a primitive.ObjectID
	dbPostId := dbPost.(primitive.ObjectID)

	// Iterate through tags and add the post ID to the tag
	for _, tag := range tags {
		err = models.DbAddPostToTag(dbClient, tag, dbPostId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK,
		gin.H{
			"status":  "success",
			"message": "successfully created post",
			"user":    &post,
			"res":     dbPost,
		})
}

// DeleteUser deletes a user from the database with the given username
func DeletePost(c *gin.Context, dbClient *mongo.Client, id string) {
	// Create a filter to find the user by ID
	filter := bson.M{"id": id}

	// Delete the post from the database
	deleteResult, err := models.DbDeletePost(dbClient, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return a success response
	c.JSON(http.StatusOK,
		gin.H{
			"status":  "success",
			"message": "successfully deleted post",
			"res":     deleteResult,
		})
}

// Returns all posts
func ReadPosts(c *gin.Context, dbClient *mongo.Client) {
	posts, err, count := models.DbQueryPosts(dbClient, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  "success",
			"message": "successfully retrieved posts",
			"count":   count,
			"posts":   posts,
		},
	)
}

// Returns all posts from specific user
func ReadUserPosts(c *gin.Context, dbClient *mongo.Client, username string) {
	// Create a filter to find all posts with given author
	filter := bson.M{"author": username}

	posts, err, count := models.DbQueryPosts(dbClient, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Author does not exist"})
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  "success",
			"message": "successfully retrieved user posts",
			"posts":   posts,
		},
	)
}

// Returns all posts with given hasthag
func ReadPostsByTag(c *gin.Context, dbClient *mongo.Client, tag string) {
	// Create a filter to find all posts with given author
	filter := bson.M{"name": tag}

	tags, err, count := models.DbQueryTags(dbClient, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if count == 0 {
		emptyPosts := []models.Post{}
		c.JSON(
			http.StatusOK,
			gin.H{
				"status":  "success",
				"message": "this tag has no posts",
				"posts":   emptyPosts,
			},
		)
		return
	}
	tagObject := tags[0]

	// Get all posts from tagObject.Posts
	postIds := make([]primitive.ObjectID, 0, len(tagObject.Posts))
	for _, id := range tagObject.Posts {
		postIds = append(postIds, id)
	}

	posts, err := models.DbDereferencePosts(dbClient, postIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  "success",
			"message": "successfully retrieved user posts",
			"posts":   posts,
		},
	)
}

// Returns post with specified ID
func ReadSinglePost(c *gin.Context, dbClient *mongo.Client, id string) {
	// Convert the string ID to a primitive ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Create a filter to find the post by ID
	filter := bson.M{"_id": objectID}

	post, err, count := models.DbQueryPosts(dbClient, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post does not exist"})
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  "success",
			"message": "successfully retrieved post",
			"post":    post,
		},
	)
}
