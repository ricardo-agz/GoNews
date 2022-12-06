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
	user := models.User{}
	err := c.Bind(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"status":  "failed",
				"message": "invalid request body",
			})
		return
	}

	user.CreatedAt, user.UpdatedAt = time.Now(), time.Now()
	dbUser := models.DbInsertUser(dbClient, user)

	c.JSON(http.StatusOK,
		gin.H{
			"status": "success",
			"user":   &user,
			"res":    dbUser,
		})
}

func ReadUsers(c *gin.Context, dbClient *mongo.Client) {
	users := models.DbQueryUsers(dbClient, bson.M{})

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": "success",
			"users":  users,
		},
	)
}
