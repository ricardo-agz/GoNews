package controllers

import (
	"gonews/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateUser(c *gin.Context, dbClient *mongo.Client) {
	// body, _ := ioutil.ReadAll(c.Request.Body)
	user := models.User{}

	// Bind the request body to the User struct
	if err := c.ShouldBindJSON(&user); err != nil {
		// If there is an error, return a Bad Request response
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.CreatedAt, user.UpdatedAt = time.Now(), time.Now()
	dbUser, err := models.DbInsertUser(dbClient, user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK,
		gin.H{
			"status":  "success",
			"message": "successfully created user",
			"user":    &user,
			"res":     dbUser,
		})
}

func UpdateUser(c *gin.Context, dbClient *mongo.Client, username string) {
	user := models.User{}

	// Bind the request body to the User struct
	if err := c.ShouldBindJSON(&user); err != nil {
		// If there is an error, return a Bad Request response
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the updated_at field
	user.UpdatedAt = time.Now()

	// Create a filter to find the user by username
	filter := bson.M{"username": username}

	// Update the user in the database
	updateResult, err := models.DbUpdateUser(dbClient, filter, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return a success response
	c.JSON(http.StatusOK,
		gin.H{
			"status":  "success",
			"message": "successfully updated user",
			"user":    &user,
			"res":     updateResult,
		})
}

// DeleteUser deletes a user from the database with the given ID
func DeleteUser(c *gin.Context, dbClient *mongo.Client, username string) {
	// Create a filter to find the user by ID
	filter := bson.M{"username": username}

	// Delete the user from the database
	deleteResult, err := models.DbDeleteUser(dbClient, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return a success response
	c.JSON(http.StatusOK,
		gin.H{
			"status":  "success",
			"message": "successfully deleted user",
			"res":     deleteResult,
		})
}

func ReadUsers(c *gin.Context, dbClient *mongo.Client) {
	users := models.DbQueryUsers(dbClient, bson.M{})

	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  "success",
			"message": "successfully retrieved users",
			"users":   users,
		},
	)
}
