package database

import (
	"context"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "github.com/lib/pq"
)

func ConnectToMongoDB(log *logrus.Entry) (*mongo.Client, error) {
	log.Infof("Connecting to MongoDB...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Errorf("MongoDB URI format incorrect: %v", err)
		return nil, err
	}

	// Test the connection by pinging the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Errorf("Could not ping MongoDB: %v", err)
		return nil, err
	}

	log.Infof("Successfully connected to MongoDB")
	return client, nil

}
