package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// "go.mongodb.org/mongo-driver/mongo/readpref"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var mongoConn *mongo.Client

type User struct {
	Username  string    `bson:"username"`
	Email     string    `bson:"email"`
	Password  string    `bson:"password"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type Users []*User

func QueryUsers(client *mongo.Client, filter bson.M) []*User {
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

func InsertUser(client *mongo.Client, user User) interface{} {
	collection := client.Database(os.Getenv("MONGODB_DATABASE")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Fatalln("Error inserting new user", err)
	}
	return res.InsertedID
}

// DB connection
func createConnection() (*mongo.Client, error) {
	credential := options.Credential{
		AuthMechanism: "SCRAM-SHA-1",
		Username:      os.Getenv("MONGODB_USERNAME"), // mongodb user
		Password:      os.Getenv("MONGODB_PASSWORD"),
	}
	connString := os.Getenv("MONGODB_URL")
	clientOpts := options.Client().ApplyURI(connString).SetAuth(credential)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return mongo.Connect(ctx, clientOpts)
}

// StartService function
func StartService(client *mongo.Client) {
	router := gin.Default()

	// Home
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to GoNews!",
		})
	})

	// User Create
	router.POST("/users", func(c *gin.Context) {
		user := User{}
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
		dbUser := InsertUser(client, user)

		c.JSON(http.StatusOK,
			gin.H{
				"status": "success",
				"user":   &user,
				"res":    dbUser,
			})
	})

	// Users List
	router.GET("/users", func(c *gin.Context) {
		users := QueryUsers(client, bson.M{})

		c.JSON(
			http.StatusOK,
			gin.H{
				"status": "success",
				"users":  users,
			},
		)
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
	mongoConn, err = createConnection()
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting Server...")
	StartService(mongoConn)
}
