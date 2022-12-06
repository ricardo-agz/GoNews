package database

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB connection
func CreateConnection() (*mongo.Client, error) {
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
