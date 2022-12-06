package main

import (
	"fmt"
	"context"
	"time"
	"os"
	"log"
	"net/http"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)


var mongoConn *mongo.Client

type User struct {
	Username       	string         	`bson:"username"`
	Email       	string         	`bson:"email"`
	Password       	string         	`bson:"password"`
	CreatedAt    	time.Time     	`bson:"created_at"`
	UpdatedAt   	time.Time       `bson:"updated_at"`
}

type Users []Users

// DB connection
func createConnection() (*mongo.Client, error) {
	fmt.Println(os.Getenv("MONGODB_USERNAME"))

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
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"status": "failed",
					"message": "invalid request body",
				},
			)
			return
		}

		user.CreatedAt, user.UpdatedAt = time.Now(), time.Now()

		collection := client.Database("gonews_users").Collection("users")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := collection.InsertOne(ctx, user)

		c.JSON(
			http.StatusOK,
			gin.H{
				"status": "success",
				"user": &user,
				"res": res,
			},
		)
	})

	// Users List
	router.GET("/users", func(c *gin.Context) {
		// users := Users{}

		collection := client.Database("gonews_users").Collection("users")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		cur, err := collection.Find(ctx, bson.D{})
		if err != nil { 
			fmt.Println("Connecting to database...")
			log.Fatal(err) 
		}

		defer cur.Close(ctx)
		for cur.Next(ctx) {
			var result bson.D
			err := cur.Decode(&result)
			if err != nil { log.Fatal(err) }
			// do something with result....
		}
		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}

		// session := connect()
		// defer session.Close()
		// err := session.DB(database).C(collection).Find(bson.M{}).All(&users)

		// if err != nil {
		// 	c.JSON(
		// 		http.StatusNotFound,
		// 		gin.H{
		// 			"status": "failed",
		// 			"message": "error getting users",
		// 		})
		// 	return
		// }

		c.JSON(
			http.StatusOK, 
			gin.H{
				"status": "success", 
				"users": collection,
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