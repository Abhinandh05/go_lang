package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")
	if MONGODB_URI == "" {
		log.Fatal("MONGODB_URI not found in environment variables")
	}

	fmt.Println("Connecting to MongoDB...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use mongo.Connect directly (mongo.NewClient is deprecated)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MONGODB_URI))
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	// Ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("❌ Could not ping MongoDB:", err)
	}

	fmt.Println("Connected to MongoDB Successfully!")

	DB = client
}

// GetCollection returns a collection from the database
func GetCollection(name string) *mongo.Collection {
	if DB == nil {
		log.Fatal("Database not initialized. Call ConnectDB() first")
	}
	return DB.Database("go-auth").Collection(name)
}

// DisconnectDB closes the database connection
func DisconnectDB() {
	if DB != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := DB.Disconnect(ctx); err != nil {
			log.Fatal("❌ Error disconnecting from MongoDB:", err)
		}
		fmt.Println("Disconnected from MongoDB")
	}
}
