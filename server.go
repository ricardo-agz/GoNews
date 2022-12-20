package main

import (
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	// "go.mongodb.org/mongo-driver/mongo/readpref"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	db "gonews/config"
	"gonews/controllers"
)

var mongoConn *mongo.Client

// StartService function
func StartService(dbClient *mongo.Client) {
	router := gin.Default()

	// Home
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to GoNews!",
		})
	})

	// Users List
	router.GET("/users", func(c *gin.Context) {
		controllers.ReadUsers(c, dbClient)
	})

	// Get Single User
	router.GET("/users/:username", func(c *gin.Context) {
		username := c.Param("username")
		controllers.ReadSingleUser(c, dbClient, username)
	})

	// User Create
	router.POST("/users", func(c *gin.Context) {
		controllers.CreateUser(c, dbClient)
	})

	// User Update
	router.PUT("/users/:username", func(c *gin.Context) {
		username := c.Param("username")
		controllers.UpdateUser(c, dbClient, username)
	})

	// User Delete
	router.DELETE("/users/:username", func(c *gin.Context) {
		username := c.Param("username")
		controllers.DeleteUser(c, dbClient, username)
	})

	// 404 Not found
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Page not found:/",
		})
	})

	router.Run(":8000")
}

func main() {
	enverr := godotenv.Load()
	if enverr != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("Connecting to database...")
	var err error
	mongoConn, err = db.CreateConnection()
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting Server...")
	StartService(mongoConn)
}
