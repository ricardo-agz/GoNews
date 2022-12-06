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

	// User Create
	router.POST("/users", func(c *gin.Context) {
		controllers.CreateUser(c, dbClient)
	})

	// Users List
	router.GET("/users", func(c *gin.Context) {
		controllers.ReadUsers(c, dbClient)
	})

	// 404 Not found
	router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
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
